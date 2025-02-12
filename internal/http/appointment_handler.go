package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/ghofaralhasyim/be-appointment-system/internal/services"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AppointmentHandler struct {
	appointmentService services.AppointmentService
}

func NewAppointmentHandler(appointmentService services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

type appointmentRequest struct {
	Title           string    `json:"title" validate:"required"`
	AppointmentTime time.Time `json:"appointment_time" validate:"required"`
	Duration        int       `json:"duration" validate:"required"`
}

func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	var req appointmentRequest

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

	dataAppointment := models.Appointment{
		Title:           req.Title,
		HostId:          userId,
		AppointmentTime: req.AppointmentTime,
		Duration:        req.Duration,
	}

	createdAppointment, err := h.appointmentService.CreateAppointment(&dataAppointment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed create appointment - internal server error",
			"details": nil,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "appointment created",
		"data":    createdAppointment,
	})
}
