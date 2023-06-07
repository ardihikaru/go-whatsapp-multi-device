package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	uRole "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/user/role"
)

// PostData defines the capture json body from the request from the POST /users
type PostData struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Contact    string `json:"contact,omitempty"`
	Background string `json:"background,omitempty"`
}

// Post creates a new user
func (s *Service) Post(ctx context.Context, userData PostData) (User, error) {
	var err error

	// validates user payload
	err = userData.validateUserPostData()
	if err != nil {
		return User{}, err
	}

	// validates email
	if s.emailExists(ctx, userData.Email) {
		return User{}, fmt.Errorf("email(=%s) exists", userData.Email)
	}

	// inserts user data to the database
	userDoc, err := buildUserDoc(userData)
	if err != nil {
		return User{}, err
	}

	user, err := s.InsertUser(ctx, userDoc)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// buildUserDoc prepares and builds the user document
func buildUserDoc(userData PostData) (User, error) {
	// var rawDoc RawDoc
	var doc User

	// inserts user data to the database
	doc = User{
		Name:       userData.Name,
		Email:      userData.Email,
		Role:       userData.Role,
		Contact:    userData.Contact,
		Background: userData.Background,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return doc, nil
}

// validateUserPostData validates the input data to login
func (p *PostData) validateUserPostData() error {
	if p.Name == "" {
		return errors.New("name is empty")
	}

	_, err := mail.ParseAddress(p.Email)
	if err != nil {
		return errors.New("email is empty")
	}

	if !uRole.GetRoleMap()[p.Role] {
		return errors.New("invalid role")
	}

	return nil
}
