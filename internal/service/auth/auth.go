package auth

import (
	"context"
	"errors"
	"time"

	accSvc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/account"
	"github.com/go-chi/jwtauth/v5"

	auth "github.com/ardihikaru/go-whatsapp-multi-device/pkg/authenticator"
)

// LoginData is the input JSON body captured from the login request
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AccessTokenPayload is the access token payload
// this payload is given when login request is validated and valid
type AccessTokenPayload struct {
	AccessToken string       `json:"access_token"`
	ExpiredIn   int64        `json:"expired_in"`
	IssuedAt    int64        `json:"issued_at"`
	Session     auth.Session `json:"session"`
}

// Service prepares the interfaces related with this auth service
type Service struct {
	accService *accSvc.Service
	jwtExpTime int64
	tokenAuth  *jwtauth.JWTAuth
}

// NewService creates a new auth service
func NewService(accService *accSvc.Service, jwtExpTime int64, tokenAuth *jwtauth.JWTAuth) *Service {
	return &Service{
		accService: accService,
		jwtExpTime: jwtExpTime,
		tokenAuth:  tokenAuth,
	}
}

func (svc *Service) Authorize(ctx context.Context, loginData LoginData) (*AccessTokenPayload, error) {
	var err error

	// validates login data
	err = loginData.Validate()
	if err != nil {
		return nil, err
	}

	// gets the user data based on the provided username
	var account accSvc.Account
	account, err = svc.accService.GetAccountByUsername(ctx, loginData.Username, false)
	if err != nil {
		return nil, err
	}

	// validate if password is incorrect
	if !auth.CheckPasswordHash(loginData.Password, account.Password) {
		return nil, errors.New("invalid password")
	}

	// builds the JWT claim options
	durationIn := time.Duration(svc.jwtExpTime) * time.Second
	jwtClaims := auth.JWTClaims{
		auth.ClaimAccountIdKey: account.ID.Hex(),
		auth.ClaimUserIdKey:    account.UserId,
		auth.ClaimUsernameKey:  account.Username,
		auth.ClaimExpiredInKey: jwtauth.ExpireIn(durationIn),
		auth.ClaimIssuedAtKey:  jwtauth.EpochNow(),
	}

	// begins to create the access token
	accessToken := auth.CreateAccessToken(svc.tokenAuth, jwtClaims)

	// builds access token payload
	payload := &AccessTokenPayload{
		AccessToken: accessToken,
		ExpiredIn:   jwtauth.ExpireIn(durationIn),
		IssuedAt:    jwtauth.EpochNow(),
		Session: auth.Session{
			AccountId: account.ID.Hex(),
			UserId:    account.UserId,
			Username:  account.Username,
		},
	}

	// once authorization succeed, update last login
	accountDoc := accSvc.Account{
		ID:        account.ID,
		LastLogin: time.Now().UTC(),
	}

	_, err = svc.accService.UpdateLastLogin(ctx, accountDoc, account.ID.Hex())
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// Validate validates the input data to login
func (d *LoginData) Validate() error {
	if d.Username == "" {
		return errors.New("username is empty")
	}
	if d.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}
