CREATE SCHEMA stg_appointment;

CREATE TABLE stg_appointment.users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'staff',
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
    invitation_id SERIAL PRIMARY KEY,
    appointment_id INT NOT NULL REFERENCES stg_appointment.appointments(appointment_id) ON DELETE CASCADE,
    invitee_id INT REFERENCES stg_appointment.users(user_id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    notes VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);