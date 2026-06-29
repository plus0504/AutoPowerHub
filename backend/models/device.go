package models

type Device struct {
	ID                 uint   `json:"id"                  gorm:"primaryKey"`
	Name               string `json:"name"                gorm:"not null"`
	MAC                string `json:"mac"                 gorm:"not null"`
	ServiceUUID        string `json:"service_uuid"        gorm:"not null"`
	CharacteristicUUID string `json:"characteristic_uuid" gorm:"not null"`
	Enabled            bool   `json:"enabled"`
}
