package account

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	userSvc "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/httputils"
)

// Account defines the account parameters
type Account struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	UserId    string             `json:"user_id,omitempty"`
	Username  string             `json:"username,omitempty"`
	Password  string             `json:"password,omitempty"`
	LastLogin time.Time          `json:"last_login"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

// storage provides the interface for account related operations
type storage interface {
	GetAccountByID(ctx context.Context, accountId primitive.ObjectID, ignorePasswd bool) (Account, error)
	GetAccounts(ctx context.Context, params httputils.GetQueryParams, ignorePasswd bool) (int64, []Account, error)
	GetAccountByUsername(ctx context.Context, username string, ignorePasswd bool) (Account, error)
	InsertAccount(ctx context.Context, accountDoc Account) (Account, error)
	UpdatePassword(ctx context.Context, accountDoc Account, accountId string) (Account, error)
	UpdateLastLogin(ctx context.Context, accountDoc Account, accountId string) (Account, error)
}

// Service prepares the interfaces related with this account service
type Service struct {
	userService *userSvc.Service
	storage     storage
	log         *logger.Logger
}

// NewService creates a account user service
func NewService(userService *userSvc.Service, storage storage, log *logger.Logger) *Service {
	return &Service{
		userService: userService,
		storage:     storage,
		log:         log,
	}
}

// HashPassword hashes the password
func HashPassword(password string) (string, error) {
	// uses the default cost factor to generate the hashed password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GetAccountByID extracts account data based on the userID
func (s *Service) GetAccountByID(ctx context.Context, accountId string, ignorePasswd bool) (Account, error) {
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	accountIdObject, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		s.log.Error("invalid ObjectID", zap.Error(err))
		return Account{}, err
	}

	return s.storage.GetAccountByID(ctx, accountIdObject, ignorePasswd)
}

// GetAccounts extracts account data based on the captured parameters
func (s *Service) GetAccounts(ctx context.Context, params httputils.GetQueryParams, ignorePasswd bool) (int64, []Account, error) {
	return s.storage.GetAccounts(ctx, params, ignorePasswd)
}

// GetAccountByUsername extracts account data based on the username
func (s *Service) GetAccountByUsername(ctx context.Context, username string, ignorePasswd bool) (Account, error) {
	return s.storage.GetAccountByUsername(ctx, username, ignorePasswd)
}

// InsertAccount stores account data
func (s *Service) InsertAccount(ctx context.Context, userDoc Account) (Account, error) {
	return s.storage.InsertAccount(ctx, userDoc)
}

// UpdatePassword updates password
func (s *Service) UpdatePassword(ctx context.Context, accountDoc Account, accountId string) (Account, error) {
	return s.storage.UpdatePassword(ctx, accountDoc, accountId)
}

// UpdateLastLogin updates last login
func (s *Service) UpdateLastLogin(ctx context.Context, accountDoc Account, accountId string) (Account, error) {
	return s.storage.UpdateLastLogin(ctx, accountDoc, accountId)
}

// usernameExists validates whether the username exists or not on the database
func (s *Service) usernameExists(ctx context.Context, username string) bool {
	// find user document by provided username
	_, err := s.GetAccountByUsername(ctx, username, true)

	// no error represents that the username exists
	return err == nil
}
