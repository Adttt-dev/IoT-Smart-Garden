package controllers

import (
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "project_iot/config"
    "project_iot/models"
)

type SensorDataRequest struct {
    DeviceID                uint       `json:"device_id" binding:"required"`
    Temperature             float64    `json:"temperature"`
    Humidity                float64    `json:"humidity"`
    TemperatureSource       string     `json:"temperature_source"`
    HumiditySource          string     `json:"humidity_source"`
    SoilMoistureRaw         int        `json:"soil_moisture_raw"`
    SoilMoisturePercent     float64    `json:"soil_moisture_percent"`
    WaterLevelCm            float64    `json:"water_level_cm"`
    WaterPercentage         float64    `json:"water_percentage"`
    TankHeightCm            float64    `json:"tank_height_cm"`
    PumpStatus              string     `json:"pump_status" binding:"required"`
    PumpPwmValue            int        `json:"pump_pwm_value"`
    PumpPercentage          int        `json:"pump_percentage"`
    SystemStatus            string     `json:"system_status" binding:"required"`
    LogicExplanation        string     `json:"logic_explanation"`
    WifiRssi                int        `json:"wifi_rssi"`
    FreeHeap                int        `json:"free_heap"`
    UptimeMs                int64      `json:"uptime_ms"`
    // UBAH: Gunakan int64 untuk menerima milliseconds dari ESP32
    DeviceTimestampMs       int64      `json:"device_timestamp"`
}

