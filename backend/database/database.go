package database

import (
	"fmt"

	"autopowerhub/config"
	"autopowerhub/models"

	"golang.org/x/crypto/bcrypt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.SQLite.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Device{}, &models.Log{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	if err := seedAdmin(db, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		return nil, fmt.Errorf("seed admin: %w", err)
	}

	if err := seedDevices(db, cfg.Devices); err != nil {
		return nil, fmt.Errorf("seed devices: %w", err)
	}

	return db, nil
}

// seedDevices inserts devices from config that don't yet exist in the DB (matched by MAC).
// Existing DB records are left untouched so manual edits are preserved.
func seedDevices(db *gorm.DB, devices []config.DeviceConfig) error {
	for _, d := range devices {
		var count int64
		db.Model(&models.Device{}).Where("mac = ?", d.MAC).Count(&count)
		if count > 0 {
			continue
		}
		if err := db.Create(&models.Device{
			Name:               d.Name,
			MAC:                d.MAC,
			ServiceUUID:        d.ServiceUUID,
			CharacteristicUUID: d.CharacteristicUUID,
			Enabled:            d.Enabled,
		}).Error; err != nil {
			return fmt.Errorf("insert device %q: %w", d.Name, err)
		}
	}
	return nil
}

func seedAdmin(db *gorm.DB, username, password string) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Create(&models.User{
		Username: username,
		Password: string(hash),
	}).Error
}
