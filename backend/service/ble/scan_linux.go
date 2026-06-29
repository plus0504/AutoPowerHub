//go:build linux

package ble

import (
	"strings"

	"tinygo.org/x/bluetooth"
)

type ScanResult struct {
	ServiceUUID         string   `json:"service_uuid"`
	CharacteristicUUIDs []string `json:"characteristic_uuids"`
}

// DumpServices connects to a device by MAC and returns all service/characteristic UUIDs.
// Intended for debugging UUID mismatches; does not cache the connection.
func (m *Manager) DumpServices(mac string) ([]ScanResult, error) {
	l := m.impl.(*linuxImpl)
	l.mu.Lock()
	defer l.mu.Unlock()

	key := strings.ToUpper(mac)
	device, err := l.connect(mac, key)
	if err != nil {
		return nil, err
	}
	defer device.Disconnect() //nolint:errcheck

	srvcs, err := device.DiscoverServices(nil)
	if err != nil {
		return nil, err
	}

	results := make([]ScanResult, 0, len(srvcs))
	for _, svc := range srvcs {
		r := ScanResult{ServiceUUID: svc.UUID().String()}
		chars, _ := svc.DiscoverCharacteristics(nil)
		for _, c := range chars {
			r.CharacteristicUUIDs = append(r.CharacteristicUUIDs, c.UUID().String())
		}
		results = append(results, r)
	}
	return results, nil
}
