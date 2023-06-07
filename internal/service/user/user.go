package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ardihikaru/go-whatsapp-multi-device/internal/logger"

	"github.com/ardihikaru/go-whatsapp-multi-device/pkg/utils/httputils"
)

// User defines the user parameters
type User struct {
	ID         primitive.ObjectID `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Email      string             `json:"email,omitempty"`
	Role       string             `json:"role"`
	Contact    string             `json:"contact"`
	Background string             `json:"background"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// storage provides the interface for User related operations
type storage interface {
	GetUserByID(ctx context.Context, userId primitive.ObjectID) (User, error)
	GetUsers(ctx context.Context, params httputils.GetQueryParams) (int64, []User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	InsertUser(ctx context.Context, userDoc User) (User, error)
	UpdateUser(ctx context.Context, userDoc User, userId string) (User, error)
}

// Service prepares the interfaces related with this user service
type Service struct {
	storage storage
	log     *logger.Logger
}

// NewService creates a new user service
func NewService(storage storage, log *logger.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

// hashPassword hashes the password
func hashPassword(password string) (string, error) {
	// uses the default cost factor to generate the hashed password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GetUserByID extracts user data based on the userID
func (s *Service) GetUserByID(ctx context.Context, userId string) (User, error) {
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	userIdObject, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		s.log.Error("invalid ObjectID", zap.Error(err))
		return User{}, err
	}

	return s.storage.GetUserByID(ctx, userIdObject)
}

// GetUsers extracts user data based on the captured parameters
func (s *Service) GetUsers(ctx context.Context, params httputils.GetQueryParams) (int64, []User, error) {
	return s.storage.GetUsers(ctx, params)
}

// GetUserByEmail extracts user data based on the email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return s.storage.GetUserByEmail(ctx, email)
}

// InsertUser stores user data
func (s *Service) InsertUser(ctx context.Context, userDoc User) (User, error) {
	return s.storage.InsertUser(ctx, userDoc)
}

// UpdateUser updates user data
func (s *Service) UpdateUser(ctx context.Context, userDoc User, userId string) (User, error) {
	return s.storage.UpdateUser(ctx, userDoc, userId)
}

// emailExists validates whether the email exists or not on the database
func (s *Service) emailExists(ctx context.Context, username string) bool {
	// find user document by provided email
	_, err := s.GetUserByEmail(ctx, username)

	// no error represents that the email exists
	return err == nil
}
