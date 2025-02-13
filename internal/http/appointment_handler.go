package http

import (
	"log"
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
	Title      string    `json:"title" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required,ISOdate"`
	EndTime    time.Time `json:"end_time" validate:"required,ISOdate"`
	InviteeIds []int     `json:"invitee_ids" validate:"required"`
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
		errMsg := "Invalid request"
		if strings.Contains(err.Error(), "parsing time") {
			errMsg = "date must in ISO 8601 format"
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": errMsg, "detail": nil})
	}

	if err := c.Validate(&req); err != nil {
		var validationErrors []map[string]string

		for _, e := range err.(validator.ValidationErrors) {
			field, friendlyMessage := utils.GetFriendlyErrorMessage(e, req)

			validationErrors = append(validationErrors, map[string]string{
				field: friendlyMessage,
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "bad request - validation failed",
			"details": validationErrors,
		})
	}

	dataAppointment := models.Appointment{
		Title:      req.Title,
		HostId:     userId,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		InviteeIds: req.InviteeIds,
	}

	createdAppointment, err := h.appointmentService.CreateAppointment(&dataAppointment)
	if err != nil {
		log.Println(err)
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

func (h *AppointmentHandler) GetAppointments(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	appointments, err := h.appointmentService.GetAppointmentsByUserId(userId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed create appointment - internal server error",
			"details": nil,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    appointments,
	})

}
