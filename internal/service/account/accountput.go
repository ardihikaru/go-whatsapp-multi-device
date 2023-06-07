package account

import (
	"context"
	"errors"
	"time"

	auth "github.com/satumedishub/sea-cucumber-api-service/pkg/authenticator"
)

// PutData defines the capture json body from the request from the POST /users
type PutData struct {
	Username         string `json:"username"`
	PasswordOriginal string `json:"password_original"`
	Password         string `json:"password"`
	PasswordConfirm  string `json:"confirm_password"`
}

// Put update password of an existing Account
func (s *Service) Put(ctx context.Context, accountData PutData, userId string) (Account, error) {
	var err error

	// validates Account payload
	err = accountData.validateAccountPutData()
	if err != nil {
		return Account{}, err
	}

	// gets the account data based on the provided username
	var accountDb Account
	accountDb, err = s.GetAccountByUsername(ctx, accountData.Username, false)
	if err != nil {
		return Account{}, err
	}

	// validate if password is incorrect
	if !auth.CheckPasswordHash(accountData.PasswordOriginal, accountDb.Password) {
		return Account{}, errors.New("incorrect username / password")
	}

	// builds Account data
	accountDoc, err := buildUpdateAccountDoc(accountData, userId)
	if err != nil {
		return Account{}, err
	}

	account, err := s.UpdatePassword(ctx, accountDoc, userId)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// buildUpdateAccountDoc prepares and builds the account document
func buildUpdateAccountDoc(accountData PutData, userId string) (Account, error) {
	var doc Account

	// hash the password
	hashedPassword, err := HashPassword(accountData.Password)
	if err != nil {
		return doc, err
	}

	// inserts user data to the database
	doc = Account{
		UserId:    userId,
		Password:  hashedPassword,
		UpdatedAt: time.Now().UTC(),
	}

	return doc, nil
}

// validateAccountPutData validates the input data to login
func (p *PutData) validateAccountPutData() error {
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
