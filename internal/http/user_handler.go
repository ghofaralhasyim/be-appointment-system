package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/ghofaralhasyim/be-appointment-system/internal/services"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *UserHandler) RefreshToken(c echo.Context) error {
	sessionId, ok := c.Get("sessionId").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	var req refreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		var validationErrors []map[string]string

		for _, e := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(e.Field())
			friendlyMessage := utils.GetFriendlyErrorMessage(e)

			validationErrors = append(validationErrors, map[string]string{
				fieldName: friendlyMessage,
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "bad request - validation failed",
			"details": validationErrors,
		})
	}

	user, newToken, err := h.userService.RefreshToken(req.RefreshToken, sessionId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "refresh failed - internal server error",
			"details": nil,
		})
	}

	user.PasswordHash = ""

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": "Refresh token success",
		"data": map[string]interface{}{
			"user":  user,
			"token": newToken,
		},
	})
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c echo.Context) error {
	var req loginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		var validationErrors []map[string]string

		for _, e := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(e.Field())
			friendlyMessage := utils.GetFriendlyErrorMessage(e)

			validationErrors = append(validationErrors, map[string]string{
				fieldName: friendlyMessage,
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "bad request - validation failed",
			"details": validationErrors,
		})
	}

	user, token, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		log.Println(err)
		// not revealing whether a user is registered or not: CWE-204 CWE-203 OWASP A07:2021
		if err.Error() == "user not found" || err.Error() == "unauthorize" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid username or password",
				"details": nil,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "login failed - internal server error",
				"details": nil,
			})
		}
	}

	user.PasswordHash = ""

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login success",
		"data": map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}
