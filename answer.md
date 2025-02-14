### Technical Question Answer

#### 1. Timezone Conflicts: How would you handle timezone conflicts between participants in an appointment?

**Data storing and transfer**: I handle timezone conflict by storing data in UTC format (db) and format in ISO 8601 when retrieve data from database to ensure consistency across differetent time zones.

**User Input:** Make sure the appointment request start time and end time is falls within working hours 08:00 – 17.00 of all participant using custom validation. If it doesn’t user will be notifying by it user interface and invitation wouldn’t sent.

#### 2. Database Optimization: How can you optimize database queries to efficiently fetch user-specific appointments?

I apply several database optimization technique to improve the performance, flexibility, and debugging. Instead use ORM I use raw SQL for gain full control over query execution and indexing, ensuring better efficiency.

- Create indexes on frequently queried columns to reduce database scan time.
- Use CTE for readability & performance: breaking complex queries into some logical step to reduce rendundant queries. Also use PostgreSQL feature like jsonb_agg to construct pre-aggregated JSON list.
- Use Exist instead of additional join table.

\*on the code appointment repository its showing how flexible and using raw SQL instead of ORM for complex query.

#### 3. Database Optimization: How can you optimize database queries to efficiently fetch user-specific appointments?

#### Application Feature & UI/UX:

- Create alternative display modes appointment. Allow user to switch between schedule view or calendar view, similar to Google calendar.
- Notification email when a user receives an invitation and also set reminders for upcoming appointments.
- Attendance confirmation, ensure the user is know the invitations and confirm their availability by reject or accept it.
- Attendant list.
- Time conflict resolution: Suggest alternative times if a new appointment conflicts with an existing one.
- JWT refresh token: Instead of automatic logout, implement a session refresh mechanism—either through a popup or auto-refresh if the user is still active

#### Security:

- Credential verification: Implement email and password verification to enhance login security.
- Session management: Allow users to terminate all active sessions. Since sessions app supported multiple session and session are stored in Redis this feature will helps secure accounts if a device is lost.

#### 4. Session Management: How would you manage user sessions securely while keeping them lightweight (e.g., avoiding large JWT payloads)?

I use Redis to store user session data for better session management, allowing easy modification or termination of sessions. Managing sessions with only JWT can be challenging since we have to wait for the token to expire. JWT is designed for verifying communicators, not for session management.

In my approach, the JWT token only contains the session ID. If the frontend needs user data, it can request it via an API call or receive it in the login response, keeping the token lightweight.
