package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims represents JWT token claims.
type Claims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	CountryCode string   `json:"country_code"`
	FranchiseID string   `json:"franchise_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair holds access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// GenerateTokenPair creates a new access/refresh token pair.
func GenerateTokenPair(cfg JWTConfig, claims Claims) (TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(cfg.AccessTokenExpiry)

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(accessExp),
		Subject:   claims.UserID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := accessToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return TokenPair{}, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.RefreshTokenExpiry)),
		Subject:   claims.UserID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return TokenPair{}, fmt.Errorf("sign refresh token: %w", err)
	}

	return TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}

// ParseToken validates and parses a JWT access token.
func ParseToken(tokenStr string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
