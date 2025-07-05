package models

import (
    "gorm.io/gorm"
    "time"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Role      string         `json:"role" gorm:"default:user"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    Devices   []Device       `json:"devices,omitempty" gorm:"foreignKey:UserID"`
}