// Endpoint untuk ESP32 (tanpa autentikasi) - FIXED VERSION
func CreateSensorDataPublic(c *gin.Context) {
    var req SensorDataRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("❌ JSON binding error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }
    
    log.Printf("📥 Received data from ESP32: DeviceID=%d, Temp=%.1f, Humidity=%.1f, SoilRaw=%d", 
        req.DeviceID, req.Temperature, req.Humidity, req.SoilMoistureRaw)
    
    // DEBUG: Check database connection first
    if config.DB == nil {
        log.Printf("❌ Database connection is nil")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
        return
    }
    
    // DEBUG: Test database connectivity
    sqlDB, err := config.DB.DB()
    if err != nil {
        log.Printf("❌ Failed to get underlying sql.DB: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
        return
    }
    
    if err := sqlDB.Ping(); err != nil {
        log.Printf("❌ Database ping failed: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
        return
    }
    
    log.Printf("✅ Database connection is healthy")
    
    // Cek apakah device exists (tanpa cek user_id)
    var device models.Device
    if err := config.DB.Where("id = ?", req.DeviceID).First(&device).Error; err != nil {
        log.Printf("❌ Device not found: DeviceID=%d, Error=%v", req.DeviceID, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
        return
    }
    
    // Use correct field name based on your Device model
    log.Printf("✅ Device found: %s (ID: %d)", device.DeviceName, device.ID)
    
    // Validate required sensor data
    if req.DeviceID == 0 {
        log.Printf("❌ Invalid DeviceID: %d", req.DeviceID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "DeviceID is required and must be greater than 0"})
        return
    }
    
    // PERBAIKAN: Konversi device timestamp dari milliseconds ke time.Time
    var deviceTimestamp *time.Time
    if req.DeviceTimestampMs > 0 {
        // Konversi milliseconds ke time.Time
        // Jika timestamp dalam format Unix milliseconds (epoch)
        if req.DeviceTimestampMs > 1000000000000 { // > year 2001 in milliseconds
            ts := time.Unix(0, req.DeviceTimestampMs*int64(time.Millisecond))
            deviceTimestamp = &ts
        } else {
            // Jika timestamp adalah uptime dalam milliseconds, konversi relatif ke waktu sekarang
            // Alternatif: bisa juga diabaikan dan menggunakan server timestamp
            ts := time.Now().Add(-time.Duration(req.DeviceTimestampMs) * time.Millisecond)
            deviceTimestamp = &ts
        }
    }
    
    // Buat sensor data dengan field yang sesuai model
    sensorData := models.SensorData{
        DeviceID:              req.DeviceID,
        Temperature:           req.Temperature,
        Humidity:              req.Humidity,
        TemperatureSource:     req.TemperatureSource,
        HumiditySource:        req.HumiditySource,
        SoilMoistureRaw:       req.SoilMoistureRaw,
        SoilMoisturePercent:   req.SoilMoisturePercent,
        WaterLevelCm:          req.WaterLevelCm,
        WaterPercentage:       req.WaterPercentage,
        TankHeightCm:          req.TankHeightCm,
        PumpStatus:            req.PumpStatus,
        PumpPwmValue:          req.PumpPwmValue,
        PumpPercentage:        req.PumpPercentage,
        SystemStatus:          req.SystemStatus,
        LogicExplanation:      req.LogicExplanation,
        WifiRssi:              req.WifiRssi,
        FreeHeap:              req.FreeHeap,
        UptimeMs:              req.UptimeMs,
        DeviceTimestamp:       deviceTimestamp, // Gunakan pointer yang sudah dikonversi
        ServerTimestamp:       time.Now(),
    }
    
    // DEBUG: Print data yang akan disimpan
    log.Printf("📊 Saving sensor data: %+v", sensorData)
    
    // DEBUG: Test query sederhana dengan error handling yang lebih baik
    var count int64
    // Use the correct table name from model
    if err := config.DB.Model(&models.SensorData{}).Count(&count).Error; err != nil {
        log.Printf("❌ Database table access error: %v", err)
        // Check if table exists using the correct table name
        tableName := (&models.SensorData{}).TableName()
        if err := config.DB.Exec("SELECT 1 FROM " + tableName + " LIMIT 1").Error; err != nil {
            log.Printf("❌ %s table doesn't exist or is not accessible: %v", tableName, err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Database table '" + tableName + "' not found or not accessible", 
                "details": err.Error(),
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Database table access error", 
            "details": err.Error(),
        })
        return
    }
    
    log.Printf("📊 Current sensor data records count: %d", count)
    
    // DEBUG: Check if DeviceID exists in foreign key constraint
    var deviceExists bool
    if err := config.DB.Model(&models.Device{}).
        Select("count(*) > 0").
        Where("id = ?", req.DeviceID).
        Find(&deviceExists).Error; err != nil {
        log.Printf("❌ Error checking device existence: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Database error while validating device", 
            "details": err.Error(),
        })
        return
    }
    
    if !deviceExists {
        log.Printf("❌ Device with ID %d does not exist for foreign key constraint", req.DeviceID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device_id: device does not exist"})
        return
    }
    
    // Simpan data dengan transaction untuk safety
    tx := config.DB.Begin()
    if tx.Error != nil {
        log.Printf("❌ Failed to begin transaction: %v", tx.Error)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database transaction error"})
        return
    }
    
    if err := tx.Create(&sensorData).Error; err != nil {
        tx.Rollback()
        log.Printf("❌ Database save error: %v", err)
        log.Printf("❌ Failed data: %+v", sensorData)
        
        // Provide more specific error messages
        errorMsg := "Failed to save sensor data"
        errStr := strings.ToLower(err.Error())
        if strings.Contains(errStr, "foreign key constraint") {
            errorMsg = "Invalid device_id: device not found"
        } else if strings.Contains(errStr, "duplicate key") {
            errorMsg = "Duplicate sensor data entry"
        } else if strings.Contains(errStr, "column") && strings.Contains(errStr, "does not exist") {
            errorMsg = "Database schema mismatch - please check sensor_readings table structure"
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": errorMsg, 
            "details": err.Error(),
        })
        return
    }
    
    if err := tx.Commit().Error; err != nil {
        log.Printf("❌ Failed to commit transaction: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database commit error"})
        return
    }
    
    log.Printf("✅ Sensor data saved successfully for device %d with ID %d", req.DeviceID, sensorData.ID)
    c.JSON(http.StatusCreated, gin.H{
        "message": "Sensor data saved successfully", 
        "data": sensorData,
    })
}

// Endpoint untuk web app (dengan autentikasi) - FIXED VERSION
func CreateSensorData(c *gin.Context) {
    var req SensorDataRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("❌ JSON binding error in authenticated endpoint: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }
    
    // Verify device belongs to user
    userID := c.GetUint("user_id")
    if userID == 0 {
        log.Printf("❌ Invalid user_id from context: %d", userID)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    var device models.Device
    if err := config.DB.Where("id = ? AND user_id = ?", req.DeviceID, userID).First(&device).Error; err != nil {
        log.Printf("❌ Device not found or not owned by user: DeviceID=%d, UserID=%d, Error=%v", 
            req.DeviceID, userID, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Device not found or access denied"})
        return
    }
    
    // PERBAIKAN: Sama seperti di endpoint public
    var deviceTimestamp *time.Time
    if req.DeviceTimestampMs > 0 {
        if req.DeviceTimestampMs > 1000000000000 { // > year 2001 in milliseconds
            ts := time.Unix(0, req.DeviceTimestampMs*int64(time.Millisecond))
            deviceTimestamp = &ts
        } else {
            ts := time.Now().Add(-time.Duration(req.DeviceTimestampMs) * time.Millisecond)
            deviceTimestamp = &ts
        }
    }
    
    // Use the same field mapping as the public endpoint
    sensorData := models.SensorData{
        DeviceID:              req.DeviceID,
        Temperature:           req.Temperature,
        Humidity:              req.Humidity,
        TemperatureSource:     req.TemperatureSource,
        HumiditySource:        req.HumiditySource,
        SoilMoistureRaw:       req.SoilMoistureRaw,
        SoilMoisturePercent:   req.SoilMoisturePercent,
        WaterLevelCm:          req.WaterLevelCm,
        WaterPercentage:       req.WaterPercentage,
        TankHeightCm:          req.TankHeightCm,
        PumpStatus:            req.PumpStatus,
        PumpPwmValue:          req.PumpPwmValue,
        PumpPercentage:        req.PumpPercentage,
        SystemStatus:          req.SystemStatus,
        LogicExplanation:      req.LogicExplanation,
        WifiRssi:              req.WifiRssi,
        FreeHeap:              req.FreeHeap,
        UptimeMs:              req.UptimeMs,
        DeviceTimestamp:       deviceTimestamp,
        ServerTimestamp:       time.Now(),
    }
    
    if err := config.DB.Create(&sensorData).Error; err != nil {
        log.Printf("❌ Failed to save sensor data for authenticated user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save sensor data"})
        return
    }
    
    log.Printf("✅ Sensor data saved successfully for authenticated user %d, device %d", userID, req.DeviceID)
    c.JSON(http.StatusCreated, gin.H{"message": "Sensor data saved successfully", "data": sensorData})
}

// UBAH: Fungsi GetSensorData menjadi public
func GetSensorData(c *gin.Context) {
    deviceID, err := strconv.ParseUint(c.Param("device_id"), 10, 32)
    if err != nil {
        log.Printf("❌ Parameter device_id tidak valid: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID tidak valid"})
        return
    }
    
    userID := c.GetUint("user_id")
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
        return
    }
    
    // UBAH: Hilangkan filter user_id - semua user bisa akses device apapun
    var device models.Device
    if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
        log.Printf("❌ Device tidak ditemukan: DeviceID=%d", deviceID)
        c.JSON(http.StatusNotFound, gin.H{"error": "Device tidak ditemukan"})
        return
    }
    
    // Ambil parameter pagination dengan validasi
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    if page < 1 {
        page = 1
    }
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
    if limit < 1 || limit > 1000 {
        limit = 50
    }
    
    offset := (page - 1) * limit
    
    var sensorData []models.SensorData
    if err := config.DB.Preload("Device").
        Where("device_id = ?", deviceID).
        Order("server_timestamp DESC").
        Limit(limit).
        Offset(offset).
        Find(&sensorData).Error; err != nil {
        log.Printf("❌ Gagal mengambil data sensor: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data sensor"})
        return
    }
    
    // Ambil total count untuk info pagination
    var totalCount int64
    config.DB.Model(&models.SensorData{}).Where("device_id = ?", deviceID).Count(&totalCount)
    
    log.Printf("✅ Data sensor berhasil diambil: DeviceID=%d, UserID=%d, Count=%d", 
        deviceID, userID, len(sensorData))
    
    c.JSON(http.StatusOK, gin.H{
        "data":        sensorData, 
        "page":        page, 
        "limit":       limit,
        "total":       totalCount,
        "total_pages": (totalCount + int64(limit) - 1) / int64(limit),
    })
}

