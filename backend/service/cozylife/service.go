package cozylife

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const defaultTimeout = 3 * time.Second

var ErrSwitchNotFound = errors.New("switch not found")

type Service struct {
	devicesPath string
}

func NewService(devicesPath string) *Service {
	return &Service{devicesPath: devicesPath}
}

func (s *Service) loadSwitches() ([]*ScanResult, error) {
	if s.devicesPath == "" {
		return nil, nil
	}
	data, err := os.ReadFile(s.devicesPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", s.devicesPath, err)
	}
	var all []*ScanResult
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, fmt.Errorf("parse devices: %w", err)
	}
	out := make([]*ScanResult, 0, len(all))
	for _, d := range all {
		if d.DeviceTypeCode == switchTypeCode {
			out = append(out, d)
		}
	}
	return out, nil
}

func (s *Service) ListSwitches() ([]*ScanResult, error) {
	return s.loadSwitches()
}

func (s *Service) findSwitch(ip string) (*ScanResult, error) {
	switches, err := s.loadSwitches()
	if err != nil {
		return nil, err
	}
	for _, d := range switches {
		if d.IP == ip {
			return d, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrSwitchNotFound, ip)
}

func (s *Service) TurnOn(ip string) error {
	return s.runCmd(ip, func(d *device) error { return d.turnOn() })
}

func (s *Service) TurnOff(ip string) error {
	return s.runCmd(ip, func(d *device) error { return d.turnOff() })
}

func (s *Service) runCmd(ip string, fn func(*device) error) error {
	r, err := s.findSwitch(ip)
	if err != nil {
		return err
	}
	d := newDeviceFromResult(r, defaultTimeout)
	if err := d.cli.connect(); err != nil {
		return fmt.Errorf("connect %s: %w", ip, err)
	}
	defer d.cli.disconnect()
	return fn(d)
}
