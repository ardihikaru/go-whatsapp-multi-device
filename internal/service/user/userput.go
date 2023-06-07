package user

import (
	"context"
	"errors"
	"time"

	uRole "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/user/role"
)

// PutData defines the capture json body from the request from the POST /users
type PutData struct {
	Name       string `json:"name"`
	Role       string `json:"role"`
	Contact    string `json:"contact,omitempty"`
	Background string `json:"background,omitempty"`
}

// Put creates a new user
func (s *Service) Put(ctx context.Context, userData PutData, userId string) (User, error) {
	var err error

	// validates user payload
	err = userData.validateUserPutData()
	if err != nil {
		return User{}, err
	}

	// inserts user data to the database
	userDoc, err := buildUpdateUserDoc(userData)
	if err != nil {
		return User{}, err
	}

	user, err := s.UpdateUser(ctx, userDoc, userId)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// buildUpdateUserDoc prepares and builds the user document
func buildUpdateUserDoc(userData PutData) (User, error) {
	var doc User

	// inserts user data to the database
	doc = User{
		Name:       userData.Name,
		Role:       userData.Role,
		Contact:    userData.Contact,
		Background: userData.Background,
		UpdatedAt:  time.Now().UTC(),
	}

	return doc, nil
}

// validateUserPutData validates the input data to login
func (p *PutData) validateUserPutData() error {
	if p.Name == "" {
		return errors.New("name is empty")
	}

	if p.Name == "" {
		return errors.New("username is empty")
	}

	if !uRole.GetRoleMap()[p.Role] {
		return errors.New("invalid role")
	}

	return nil
}
