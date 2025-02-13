package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
)

type UserRepository interface {
	GetUsers() ([]models.User, error)
	GetUserById(userId int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserTimezone(userId int, timezone string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT
			u.user_id, u.name, u.username, u.role, u.timezone, timezone(u.timezone, u.created_at) as created_at,
			timezone(u.timezone, u.updated_at) as updated_at
		FROM stg_appointment.users u
			WHERE u.username = $1 AND u.deleted_at IS NULL
		LIMIT 1;
	`

	var user models.User
	var updated, deleted sql.NullString

	err := r.db.QueryRow(query, username).Scan(
		&user.UserId, &user.Name, &user.Username, &user.Role, &user.Timezone,
		&user.CreatedAt, &updated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	layoutTime := "2006-01-02T15:04:05Z"
	if updated.Valid {
		parsedTime, err := time.Parse(layoutTime, updated.String)
		if err != nil {
			return nil, err
		}
		user.UpdatedAt = parsedTime
	}

	if deleted.Valid {
		parsedTime, err := time.Parse(layoutTime, deleted.String)
		if err != nil {
			return nil, err
		}
		user.DeletedAt = parsedTime
	}

	return &user, nil
}

func (r *userRepository) GetUserById(userId int) (*models.User, error) {
	query := `
		SELECT
			u.user_id, u.name, u.username, u.timezone, timezone(u.timezone, u.created_at) as created_at,
			timezone(u.timezone, u.updated_at) as updated_at
		FROM stg_appointment.users u WHERE u.user_id = $1 AND u.deleted_at IS NULL
		LIMIT 1;
	`

	var user models.User
	var updated sql.NullString

	err := r.db.QueryRow(query, userId).Scan(
		&user.UserId, &user.Name, &user.Username, &user.Timezone, &user.CreatedAt, &updated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	layoutTime := "2006-01-02T15:04:05Z"
	if updated.Valid {
		parsedTime, err := time.Parse(layoutTime, updated.String)
		if err != nil {
			return nil, err
		}
		user.UpdatedAt = parsedTime
	}

	return &user, nil
}

func (r *userRepository) GetUsers() ([]models.User, error) {
	query := `
		SELECT
			u.user_id, u.name, u.username, u.timezone, 
			timezone(u.timezone, u.created_at) as created_at,
			timezone(u.timezone, u.updated_at) as updated_at
		FROM stg_appointment.users u 
		WHERE u.deleted_at IS NULL;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var updated sql.NullTime

		err := rows.Scan(
			&user.UserId, &user.Name, &user.Username, &user.Timezone, &user.CreatedAt, &updated,
		)
		if err != nil {
			return nil, err
		}

		if updated.Valid {
			user.UpdatedAt = updated.Time
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) UpdateUserTimezone(userId int, timezone string) error {
	query := `
		UPDATE stg_appointment.users
		SET
			timezone = $1
		WHERE user_id = $2;
	`

	_, err := r.db.Exec(query, timezone, userId)
	return err
}
