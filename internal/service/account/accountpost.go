package account

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// PostData defines the capture json body from the request from the POST /users
type PostData struct {
	UserId          string `json:"user_id,omitempty"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// Post creates a new user
func (s *Service) Post(ctx context.Context, accountData PostData) (Account, error) {
	var err error

	// validates user payload
	err = accountData.validateUserPostData()
	if err != nil {
		return Account{}, err
	}

	// validates if user exists
	_, err = s.userService.GetUserByID(ctx, accountData.UserId)
	if err != nil {
		return Account{}, fmt.Errorf("user not found")
	}

	// validates username
	if s.usernameExists(ctx, accountData.Username) {
		return Account{}, fmt.Errorf("username(=%s) exists", accountData.Username)
	}

	// inserts user data to the database
	accountDoc, err := buildAccountDoc(accountData)
	if err != nil {
		return Account{}, err
	}

	account, err := s.InsertAccount(ctx, accountDoc)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// buildAccountDoc prepares and builds the user document
func buildAccountDoc(accountData PostData) (Account, error) {
	var doc Account

	// hash the password
	hashedPassword, err := HashPassword(accountData.Password)
	if err != nil {
		return doc, err
	}

	// inserts user data to the database
	doc = Account{
		UserId:    accountData.UserId,
		Username:  accountData.Username,
		Password:  hashedPassword,
		LastLogin: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return doc, nil
}

// validateUserPostData validates the input data to login
func (p *PostData) validateUserPostData() error {
	if p.Username == "" {
		return errors.New("username is empty")
	}

	if p.Password == "" {
		return errors.New("password is empty")
	}

	if p.PasswordConfirm == "" {
		return errors.New("password confirmation is empty")
	}

	if p.Password != p.PasswordConfirm {
		return errors.New("password confirmation does not match")
	}

	return nil
}
