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
}

type appointmentService struct {
	appointmentRepository repositories.AppointmentRepository
}

func NewAppointmentService(appointmentRepository repositories.AppointmentRepository) AppointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
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

	return createdAppointment, nil
}
