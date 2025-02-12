package models

import "time"

type Invitation struct {
	InvitationId  int       `json:"invitation_id"`
	AppointmentId int       `json:"appointment_id"`
	InviteeId     int       `json:"invitee_id"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
