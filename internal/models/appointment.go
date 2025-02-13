package models

import "time"

type Appointment struct {
	AppointmentId int       `json:"appointment_id"`
	HostId        int       `json:"host_id"`
	Title         string    `json:"title"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	CreatedAt     time.Time `json:"created_at"`
	InviteeIds    []int     `json:"invitee_ids,omitempty"`
}

type AppointmentInvitation struct {
	Appointment
	TotalAttendants int    `json:"total_attendants"`
	InvitationId    int    `json:"invitation_id"`
	Invitee_id      int    `json:"invitee_id"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
	Host            User   `json:"host"`
	Attendants      []User `json:"attendants"`
}
