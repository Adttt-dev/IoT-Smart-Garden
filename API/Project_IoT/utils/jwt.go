package utils

import (
    "log"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID uint, username, role string) (string, error) {
    // Buat claims JWT
    claims := jwt.MapClaims{
        "user_id":  userID,  // Pastikan ini uint
        "username": username,
        "role":     role,
        "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 hari
        "iat":      time.Now().Unix(),
    }
    
    log.Printf("🔍 DEBUG: Membuat JWT dengan claims: %+v", claims)
    
    // Buat token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token dengan secret
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Println("❌ DEBUG: JWT_SECRET tidak ditemukan di environment")
        return "", jwt.ErrInvalidKey
    }
    
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        log.Printf("❌ DEBUG: Error signing token: %v", err)
        return "", err
    }
    
    log.Printf("✅ DEBUG: Token berhasil dibuat untuk user ID: %d", userID)
    return tokenString, nil
}

// Fungsi untuk memverifikasi token (opsional, untuk debug)
func VerifyToken(tokenString string) (*jwt.Token, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, jwt.ErrInvalidKey
    }
    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(secret), nil
    })
    
    return token, err
}