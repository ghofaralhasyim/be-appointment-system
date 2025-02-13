package services

import (
	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/ghofaralhasyim/be-appointment-system/internal/repositories"
)

type InvitationService interface {
	GetInvitations(userId int) ([]models.AppointmentInvitation, error)
	UpdateStatusInvitation(userId int, invId int, status string) error
}

type invitationService struct {
	invitationRepository repositories.InvitationRepository
}

func NewInvitationService(invitationRepository repositories.InvitationRepository) InvitationService {
	return &invitationService{
		invitationRepository: invitationRepository,
	}
}

func (s *invitationService) GetInvitations(userId int) ([]models.AppointmentInvitation, error) {
	return s.invitationRepository.GetInvitations(userId)
}

func (s *invitationService) UpdateStatusInvitation(userId int, invId int, status string) error {
	return s.invitationRepository.UpdateStatusInvitation(userId, invId, status)
}
