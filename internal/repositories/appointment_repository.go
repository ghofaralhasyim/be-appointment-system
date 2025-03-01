package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
)

type AppointmentRepository interface {
	InsertAppointment(tx *sql.Tx, appointment *models.Appointment) (*models.Appointment, error)
	BeginAppointmentTx() (*sql.Tx, error)

	GetAppointmentsByUserId(userId int, startDate, endDate time.Time) ([]models.AppointmentInvitation, error)
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
			(host_id, title, start_time, end_time, created_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING appointment_id;
	`

	err := r.db.QueryRow(
		query, appointment.HostId, appointment.Title, appointment.StartTime, appointment.EndTime,
		appointment.CreatedAt,
	).Scan(&appointment.AppointmentId)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *appointmentRepository) GetAppointmentsByUserId(userId int, startDate, endDate time.Time) ([]models.AppointmentInvitation, error) {
	query := `
		WITH user_tz AS (
			SELECT timezone
			FROM stg_appointment.users
			WHERE user_id = $1
		),
		appointment_details AS (
			SELECT 
				a.appointment_id,
				a.title,
				timezone((SELECT timezone FROM user_tz), a.start_time) AS start_time,
				timezone((SELECT timezone FROM user_tz), a.end_time) AS end_time,
				a.created_at AS appointment_created_at,
				a.host_id,
				-- Host information
				jsonb_build_object(
					'username', host.username,
					'name', host.name,
					'timezone', host.timezone
				) AS host,
				-- Get total attendants count
				(
					SELECT COUNT(*)
					FROM stg_appointment.invitations inv
					WHERE inv.appointment_id = a.appointment_id
				) AS total_attendants,
				-- Limited attendants list (only 3)
				(
					SELECT jsonb_agg(attendant_info)
					FROM (
						SELECT jsonb_build_object(
							'username', u.username,
							'name', u.name,
							'timezone', u.timezone,
							'status', inv.status,
							'invitation_id', inv.invitation_id,
							'invitee_id', inv.invitee_id
						) as attendant_info
						FROM stg_appointment.invitations inv
						JOIN stg_appointment.users u ON inv.invitee_id = u.user_id
						WHERE inv.appointment_id = a.appointment_id
						LIMIT 3
					) limited_attendants
				) AS limited_attendants
			FROM stg_appointment.appointments a
			JOIN stg_appointment.users host ON a.host_id = host.user_id
			WHERE a.start_time BETWEEN $2 AND $3
				AND (
					a.host_id = $1  -- User is host
					OR EXISTS (
						SELECT 1 
						FROM stg_appointment.invitations i 
						WHERE i.appointment_id = a.appointment_id 
						AND i.invitee_id = $1
						AND (
							a.host_id = $1  -- Allow all statuses if host
							OR i.status = 'accepted'  -- Only accepted if not host
						)
					)
				)
		)
		SELECT 
			ad.appointment_id,
			ad.title,
			ad.start_time,
			ad.end_time,
			ad.appointment_created_at,
			ad.host,
			ad.total_attendants,
			COALESCE(ad.limited_attendants, '[]'::jsonb) as attendants,
			-- Invitation details for the current user
			COALESCE(i.invitation_id, 0) AS invitation_id,
			COALESCE(i.invitee_id, ad.host_id) AS invitee_id,
			COALESCE(i.status, 
				CASE 
					WHEN ad.host_id = $1 THEN 'host'
					ELSE NULL 
				END
			) AS status,
			COALESCE(i.created_at, ad.appointment_created_at) AS invitation_created_at
		FROM appointment_details ad
		LEFT JOIN stg_appointment.invitations i ON 
			ad.appointment_id = i.appointment_id 
			AND i.invitee_id = $1
		ORDER BY ad.start_time;
	`

	rows, err := r.db.Query(query, userId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying appointments: %w", err)
	}
	defer rows.Close()

	var appointments []models.AppointmentInvitation

	for rows.Next() {
		var appointment models.AppointmentInvitation
		var hostJSON, attendantsJSON []byte
		var invitationID sql.NullInt64

		err := rows.Scan(
			&appointment.AppointmentId,
			&appointment.Title,
			&appointment.StartTime,
			&appointment.EndTime,
			&appointment.CreatedAt,
			&hostJSON,
			&appointment.TotalAttendants,
			&attendantsJSON,
			&invitationID,
			&appointment.Invitee_id,
			&appointment.Status,
			&appointment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning appointment row: %w", err)
		}

		if invitationID.Valid {
			appointment.InvitationId = int(invitationID.Int64)
		}

		if err := json.Unmarshal(hostJSON, &appointment.Host); err != nil {
			return nil, fmt.Errorf("error unmarshaling host data: %w", err)
		}

		if err := json.Unmarshal(attendantsJSON, &appointment.Attendants); err != nil {
			return nil, fmt.Errorf("error unmarshaling attendants data: %w", err)
		}

		appointments = append(appointments, appointment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointment rows: %w", err)
	}

	return appointments, nil
}
