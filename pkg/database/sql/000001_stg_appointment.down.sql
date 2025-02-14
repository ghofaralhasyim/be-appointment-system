DROP TABLE IF EXISTS stg_appointment.invitations;
DROP TABLE IF EXISTS stg_appointment.appointments;
DROP TABLE IF EXISTS stg_appointment.users;

DROP INDEX IF EXISTS stg_appointment.idx_users_username;
DROP INDEX IF EXISTS stg_appointment.idx_users_role;
DROP INDEX IF EXISTS stg_appointment.idx_users_timezone;

DROP INDEX IF EXISTS stg_appointment.idx_appointments_host_id;
DROP INDEX IF EXISTS stg_appointment.idx_appointments_start_time;
DROP INDEX IF EXISTS stg_appointment.idx_appointments_end_time;

DROP INDEX IF EXISTS stg_appointment.idx_invitations_appointment_id;
DROP INDEX IF EXISTS stg_appointment.idx_invitations_invitee_id;
DROP INDEX IF EXISTS stg_appointment.idx_invitations_status;

DROP SCHEMA IF EXISTS stg_appointment CASCADE;