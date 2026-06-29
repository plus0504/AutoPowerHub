package models

type User struct {
	ID       uint   `json:"-"        gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"-"        gorm:"not null"`
}
