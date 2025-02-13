package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ghofaralhasyim/be-appointment-system/internal/services"
	"github.com/labstack/echo/v4"
)

type InvitationHandler struct {
	invitationService services.InvitationService
}

func NewInvitationHandler(invitationService services.InvitationService) *InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
	}
}

func (h *InvitationHandler) GetInvitations(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	invitations, err := h.invitationService.GetInvitations(userId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed retrieve invitations - internal server error",
			"details": nil,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "appointment created",
		"data":    invitations,
	})
}

func (h *InvitationHandler) AcceptInvitation(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	invId := c.Param("invitationId")
	invIdInt, err := strconv.Atoi(invId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid address id", "detail": err})
	}

	err = h.invitationService.UpdateStatusInvitation(userId, invIdInt, "accepted")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed accept invitation - internal server error",
			"details": nil,
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": "invitation accepted",
		"data":    nil,
	})
}

func (h *InvitationHandler) RejectInvitation(c echo.Context) error {
	userId, ok := c.Get("userId").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid token or malformed token",
			"details": nil,
		})
	}

	invId := c.Param("invitationId")
	invIdInt, err := strconv.Atoi(invId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid address id", "detail": err})
	}

	err = h.invitationService.UpdateStatusInvitation(userId, invIdInt, "rejected")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed reject invitation - internal server error",
			"details": nil,
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": "invitation rejected",
		"data":    nil,
	})
}
