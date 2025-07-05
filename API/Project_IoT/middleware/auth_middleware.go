package middleware

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            log.Println("❌ DEBUG: Header Authorization tidak ada")
            c.JSON(401, gin.H{"error": "Header Authorization diperlukan"})
            c.Abort()
            return
        }

        // Ekstrak token dari "Bearer <token>"
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            log.Println("❌ DEBUG: Format token tidak valid (harus Bearer <token>)")
            c.JSON(401, gin.H{"error": "Format token tidak valid"})
            c.Abort()
            return
        }
        
        log.Printf("🔍 DEBUG: Token string: %s...", tokenString[:min(len(tokenString), 20)])

        // Parse dan validasi token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Pastikan signing method adalah HMAC
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("metode signing tidak dikenali: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil {
            log.Printf("❌ DEBUG: Error parse JWT: %v", err)
            c.JSON(401, gin.H{"error": "Token tidak valid"})
            c.Abort()
            return
        }

        // Ekstrak claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            log.Printf("🔍 DEBUG: JWT Claims: %+v", claims)
            
            // Ambil user_id dari claims
            userIDClaim, exists := claims["user_id"]
            if !exists {
                log.Println("❌ DEBUG: user_id tidak ada di JWT claims")
                c.JSON(401, gin.H{"error": "Token tidak valid - user_id hilang"})
                c.Abort()
                return
            }
            
            log.Printf("🔍 DEBUG: user_id dari claims: %v (tipe: %T)", userIDClaim, userIDClaim)
            
            // Konversi user_id ke uint
            var userID uint
            switch v := userIDClaim.(type) {
            case float64:
                userID = uint(v)
                log.Printf("🔍 DEBUG: Konversi float64 ke uint: %d", userID)
            case int:
                userID = uint(v)
                log.Printf("🔍 DEBUG: Konversi int ke uint: %d", userID)
            case uint:
                userID = v
                log.Printf("🔍 DEBUG: user_id sudah uint: %d", userID)
            default:
                log.Printf("❌ DEBUG: Tipe user_id tidak didukung: %T", v)
                c.JSON(500, gin.H{"error": "Tipe user_id tidak valid"})
                c.Abort()
                return
            }
            
            // Set user_id di context sebagai uint
            c.Set("user_id", userID)
            log.Printf("✅ DEBUG: Set user_id di context: %d", userID)
            
            // Optional: set claims lain jika diperlukan
            if username, ok := claims["username"].(string); ok {
                c.Set("username", username)
            }
            if role, ok := claims["role"].(string); ok {
                c.Set("role", role)
            }
            
        } else {
            log.Println("❌ DEBUG: Claims JWT tidak valid")
            c.JSON(401, gin.H{"error": "Token tidak valid"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Helper function untuk min
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}