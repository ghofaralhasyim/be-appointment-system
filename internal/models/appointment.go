package models

import "time"

type Appointment struct {
	AppointmentId   int       `json:"appointment_id"`
	HostId          int       `json:"host_id"`
	Title           string    `json:"title"`
	AppointmentTime time.Time `json:"appointment_time"`
	Duration        int       `json:"duration"`
	CreatedAt       time.Time `json:"created_at"`
}
