package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ghofaralhasyim/be-appointment-system/internal/models"
	"github.com/ghofaralhasyim/be-appointment-system/internal/repositories"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Authenticate(email string, password string) (*models.User, *models.JwtToken, error)
	RefreshToken(refreshToken string, sessionId string) (*models.User, *models.JwtToken, error)
}

type userService struct {
	userRepository  repositories.UserRepository
	redisRepository repositories.RedisRepository
}

func NewUserService(userRepository repositories.UserRepository, redisRepository repositories.RedisRepository) UserService {
	return &userService{
		userRepository:  userRepository,
		redisRepository: redisRepository,
	}
}

func (s *userService) RefreshToken(refreshToken string, sessionId string) (*models.User, *models.JwtToken, error) {

	_, err := utils.VerifyToken(refreshToken, true)
	if err != nil {
		return nil, nil, err
	}

	dataJSON, err := s.redisRepository.Get(context.Background(), sessionId)
	if err != nil {
		return nil, nil, err
	}

	var dataSession map[string]interface{}
	if err := json.Unmarshal([]byte(dataJSON), &dataSession); err != nil {
		return nil, nil, fmt.Errorf("refresh token: failed to unmarshal data session : %w", err)
	}

	userId := dataSession["user_id"].(float64)
	user, err := s.userRepository.GetUserById(int(userId))
	if err != nil {
		return nil, nil, fmt.Errorf("refresh token: failed to get user data - %w", err)
	}

	newToken, err := utils.GenerateSessionToken(sessionId)
	if err != nil {
		return nil, nil, fmt.Errorf("refresh token: failed to generate new session - %w", err)
	}

	dataUser, err := json.Marshal(map[string]interface{}{
		"user_id":       user.UserId,
		"timezone":      user.Timezone,
		"role":          user.Role,
		"access_token":  newToken.AccessToken,
		"refresh_token": newToken.RefreshToken,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling user: %w", err)
	}

	expTime, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRE_HOURS"))
	if err != nil {
		return nil, nil, fmt.Errorf("auth error: %w", err)
	}

	err = s.redisRepository.Set(context.Background(), sessionId, dataUser, time.Duration(expTime*int(time.Hour)))
	if err != nil {
		return nil, nil, fmt.Errorf("auth error: storing redis: %w", err)
	}

	return user, newToken, nil
}

func (s *userService) Authenticate(username string, password string) (*models.User, *models.JwtToken, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, fmt.Errorf("unauthorize")
	}

	timestamp := time.Now().Unix()
	shortTimestamp := timestamp % 10000

	sessionId := fmt.Sprintf("session:%d-%04d", user.UserId, shortTimestamp)
	jwt, err := utils.GenerateSessionToken(sessionId)
	if err != nil {
		return nil, nil, err
	}

	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if err != nil {
		return nil, nil, fmt.Errorf("auth error: %w", err)
	}

	dataUser, err := json.Marshal(map[string]interface{}{
		"user_id":       user.UserId,
		"timezone":      user.Timezone,
		"role":          user.Role,
		"access_token":  jwt.AccessToken,
		"refresh_token": jwt.RefreshToken,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling user: %w", err)
	}

	err = s.redisRepository.Set(context.Background(), sessionId, dataUser, time.Duration(expTime*int(time.Hour)))
	if err != nil {
		return nil, nil, fmt.Errorf("auth error: storing redis: %w", err)
	}

	return user, jwt, nil
}
