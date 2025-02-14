CREATE SCHEMA stg_appointment;

CREATE TABLE stg_appointment.users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'staff',
    timezone TEXT NOT NULL,  -- timezone (e.g., 'America/New_York', 'Asia/Jakarta')
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE stg_appointment.appointments (
    appointment_id SERIAL PRIMARY KEY,
    host_id INT NOT NULL REFERENCES stg_appointment.users(user_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE stg_appointment.invitations (
    invitation_id SERIAL PRIMARY KEY,
    appointment_id INT NOT NULL REFERENCES stg_appointment.appointments(appointment_id) ON DELETE CASCADE,
    invitee_id INT REFERENCES stg_appointment.users(user_id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    notes VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_username ON stg_appointment.users (username);
CREATE INDEX idx_users_timezone ON stg_appointment.users (timezone);

CREATE INDEX idx_appointments_host_id ON stg_appointment.appointments (host_id);
CREATE INDEX idx_appointments_start_time ON stg_appointment.appointments (start_time);
CREATE INDEX idx_appointments_end_time ON stg_appointment.appointments (end_time);

CREATE INDEX idx_invitations_appointment_id ON stg_appointment.invitations (appointment_id);
CREATE INDEX idx_invitations_invitee_id ON stg_appointment.invitations (invitee_id);
CREATE INDEX idx_invitations_status ON stg_appointment.invitations (status);