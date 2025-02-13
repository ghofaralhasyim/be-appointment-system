package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ghofaralhasyim/be-appointment-system/internal/repositories"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(redisRepo repositories.RedisRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			}

			tokenString := parts[1]
			token, err := utils.VerifyToken(tokenString, false)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			claims, ok := utils.ExtractClaims(token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Failed to extract claims from token"})
			}

			sessionId, ok := claims["session_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			dataJSON, err := redisRepo.Get(context.Background(), sessionId)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired session"})
			}

			var dataUser map[string]interface{}
			if err := json.Unmarshal([]byte(dataJSON), &dataUser); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed get session data"})
			}

			userId, ok := dataUser["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
			}

			c.Set("userId", int(userId))
			c.Set("sessionId", sessionId)

			return next(c)
		}
	}
}
