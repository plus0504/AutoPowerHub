package models

import "time"

type Log struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	Username  string    `json:"username"`
	Device    string    `json:"device"`
	Command   string    `json:"command"`
	Result    string    `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}
