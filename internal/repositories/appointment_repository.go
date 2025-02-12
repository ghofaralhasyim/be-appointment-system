package repositories

import (
	"database/sql"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
)

type AppointmentRepository interface {
	InsertAppointment(tx *sql.Tx, appointment *models.Appointment) (*models.Appointment, error)
	BeginAppointmentTx() (*sql.Tx, error)
}

type appointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) BeginAppointmentTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *appointmentRepository) InsertAppointment(tx *sql.Tx, appointment *models.Appointment) (*models.Appointment, error) {
	query := `
		INSERT INTO stg_appointment.appointments
			(host_id, title, appointment_time, duration, created_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING appointment_id;
	`

	err := r.db.QueryRow(
		query, appointment.HostId, appointment.Title, appointment.AppointmentTime, appointment.Duration,
		appointment.CreatedAt,
	).Scan(&appointment.AppointmentId)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}
