//go:build linux

package ble

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

type charKey struct{ svc, chr string }

type connEntry struct {
	device bluetooth.Device
	chars  map[charKey]bluetooth.DeviceCharacteristic
}

type linuxImpl struct {
	mu        sync.Mutex
	adapter   *bluetooth.Adapter
	addrCache map[string]bluetooth.Address // upper(MAC) → Address
	connMap   map[string]*connEntry        // upper(MAC) → active connection
}

func NewManager() (*Manager, error) {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("enable BLE adapter: %w", err)
	}
	return &Manager{impl: &linuxImpl{
		adapter:   adapter,
		addrCache: make(map[string]bluetooth.Address),
		connMap:   make(map[string]*connEntry),
	}}, nil
}

func (l *linuxImpl) sendCommand(mac, serviceUUID, charUUID, command string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	svcUUID, err := bluetooth.ParseUUID(serviceUUID)
	if err != nil {
		return fmt.Errorf("parse service UUID: %w", err)
	}
	charUUIDParsed, err := bluetooth.ParseUUID(charUUID)
	if err != nil {
		return fmt.Errorf("parse characteristic UUID: %w", err)
	}

	key := strings.ToUpper(mac)
	ckey := charKey{svc: serviceUUID, chr: charUUID}

	// Try the cached connection and characteristic first (fastest path).
	if entry, ok := l.connMap[key]; ok {
		if c, ok := entry.chars[ckey]; ok {
			if _, err := c.WriteWithoutResponse([]byte(command)); err == nil {
				return nil
			}
		}
		// Write failed – connection has dropped; clean up and fall through to reconnect.
		entry.device.Disconnect() //nolint:errcheck
		delete(l.connMap, key)
	}

	// Connect, discover services/chars, and cache the connection.
	entry, err := l.connectAndDiscover(mac, key, svcUUID, charUUIDParsed, serviceUUID, charUUID)
	if err != nil {
		return err
	}
	l.connMap[key] = entry

	if _, err := entry.chars[ckey].WriteWithoutResponse([]byte(command)); err != nil {
		entry.device.Disconnect() //nolint:errcheck
		delete(l.connMap, key)
		return fmt.Errorf("write command: %w", err)
	}

	return nil
}

func (l *linuxImpl) connectAndDiscover(
	mac, key string,
	svcUUID, charUUIDParsed bluetooth.UUID,
	serviceUUID, charUUID string,
) (*connEntry, error) {
	device, err := l.connect(mac, key)
	if err != nil {
		return nil, err
	}

	srvcs, err := device.DiscoverServices([]bluetooth.UUID{svcUUID})
	if err != nil || len(srvcs) == 0 {
		device.Disconnect() //nolint:errcheck
		all, _ := device.DiscoverServices(nil)
		uuids := make([]string, 0, len(all))
		for _, s := range all {
			uuids = append(uuids, s.UUID().String())
		}
		if len(uuids) == 0 {
			return nil, fmt.Errorf("service %s not found (no services discovered)", serviceUUID)
		}
		return nil, fmt.Errorf("service %s not found; device exposes: %v", serviceUUID, uuids)
	}

	chars, err := srvcs[0].DiscoverCharacteristics([]bluetooth.UUID{charUUIDParsed})
	if err != nil || len(chars) == 0 {
		device.Disconnect() //nolint:errcheck
		all, _ := srvcs[0].DiscoverCharacteristics(nil)
		uuids := make([]string, 0, len(all))
		for _, c := range all {
			uuids = append(uuids, c.UUID().String())
		}
		return nil, fmt.Errorf("characteristic %s not found; service exposes: %v", charUUID, uuids)
	}

	ckey := charKey{svc: serviceUUID, chr: charUUID}
	return &connEntry{
		device: device,
		chars:  map[charKey]bluetooth.DeviceCharacteristic{ckey: chars[0]},
	}, nil
}

// connect returns a bluetooth.Device, using a cached address to skip the scan
// when possible. Falls back to scanning if the cached address fails to connect.
func (l *linuxImpl) connect(mac, key string) (bluetooth.Device, error) {
	if addr, ok := l.addrCache[key]; ok {
		if dev, err := l.adapter.Connect(addr, bluetooth.ConnectionParams{}); err == nil {
			return dev, nil
		}
		delete(l.addrCache, key)
	}

	addr, err := l.scan(mac)
	if err != nil {
		return bluetooth.Device{}, err
	}
	l.addrCache[key] = addr

	dev, err := l.adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return bluetooth.Device{}, fmt.Errorf("connect to %s: %w", mac, err)
	}
	return dev, nil
}

// scan performs a BLE scan until the target MAC is found or 30 s elapse.
func (l *linuxImpl) scan(mac string) (bluetooth.Address, error) {
	var addr bluetooth.Address
	found := make(chan struct{}, 1)
	scanDone := make(chan error, 1)

	go func() {
		scanDone <- l.adapter.Scan(func(_ *bluetooth.Adapter, result bluetooth.ScanResult) {
			if strings.EqualFold(result.Address.String(), mac) {
				addr = result.Address
				select {
				case found <- struct{}{}:
				default:
				}
				l.adapter.StopScan() //nolint:errcheck
			}
		})
	}()

	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	select {
	case <-found:
	case <-timer.C:
		l.adapter.StopScan() //nolint:errcheck
		<-scanDone
		return bluetooth.Address{}, fmt.Errorf("device %s not found within scan timeout", mac)
	case err := <-scanDone:
		if err != nil {
			return bluetooth.Address{}, fmt.Errorf("scan: %w", err)
		}
		return bluetooth.Address{}, fmt.Errorf("device %s not found", mac)
	}
	<-scanDone
	return addr, nil
}

func (l *linuxImpl) close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for key, entry := range l.connMap {
		entry.device.Disconnect() //nolint:errcheck
		delete(l.connMap, key)
	}
	return nil
}
