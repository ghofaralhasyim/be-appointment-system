package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
)

type UserRepository interface {
	GetUserById(userId int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
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
			u.user_id, u.name, u.username, u.role, u.password_hash, u.timezone, timezone(u.timezone, u.created_at) as created_at,
			timezone(u.timezone, u.updated_at) as updated_at
		FROM stg_appointment.users u
			WHERE u.username = $1 AND u.deleted_at IS NULL
		LIMIT 1;
	`

	var user models.User
	var updated, deleted sql.NullString

	err := r.db.QueryRow(query, username).Scan(
		&user.UserId, &user.Name, &user.Username, &user.Role, &user.PasswordHash, &user.Timezone,
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
			u.user_id, u.name, u.username, u.password_hash, u.timezone, timezone(u.timezone, u.created_at) as created_at,
			timezone(u.timezone, u.updated_at) as updated_at
		FROM stg_appointment.users u WHERE u.user_id = $1 AND u.deleted_at IS NULL
		LIMIT 1;
	`

	var user models.User
	var updated sql.NullString

	err := r.db.QueryRow(query, userId).Scan(
		&user.UserId, &user.Name, &user.Username, &user.PasswordHash, &user.Timezone, &user.CreatedAt, &updated,
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