// UBAH: Fungsi GetLatestSensorData menjadi public
func GetLatestSensorData(c *gin.Context) {
    deviceID, err := strconv.ParseUint(c.Param("device_id"), 10, 32)
    if err != nil {
        log.Printf("❌ Parameter device_id tidak valid: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID tidak valid"})
        return
    }
    
    userID := c.GetUint("user_id")
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
        return
    }
    
    // UBAH: Hilangkan filter user_id - semua user bisa akses device apapun
    var device models.Device
    if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
        log.Printf("❌ Device tidak ditemukan: DeviceID=%d", deviceID)
        c.JSON(http.StatusNotFound, gin.H{"error": "Device tidak ditemukan"})
        return
    }
    
    var sensorData models.SensorData
    if err := config.DB.Preload("Device").
        Where("device_id = ?", deviceID).
        Order("server_timestamp DESC").
        First(&sensorData).Error; err != nil {
        log.Printf("❌ Tidak ada data sensor untuk device %d", deviceID)
        c.JSON(http.StatusNotFound, gin.H{"error": "Data sensor tidak ditemukan"})
        return
    }
    
    log.Printf("✅ Data sensor terbaru berhasil diambil: DeviceID=%d, UserID=%d", 
        deviceID, userID)
    
    c.JSON(http.StatusOK, gin.H{"data": sensorData})
}


