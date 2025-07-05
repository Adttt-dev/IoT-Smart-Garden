// // models/device.go
// package models

// import (
//     "gorm.io/gorm"
//     "time"
// )

// type Device struct {
//     ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
//     DeviceName  string         `json:"device_name" gorm:"not null"`
//     DeviceType  string         `json:"device_type" gorm:"default:irrigation"`
//     Location    string         `json:"location"`
//     IsActive    bool           `json:"is_active" gorm:"default:true"`
//     IPAddress   string         `json:"ip_address"`  
//     UserID      uint           `json:"user_id" gorm:"not null;index"`
//     User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
//     CreatedAt   time.Time      `json:"created_at"`
//     UpdatedAt   time.Time      `json:"updated_at"`
//     DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
//     SensorData  []SensorData   `json:"sensor_data,omitempty" gorm:"foreignKey:DeviceID"`
// }


// // models/device.go
// package models

// import (
// 	"gorm.io/gorm"
// 	"time"
// )

// type Device struct {
// 	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
// 	DeviceName  string `json:"device_name" gorm:"not null"`
// 	DeviceType  string `json:"device_type" gorm:"default:irrigation"`
// 	Location    string `json:"location"`
// 	IsActive    bool   `json:"is_active" gorm:"default:true"`
// 	IPAddress   string `json:"ip_address"`
// 	// --- FIELD YANG DITAMBAHKAN UNTUK KONTROL MANUAL ---
// 	AutoMode    bool   `json:"auto_mode" gorm:"default:true"`
// 	LastCommand string `json:"last_command,omitempty" gorm:"type:varchar(50);default:null"`
// 	// --- AKHIR FIELD YANG DITAMBAHKAN ---
// 	UserID     uint         `json:"user_id" gorm:"not null;index"`
// 	User       User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
// 	CreatedAt  time.Time    `json:"created_at"`
// 	UpdatedAt  time.Time    `json:"updated_at"`
// 	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
// 	SensorData []SensorData `json:"sensor_data,omitempty" gorm:"foreignKey:DeviceID"`
// }

// models/device.go
package models

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceName  string `json:"device_name" gorm:"not null"`
	DeviceType  string `json:"device_type" gorm:"default:irrigation"`
	Location    string `json:"location"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	IPAddress   string `json:"ip_address"`
	// --- FIELD YANG DITAMBAHKAN UNTUK KONTROL MANUAL ---
	AutoMode    bool   `json:"auto_mode" gorm:"default:true"`
	LastCommand string `json:"last_command,omitempty" gorm:"type:varchar(50);default:null"`
	// --- AKHIR FIELD YANG DITAMBAHKAN ---
	UserID     uint         `json:"user_id" gorm:"not null;index"`
	User       User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	SensorData []SensorData `json:"sensor_data,omitempty" gorm:"foreignKey:DeviceID"`
}