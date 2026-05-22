package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func loadSecret() ([]byte, error) {
	secretString := os.Getenv("GOCART_JWT_SECRET")
	secretBytes := []byte(secretString)
	if len(secretBytes) < 64 {
		err := errors.New("JWT secret is too short or was not loaded correctly")
		log.Println(err)
		return nil, err
	}
	return secretBytes, nil
}

func loadIssuer() (string, error) {
	issuer := os.Getenv("GOCART_JWT_ISSUER")
	if issuer == "" {
		return "", errors.New("missing required env var: GOCART_JWT_ISSUER")
	}
	return issuer, nil
}

func GenerateJwtToken(userID uuid.UUID, purpose string, validForHours int) (string, error) {
	secret, err := loadSecret()
	if err != nil {
		return "", err
	}

	issuer, err := loadIssuer()
	if err != nil {
		return "", err
	}

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "iss": issuer,
        "sub": userID.String(),
        "iat": time.Now().Unix(),
        "exp": time.Now().Add(time.Duration(validForHours) * time.Hour).Unix(),
        "purpose": purpose,
    })

    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ValidateJwtToken(tokenString string, expectedPurpose string) (uuid.UUID, time.Time, error) {
	secret, err := loadSecret()
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

	issuer, err := loadIssuer()
	if err != nil {
		return uuid.Nil, time.Time{}, err
	}

    // Parse token and check validity
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
    
        return secret, nil
    })
    if err != nil || !token.Valid {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: Expired or invalid signature")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, time.Time{}, errors.New("Invalid token claims")
    }

    if iss, ok := claims["iss"].(string); !ok || iss != issuer {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: Issuer mismatch")
    }

    if purpose, ok := claims["purpose"].(string); !ok || purpose != expectedPurpose {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: Purpose mismatch")
    }

    // Extract "sub" (the user ID) to return it
    sub, ok := claims["sub"].(string)
    if !ok {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: Subject missing")
    }

    userID, err := uuid.Parse(sub)
    if err != nil {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: Subject is not a UUID")
    }

    // Extract "iat" (the issued at timestamp) to return it
    iatFloat, ok := claims["iat"].(float64) // JSON numbers become float64
    if !ok {
        return uuid.Nil, time.Time{}, errors.New("Invalid token: IssuedAt missing or not a number")
    }
    iat := time.Unix(int64(iatFloat), 0)

    return userID, iat, nil
}
