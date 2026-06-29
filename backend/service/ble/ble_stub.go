//go:build !linux

package ble

import "errors"

type stubImpl struct{}

// NewManager returns a stub manager on non-Linux platforms.
// Deploy on Raspberry Pi OS (Linux) for real BLE functionality.
func NewManager() (*Manager, error) {
	return &Manager{impl: &stubImpl{}}, nil
}

func (s *stubImpl) sendCommand(_, _, _, _ string) error {
	return errors.New("BLE is only supported on Linux (Raspberry Pi); this is a stub")
}

func (s *stubImpl) close() error {
	return nil
}
