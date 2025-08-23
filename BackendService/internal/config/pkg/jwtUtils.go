package pkg

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/types"
)

type AppClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type MQTTClaims struct {
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	ACL      []types.ACL `json:"acl"`
	jwt.RegisteredClaims
}

// GenerateJWKToken generates a JWT token with the given kid, alg, and key
func GenerateJWKToken(user *types.User, acl *[]types.ACL, path string, ttl time.Duration) (string, error) {
	var customClaims jwt.Claims
	if acl == nil {
		customClaims = &AppClaims{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.Username,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
				Issuer:    "app",
				Audience:  jwt.ClaimStrings{"app"},
			},
		}
	} else {
		customClaims = &MQTTClaims{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			ACL:      *acl,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.Username,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
				Issuer:    "mqtt",
				Audience:  jwt.ClaimStrings{"mqtt"},
			},
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, customClaims)

	token.Header["kid"] = "ec256-2025-01"

	priv, err := LoadECPrivateKey(path)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(priv)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWKToken(tokenString string, path string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodES256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		pub, _ := LoadECPublicKey(path)
		return pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func LoadECPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func LoadECPublicKey(path string) (*ecdsa.PublicKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if pubKey, ok := pubKey.(*ecdsa.PublicKey); ok {
		return pubKey, nil
	}
	return nil, fmt.Errorf("invalid public key")
}