// TAMBAHAN: Fungsi untuk mendapatkan semua devices (public)
func GetAllDevices(c *gin.Context) {
    userID := c.GetUint("user_id")
    if userID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak terautentikasi"})
        return
    }
    
    var devices []models.Device
    if err := config.DB.Find(&devices).Error; err != nil {
        log.Printf("❌ Gagal mengambil daftar devices: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar devices"})
        return
    }
    
    // Format response dengan informasi basic device
    var responseDevices []gin.H
    for _, device := range devices {
        responseDevices = append(responseDevices, gin.H{
            "id":         device.ID,
            "name":       device.DeviceName,
            "user_id":    device.UserID,
            "created_at": device.CreatedAt,
        })
    }
    
    log.Printf("✅ Daftar devices berhasil diambil: UserID=%d, Count=%d", 
        userID, len(devices))
    
    c.JSON(http.StatusOK, gin.H{
        "devices": responseDevices,
        "total":   len(devices),
    })
}

// PERBAIKAN: DebugDeviceRelation juga dibuat public
func DebugDeviceRelation(c *gin.Context) {
    deviceID, err := strconv.ParseUint(c.Param("device_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID tidak valid"})
        return
    }
    
    // Cek device ada atau tidak (tanpa filter user_id)
    var device models.Device
    if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Device tidak ditemukan",
            "device_id": deviceID,
        })
        return
    }
    
    // Cek sensor data ada atau tidak
    var sensorCount int64
    config.DB.Model(&models.SensorData{}).Where("device_id = ?", deviceID).Count(&sensorCount)
    
    // Cek foreign key constraint
    var foreignKeyExists bool
    config.DB.Raw(`
        SELECT EXISTS(
            SELECT 1 FROM information_schema.table_constraints 
            WHERE constraint_name LIKE '%sensor_readings%device_id%' 
            AND table_name = 'sensor_readings'
        )
    `).Scan(&foreignKeyExists)
    
    c.JSON(http.StatusOK, gin.H{
        "device_found":         true,
        "device_id":           device.ID,
        "device_name":         device.DeviceName,
        "device_user_id":      device.UserID,
        "sensor_data_count":   sensorCount,
        "foreign_key_exists":  foreignKeyExists,
        "message":             "Debug info untuk relasi Device-SensorData (PUBLIC ACCESS)",
    })
}
