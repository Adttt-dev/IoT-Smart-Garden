// // controllers/device_controller.go
// package controllers

// import (
//     "net/http"
//     "strconv"
    
//     "github.com/gin-gonic/gin"
//     "project_iot/config"
//     "project_iot/models"
// )

// type DeviceRequest struct {
//     DeviceName string `json:"device_name" binding:"required"`
//     DeviceType string `json:"device_type"`
//     Location   string `json:"location"`
//     IPAddress  string `json:"ip_address"`
// }

// // CreateDevice - Membuat device baru (tetap milik user yang membuat)
// func CreateDevice(c *gin.Context) {
//     var req DeviceRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     userID := c.GetUint("user_id")
    
//     device := models.Device{
//         DeviceName: req.DeviceName,
//         DeviceType: req.DeviceType,
//         Location:   req.Location,
//         IPAddress:  req.IPAddress,
//         UserID:     userID,
//         IsActive:   true,
//     }
    
//     if err := config.DB.Create(&device).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
//         return
//     }
    
//     c.JSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
// }

// // GetDevices - Menampilkan semua device dari semua user (public access)
// func GetDevices(c *gin.Context) {
//     var devices []models.Device
    
//     // Mengambil semua device tanpa filter user_id, termasuk data user
//     if err := config.DB.Preload("User").Find(&devices).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{"devices": devices})
// }

// // GetDevice - Menampilkan detail device berdasarkan ID (public access)
// func GetDevice(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }
    
//     var device models.Device
    
//     // Mengambil device berdasarkan ID saja, tanpa filter user_id
//     if err := config.DB.Where("id = ?", deviceID).Preload("User").Preload("SensorData").First(&device).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{"device": device})
// }

// // UpdateDevice - Hanya owner device yang bisa update
// func UpdateDevice(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }
    
//     userID := c.GetUint("user_id")
    
//     var device models.Device
//     // Tetap mempertahankan filter user_id untuk update (hanya owner yang bisa update)
//     if err := config.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Device not found or you don't have permission to update this device"})
//         return
//     }
    
//     var req DeviceRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     device.DeviceName = req.DeviceName
//     device.DeviceType = req.DeviceType
//     device.Location = req.Location
//     device.IPAddress = req.IPAddress
    
//     if err := config.DB.Save(&device).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully", "device": device})
// }

// // DeleteDevice - Hanya owner device yang bisa delete
// func DeleteDevice(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }
    
//     userID := c.GetUint("user_id")
    
//     // Tetap mempertahankan filter user_id untuk delete (hanya owner yang bisa delete)
//     if err := config.DB.Where("id = ? AND user_id = ?", deviceID, userID).Delete(&models.Device{}).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device or you don't have permission to delete this device"})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
// }

// // GetMyDevices - Menampilkan device milik user yang sedang login
// func GetMyDevices(c *gin.Context) {
//     userID := c.GetUint("user_id")
    
//     var devices []models.Device
//     if err := config.DB.Where("user_id = ?", userID).Preload("User").Find(&devices).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch your devices"})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{"devices": devices})
// }

// // ToggleDeviceStatus - Hanya owner device yang bisa mengubah status
// func ToggleDeviceStatus(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }
    
//     userID := c.GetUint("user_id")
    
//     var device models.Device
//     // Hanya owner yang bisa mengubah status device
//     if err := config.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Device not found or you don't have permission to modify this device"})
//         return
//     }
    
//     // Toggle status
//     device.IsActive = !device.IsActive
    
//     if err := config.DB.Save(&device).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device status"})
//         return
//     }
        
//     status := "deactivated"
//     if device.IsActive {
//         status = "activated"
//     }
    
//     c.JSON(http.StatusOK, gin.H{
//         "message": "Device " + status + " successfully", 
//         "device": device,
//     })
// }

// controllers/device_controller.go
// package controllers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"project_iot/config"
// 	"project_iot/models"
// )

// type DeviceRequest struct {
// 	DeviceName string `json:"device_name" binding:"required"`
// 	DeviceType string `json:"device_type"`
// 	Location   string `json:"location"`
// 	IPAddress  string `json:"ip_address"`
// }

// // --- STRUCT BARU UNTUK MENANGANI PERINTAH ---
// type CommandRequest struct {
// 	Command string `json:"command" binding:"required,oneof=PUMP_ON PUMP_OFF AUTO_ON"`
// }

// // CreateDevice - Membuat device baru
// func CreateDevice(c *gin.Context) {
// 	var req DeviceRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userID := c.GetUint("user_id")

// 	device := models.Device{
// 		DeviceName: req.DeviceName,
// 		DeviceType: req.DeviceType,
// 		Location:   req.Location,
// 		IPAddress:  req.IPAddress,
// 		UserID:     userID,
// 		IsActive:   true,
// 		AutoMode:   true, // Default mode is AUTO
// 	}

// 	if err := config.DB.Create(&device).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
// }

// // GetDevices - Menampilkan semua device
// func GetDevices(c *gin.Context) {
// 	var devices []models.Device
// 	if err := config.DB.Preload("User").Find(&devices).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"devices": devices})
// }

// // GetDevice - Menampilkan detail device berdasarkan ID
// func GetDevice(c *gin.Context) {
// 	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
// 		return
// 	}

// 	var device models.Device
// 	if err := config.DB.Where("id = ?", deviceID).Preload("User").First(&device).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"device": device})
// }

// // UpdateDevice - Owner atau Admin bisa update
// func UpdateDevice(c *gin.Context) {
// 	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
// 		return
// 	}

