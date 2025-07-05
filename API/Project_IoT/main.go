package main

import (
	"log"
	"os"
	"time" // Diperlukan untuk penjadwal

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"project_iot/config"
	"project_iot/controllers"
	"project_iot/middleware"
	"project_iot/models" // Diperlukan untuk mengakses model SensorData
)

// --- FUNGSI BACKGROUND JOB UNTUK PEMBERSIHAN DATA ---

// Fungsi untuk membersihkan data lama di database
func cleanupOldSensorData() {
	// Batas maksimal baris yang ingin Anda simpan
	const maxRows = 50

	log.Println("SCHEDULER: Running periodic data cleanup check...")

	// Hitung jumlah baris saat ini
	var currentCount int64
	if err := config.DB.Model(&models.SensorData{}).Count(&currentCount).Error; err != nil {
		log.Printf("SCHEDULER ERROR: Could not count sensor data: %v", err)
		return
	}

	log.Printf("SCHEDULER INFO: Current sensor data rows: %d", currentCount)

	// Jika jumlah baris melebihi batas, hapus yang terlama
	if currentCount > maxRows {
		rowsToDelete := currentCount - maxRows
		log.Printf("SCHEDULER INFO: Row count exceeds limit. Deleting %d oldest records...", rowsToDelete)

		// Hapus 'N' baris data yang PALING LAMA
		result := config.DB.Exec(
			"DELETE FROM sensor_readings ORDER BY server_timestamp ASC LIMIT ?",
			rowsToDelete,
		)

		if result.Error != nil {
			log.Printf("SCHEDULER ERROR: Failed to delete old records: %v", result.Error)
			return
		}

		log.Printf("SCHEDULER SUCCESS: Deleted %d old records.", result.RowsAffected)
	} else {
		log.Println("SCHEDULER INFO: Row count is within the limit. No action needed.")
	}
}

// Fungsi untuk membuat dan menjalankan penjadwal (ticker)
func runCleanupScheduler() {
	// --- PERUBAHAN UTAMA DI SINI ---
	// Jadwal diatur untuk berjalan setiap 10 menit
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Loop selamanya untuk menjalankan pembersihan sesuai jadwal
	for range ticker.C {
		cleanupOldSensorData()
	}
}

// --- AKHIR FUNGSI BACKGROUND JOB ---


func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	config.ConnectDatabase()

	// Gin router
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "Smart Irrigation API is running"})
	})

	// Public (tanpa login, untuk ESP32)
	public := r.Group("/api/public")
	{
		public.POST("/sensor-data", controllers.CreateSensorDataPublic)
		public.POST("/sensor-readings", controllers.CreateSensorDataPublic)
	}

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes (dengan JWT)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Rute-rute Anda yang lain...
		api.GET("/users", controllers.GetAllUsers)
		api.DELETE("/users", controllers.DeleteAllUsers)
		api.DELETE("/users/:id", controllers.DeleteUserByID)
		api.GET("/profile", controllers.GetProfile)
		api.POST("/devices", controllers.CreateDevice)
		api.GET("/devices", controllers.GetDevices)
		api.GET("/devices/:id", controllers.GetDevice)
		api.PUT("/devices/:id", controllers.UpdateDevice)
		api.DELETE("/devices/:id", controllers.DeleteDevice)
		api.PUT("/devices/:id/command", controllers.DeviceCommand)
		api.PUT("/devices/:id/command-ack", controllers.AcknowledgeCommand)
		api.POST("/sensor-data", controllers.CreateSensorData)
		api.GET("/sensor-data/device/:device_id", controllers.GetSensorData)
		api.GET("/sensor-data/device/:device_id/latest", controllers.GetLatestSensorData)
		api.POST("/sensor-readings", controllers.CreateSensorData)
		api.GET("/sensor-readings/device/:device_id", controllers.GetSensorData)
		api.GET("/sensor-readings/device/:device_id/latest", controllers.GetLatestSensorData)
	}

	// Jalankan penjadwal pembersihan data di latar belakang
	go runCleanupScheduler()

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	// --- PERUBAHAN PESAN LOG ---
	log.Printf("🧹 Data cleanup scheduler is running in the background (check runs every 10 minutes).")
	log.Printf("📡 ESP32 endpoints (public):")
	log.Printf("   - Sensor data: http://localhost:%s/api/public/sensor-data", port)
	log.Printf("   - Device settings: http://localhost:%s/api/public/device-settings/{device_id}", port)
	log.Printf("🔐 Protected endpoints:")
	log.Printf("   - User management: http://localhost:%s/api/users", port)
	log.Printf("   - Profile: http://localhost:%s/api/profile", port)
	log.Printf("   - Devices: http://localhost:%s/api/devices", port)
	log.Printf("   - Sensor data: http://localhost:%s/api/sensor-data", port)
	log.Printf("   - Device Command: http://localhost:%s/api/devices/{id}/command", port)

	r.Run(":" + port)
}