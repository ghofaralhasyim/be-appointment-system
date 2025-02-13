package services

import (
	"fmt"
	"log"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/ghofaralhasyim/be-appointment-system/internal/repositories"
)

type AppointmentService interface {
	CreateAppointment(appointment *models.Appointment) (*models.Appointment, error)
	GetAppointmentsByUserId(userId int) ([]models.AppointmentInvitation, error)
}

type appointmentService struct {
	appointmentRepository repositories.AppointmentRepository
	invitationRepository  repositories.InvitationRepository
}

func NewAppointmentService(appointmentRepository repositories.AppointmentRepository, invitationRepository repositories.InvitationRepository) AppointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
		invitationRepository:  invitationRepository,
	}
}

func (s *appointmentService) CreateAppointment(appointment *models.Appointment) (*models.Appointment, error) {

	tx, err := s.appointmentRepository.BeginAppointmentTx()
	if err != nil {
		return nil, fmt.Errorf("error create appointment: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", p)
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				err = commitErr
			}
		}
	}()

	appointment.CreatedAt = time.Now().UTC()

	createdAppointment, err := s.appointmentRepository.InsertAppointment(tx, appointment)
	if err != nil {
		return nil, err
	}

	var invitees []models.Invitation
	for _, item := range appointment.InviteeIds {
		invite := models.Invitation{
			AppointmentId: createdAppointment.AppointmentId,
			InviteeId:     item,
			Status:        "pending",
			Notes:         "",
			CreatedAt:     time.Now(),
		}
		invitees = append(invitees, invite)
	}

	err = s.invitationRepository.InsertInvitation(tx, invitees)
	if err != nil {
		return nil, err
	}

	return createdAppointment, nil
}

func (s *appointmentService) GetAppointmentsByUserId(userId int) ([]models.AppointmentInvitation, error) {
	date := "2025-02-13"

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, err
	}

	endDate := parsedDate.AddDate(0, 0, 4)

	return s.appointmentRepository.GetAppointmentsByUserId(userId, parsedDate, endDate)
}
