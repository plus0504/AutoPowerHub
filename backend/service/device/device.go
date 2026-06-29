package device

import (
	"errors"
	"fmt"
	"time"

	"autopowerhub/models"
	"autopowerhub/repository"
	"autopowerhub/service/ble"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceDisabled = errors.New("device is disabled")
)

type Service struct {
	deviceRepo *repository.DeviceRepository
	logRepo    *repository.LogRepository
	bleMgr     *ble.Manager
}

func NewService(
	deviceRepo *repository.DeviceRepository,
	logRepo *repository.LogRepository,
	bleMgr *ble.Manager,
) *Service {
	return &Service{deviceRepo: deviceRepo, logRepo: logRepo, bleMgr: bleMgr}
}

func (s *Service) ListDevices() ([]models.Device, error) {
	return s.deviceRepo.FindAll()
}

func (s *Service) SendCommand(deviceID uint, username, command string) error {
	dev, err := s.deviceRepo.FindByID(deviceID)
	if err != nil {
		return fmt.Errorf("%w: id=%d", ErrDeviceNotFound, deviceID)
	}
	if !dev.Enabled {
		return ErrDeviceDisabled
	}

	bleErr := s.bleMgr.SendCommand(dev.MAC, dev.ServiceUUID, dev.CharacteristicUUID, command)

	result := "success"
	if bleErr != nil {
		result = bleErr.Error()
	}

	// Write audit log regardless of BLE outcome.
	_ = s.logRepo.Create(&models.Log{
		Username:  username,
		Device:    dev.Name,
		Command:   command,
		Result:    result,
		CreatedAt: time.Now(),
	})

	return bleErr
}