// 	userID := c.GetUint("user_id")
// 	userRole := c.GetString("user_role")

// 	var device models.Device
// 	if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
// 		return
// 	}

// 	if device.UserID != userID && userRole != "admin" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this device"})
// 		return
// 	}

// 	var req DeviceRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	device.DeviceName = req.DeviceName
// 	device.DeviceType = req.DeviceType
// 	device.Location = req.Location
// 	device.IPAddress = req.IPAddress

// 	if err := config.DB.Save(&device).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully", "device": device})
// }

// // DeleteDevice - Owner atau Admin bisa delete
// func DeleteDevice(c *gin.Context) {
// 	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
// 		return
// 	}

// 	userID := c.GetUint("user_id")
// 	userRole := c.GetString("user_role")

// 	var device models.Device
// 	if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
// 		return
// 	}

// 	if device.UserID != userID && userRole != "admin" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this device"})
// 		return
// 	}

// 	if err := config.DB.Delete(&device).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
// }


// // --- FUNGSI BARU UNTUK KONTROL POMPA MANUAL ---
// func DeviceCommand(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }

//     userID := c.GetUint("user_id")
//     userRole := c.GetString("user_role")

//     var device models.Device
//     if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
//         return
//     }

//     if device.UserID != userID && userRole != "admin" {
//         c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied to send command to this device"})
//         return
//     }

//     var req CommandRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     switch req.Command {
//     case "PUMP_ON", "PUMP_OFF":
//         device.AutoMode = false // Manual mode engaged
//         device.LastCommand = req.Command
//     case "AUTO_ON":
//         device.AutoMode = true
//         device.LastCommand = req.Command // ESP will turn off the pump and switch to auto
//     default:
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command"})
//         return
//     }

//     if err := config.DB.Save(&device).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save command to device"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Command sent successfully", "device": device})
// }

// // --- FUNGSI BARU UNTUK KONFIRMASI DARI ESP32 ---
// func AcknowledgeCommand(c *gin.Context) {
//     deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
//         return
//     }

// 	// Any authenticated user (assumed to be the device owner via JWT) can ack
//     var device models.Device
//     if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
//         return
//     }

//     device.LastCommand = "" // Hapus perintah setelah dieksekusi oleh ESP32

//     if err := config.DB.Save(&device).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledge command"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "Command acknowledged"})
// }


// controllers/device_controller.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"project_iot/config"
	"project_iot/models"
)

type DeviceRequest struct {
	DeviceName string `json:"device_name" binding:"required"`
	DeviceType string `json:"device_type"`
	Location   string `json:"location"`
	IPAddress  string `json:"ip_address"`
}

// --- STRUCT BARU UNTUK MENANGANI PERINTAH ---
type CommandRequest struct {
	Command string `json:"command" binding:"required,oneof=PUMP_ON PUMP_OFF AUTO_ON"`
}

// CreateDevice - Membuat device baru
func CreateDevice(c *gin.Context) {
	var req DeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	device := models.Device{
		DeviceName: req.DeviceName,
		DeviceType: req.DeviceType,
		Location:   req.Location,
		IPAddress:  req.IPAddress,
		UserID:     userID,
		IsActive:   true,
		AutoMode:   true, // Default mode is AUTO
	}

	if err := config.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Device created successfully", "device": device})
}

// GetDevices - Menampilkan semua device
func GetDevices(c *gin.Context) {
	var devices []models.Device
	if err := config.DB.Preload("User").Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

// GetDevice - Menampilkan detail device berdasarkan ID
func GetDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var device models.Device
	if err := config.DB.Where("id = ?", deviceID).Preload("User").First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"device": device})
}

// UpdateDevice - Owner atau Admin bisa update
func UpdateDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var device models.Device
	if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if device.UserID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this device"})
		return
	}

	var req DeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device.DeviceName = req.DeviceName
	device.DeviceType = req.DeviceType
	device.Location = req.Location
	device.IPAddress = req.IPAddress

	if err := config.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully", "device": device})
}

// DeleteDevice - Owner atau Admin bisa delete
func DeleteDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var device models.Device
	if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if device.UserID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this device"})
		return
	}

	if err := config.DB.Delete(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}


// --- FUNGSI BARU UNTUK KONTROL POMPA MANUAL ---
func DeviceCommand(c *gin.Context) {
    deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
        return
    }

    userID := c.GetUint("user_id")
    userRole := c.GetString("user_role")

    var device models.Device
    if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
        return
    }

    if device.UserID != userID && userRole != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied to send command to this device"})
        return
    }

    var req CommandRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    switch req.Command {
    case "PUMP_ON", "PUMP_OFF":
        device.AutoMode = false // Manual mode engaged
        device.LastCommand = req.Command
    case "AUTO_ON":
        device.AutoMode = true
        device.LastCommand = req.Command // ESP will turn off the pump and switch to auto
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command"})
        return
    }

    if err := config.DB.Save(&device).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save command to device"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Command sent successfully", "device": device})
}

// --- FUNGSI BARU UNTUK KONFIRMASI DARI ESP32 ---
func AcknowledgeCommand(c *gin.Context) {
    deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
        return
    }

	// Any authenticated user (assumed to be the device owner via JWT) can ack
    var device models.Device
    if err := config.DB.Where("id = ?", deviceID).First(&device).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
        return
    }

    device.LastCommand = "" // Hapus perintah setelah dieksekusi oleh ESP32

    if err := config.DB.Save(&device).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledge command"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Command acknowledged"})
}