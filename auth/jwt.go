// Updated jwt.go
package auth

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

// JWT token errors
var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token has expired")
)

// Default JWT expiration time if not specified
const DefaultTokenExpiration = 24 * time.Hour

// JWT claims struct
type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Get JWT secret key from environment variable with fallback
func getJWTSecret() string {
    secretKey := os.Getenv("JWT_SECRET")
    if secretKey == "" {
        // Fallback to default only if not set
        secretKey = "your-default-secret-key-for-development-only"
    }
    return secretKey
}

// Generate a new JWT token
func GenerateToken(userID int, email string, expiration ...time.Duration) (string, error) {
    // Use provided expiration or default
    tokenExpiration := DefaultTokenExpiration
    if len(expiration) > 0 {
        tokenExpiration = expiration[0]
    }

    // Get secret key from environment variable
    secretKey := getJWTSecret()

    // Create token claims
    claims := &Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    // Create token with claims and sign it
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// Validate JWT token and extract claims
func ValidateToken(tokenString string) (*Claims, error) {
    // Get secret key from environment variable
    secretKey := getJWTSecret()

    // Parse and validate token
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, ErrInvalidToken
            }
            return []byte(secretKey), nil
        },
    )

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

    // Extract claims
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}