package repository

import (
	"autopowerhub/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) FindAll() ([]models.Device, error) {
	var devices []models.Device
	if err := r.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *DeviceRepository) FindByID(id uint) (*models.Device, error) {
	var device models.Device
	if err := r.db.First(&device, id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
