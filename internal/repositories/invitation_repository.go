package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/lib/pq"
)

type InvitationRepository interface {
	InsertInvitation(tx *sql.Tx, invitations []models.Invitation) error
	GetInvitations(userId int) ([]models.AppointmentInvitation, error)
	UpdateStatusInvitation(userId int, invId int, status string) error
}

type invitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) InsertInvitation(tx *sql.Tx, invitations []models.Invitation) error {
	if len(invitations) == 0 {
		return nil
	}

	var appointmentIDs []int64
	var inviteeIDs []int64
	var statuses []string
	var notes []string
	var createdAts []time.Time

	for _, inv := range invitations {
		appointmentIDs = append(appointmentIDs, int64(inv.AppointmentId))
		inviteeIDs = append(inviteeIDs, int64(inv.InviteeId))
		statuses = append(statuses, inv.Status)
		notes = append(notes, "")
		createdAts = append(createdAts, time.Now())
	}

	query := `
		INSERT INTO stg_appointment.invitations 
			(appointment_id, invitee_id, status, notes, created_at)
		SELECT * FROM UNNEST($1::bigint[], $2::bigint[], $3::text[], $4::text[], $5::timestamptz[])
	`

	_, err := tx.Exec(query, pq.Array(appointmentIDs), pq.Array(inviteeIDs), pq.Array(statuses), pq.Array(notes), pq.Array(createdAts))
	return err
}

func (r *invitationRepository) GetInvitations(userId int) ([]models.AppointmentInvitation, error) {
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
			JOIN stg_appointment.invitations i ON a.appointment_id = i.appointment_id
			WHERE i.invitee_id = $1     
				AND a.host_id != $1        
				AND i.status = 'pending' 
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
			i.invitation_id,
			i.invitee_id,
			i.status,
			i.created_at AS invitation_created_at
		FROM appointment_details ad
		JOIN stg_appointment.invitations i ON 
			ad.appointment_id = i.appointment_id 
			AND i.invitee_id = $1
		ORDER BY ad.start_time;
	`

	rows, err := r.db.Query(query, userId)
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

func (r *invitationRepository) UpdateStatusInvitation(userId int, invId int, status string) error {
	query := `
		UPDATE  stg_appointment.invitations
		SET 
			status = $1
		WHERE 
			invitee_id = $2 AND invitation_id = $3;
	`

	_, err := r.db.Exec(query, status, userId, invId)
	return err
}
