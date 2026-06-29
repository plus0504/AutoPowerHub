//go:build !linux

package ble

import "errors"

type ScanResult struct {
	ServiceUUID         string   `json:"service_uuid"`
	CharacteristicUUIDs []string `json:"characteristic_uuids"`
}

func (m *Manager) DumpServices(_ string) ([]ScanResult, error) {
	return nil, errors.New("BLE not supported on this platform")
}
