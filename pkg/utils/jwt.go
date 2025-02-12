package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = os.Getenv("JWT_SECRET_KEY")
var refreshSecretKey = os.Getenv("JWT_REFRESH_KEY")

func GenerateSessionToken(sessionId string) (*models.JwtToken, error) {
	var jwtToken models.JwtToken

	hoursCount, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	claims["session_id"] = sessionId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	jwtToken.AccessToken = accessTokenString

	hoursCount, err = strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRE_HOURS"))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"session_id": sessionId,
		"exp":        time.Now().Add(time.Hour * 24 * time.Duration(hoursCount)).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecretKey))
	if err != nil {
		return nil, err
	}

	jwtToken.RefreshToken = refreshTokenString

	return &jwtToken, nil
}

func VerifyToken(tokenString string, isRefreshToken bool) (*jwt.Token, error) {

	key := secretKey
	if isRefreshToken {
		key = refreshSecretKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else {
				return nil, errors.New("token is invalid")
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}

func ExtractClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok && token.Valid
}
