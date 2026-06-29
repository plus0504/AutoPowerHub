package ble

// Manager wraps platform-specific BLE communication.
// On Linux it uses tinygo.org/x/bluetooth via BlueZ D-Bus.
// On other platforms it returns a stub error (build on Raspberry Pi for full functionality).
type Manager struct {
	impl bleImpl
}

type bleImpl interface {
	sendCommand(mac, serviceUUID, charUUID, command string) error
	close() error
}

func (m *Manager) SendCommand(mac, serviceUUID, charUUID, command string) error {
	return m.impl.sendCommand(mac, serviceUUID, charUUID, command)
}

func (m *Manager) Close() error {
	return m.impl.close()
}
