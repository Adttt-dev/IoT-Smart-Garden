// models/sensor_data.go
package models

import (
    "time"
)

type SensorData struct {
    ID                    uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    DeviceID              uint           `json:"device_id" gorm:"not null;index"`
    Device                Device         `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
    Temperature           float64        `json:"temperature" gorm:"type:decimal(5,2)"`
    Humidity              float64        `json:"humidity" gorm:"type:decimal(5,2)"`
    TemperatureSource     string         `json:"temperature_source" gorm:"type:enum('sensor','cached');default:sensor"`
    HumiditySource        string         `json:"humidity_source" gorm:"type:enum('sensor','cached');default:sensor"`
    SoilMoistureRaw       int            `json:"soil_moisture_raw"`
    SoilMoisturePercent   float64        `json:"soil_moisture_percent" gorm:"type:decimal(5,2)"`
    WaterLevelCm          float64        `json:"water_level_cm" gorm:"type:decimal(6,2)"`
    WaterPercentage       float64        `json:"water_percentage" gorm:"type:decimal(5,2)"`
    TankHeightCm          float64        `json:"tank_height_cm" gorm:"type:decimal(6,2);default:100.00"`
    PumpStatus            string         `json:"pump_status" gorm:"type:enum('OFF','MED','HIGH','MAX','NO_WATER');not null"`
    PumpPwmValue          int            `json:"pump_pwm_value" gorm:"type:smallint;default:0"`
    PumpPercentage        int            `json:"pump_percentage" gorm:"type:tinyint;default:0"`
    SystemStatus          string         `json:"system_status" gorm:"type:varchar(50);not null"`
    LogicExplanation      string         `json:"logic_explanation" gorm:"type:text"`
    WifiRssi              int            `json:"wifi_rssi" gorm:"type:smallint"`
    FreeHeap              int            `json:"free_heap"`
    UptimeMs              int64          `json:"uptime_ms" gorm:"type:bigint"`
    DeviceTimestamp       *time.Time     `json:"device_timestamp"`
    ServerTimestamp       time.Time      `json:"server_timestamp" gorm:"default:CURRENT_TIMESTAMP"`
    
    // HAPUS field GORM standar karena tidak ada di tabel database
    // CreatedAt             time.Time      `json:"created_at"`
    // UpdatedAt             time.Time      `json:"updated_at"`
    // DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

// Method untuk menggunakan tabel yang benar
func (SensorData) TableName() string {
    return "sensor_readings"
}

// OPTIONAL: Jika Anda ingin menggunakan GORM timestamps tapi dengan nama kolom custom
// Uncomment bagian di bawah dan comment bagian di atas:

/*
type SensorData struct {
    ID                    uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    DeviceID              uint           `json:"device_id" gorm:"not null;index"`
    Device                Device         `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
    Temperature           float64        `json:"temperature" gorm:"type:decimal(5,2)"`
    Humidity              float64        `json:"humidity" gorm:"type:decimal(5,2)"`
    TemperatureSource     string         `json:"temperature_source" gorm:"type:enum('sensor','cached');default:sensor"`
    HumiditySource        string         `json:"humidity_source" gorm:"type:enum('sensor','cached');default:sensor"`
    SoilMoistureRaw       int            `json:"soil_moisture_raw"`
    SoilMoisturePercent   float64        `json:"soil_moisture_percent" gorm:"type:decimal(5,2)"`
    WaterLevelCm          float64        `json:"water_level_cm" gorm:"type:decimal(6,2)"`
    WaterPercentage       float64        `json:"water_percentage" gorm:"type:decimal(5,2)"`
    TankHeightCm          float64        `json:"tank_height_cm" gorm:"type:decimal(6,2);default:100.00"`
    PumpStatus            string         `json:"pump_status" gorm:"type:enum('OFF','MED','HIGH','MAX','NO_WATER');not null"`
    PumpPwmValue          int            `json:"pump_pwm_value" gorm:"type:smallint;default:0"`
    PumpPercentage        int            `json:"pump_percentage" gorm:"type:tinyint;default:0"`
    SystemStatus          string         `json:"system_status" gorm:"type:varchar(50);not null"`
    LogicExplanation      string         `json:"logic_explanation" gorm:"type:text"`
    WifiRssi              int            `json:"wifi_rssi" gorm:"type:smallint"`
    FreeHeap              int            `json:"free_heap"`
    UptimeMs              int64          `json:"uptime_ms" gorm:"type:bigint"`
    DeviceTimestamp       *time.Time     `json:"device_timestamp"`
    ServerTimestamp       time.Time      `json:"server_timestamp" gorm:"default:CURRENT_TIMESTAMP"`
}

// Disable GORM timestamps karena kita menggunakan server_timestamp
func (SensorData) BeforeCreate(tx *gorm.DB) error {
    return nil
}

func (SensorData) TableName() string {
    return "sensor_readings"
}
*/