CREATE SCHEMA stg_appointment;

CREATE TABLE stg_appointment.users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    timezone TEXT NOT NULL,  -- Stores user's timezone (e.g., 'America/New_York', 'Asia/Jakarta')
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE stg_appointment.appointments (
    appointment_id SERIAL PRIMARY KEY,
    host_id INT NOT NULL REFERENCES stg_appointment.users(user_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    appointment_time TIMESTAMPTZ NOT NULL,  -- Appointment time stored with timezone
    duration INT NOT NULL, -- Duration in minutes
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE stg_appointment.invitations (
    id SERIAL PRIMARY KEY,
    appointment_id INT NOT NULL REFERENCES stg_appointment.appointments(appointment_id) ON DELETE CASCADE,
    invitee_id INT NOT NULL REFERENCES stg_appointment.users(user_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_timezone ON stg_appointment.users(timezone);
CREATE INDEX idx_appointments_host ON stg_appointment.appointments(host_id);
CREATE INDEX idx_invitations_appointment ON stg_appointment.invitations(appointment_id);
CREATE INDEX idx_invitaions_invitee ON stg_appointment.invitations(invitee_id);