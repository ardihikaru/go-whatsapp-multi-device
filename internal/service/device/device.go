package account

import (
	"context"
	"fmt"
	"time"

	"github.com/ardihikaru/go-modules/pkg/logger"
	"github.com/ardihikaru/go-modules/pkg/utils/common"
	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
)

// RegisterPayload is the input JSON body captured from the register request
type RegisterPayload struct {
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	WebhookUrl string `json:"webhook_url"`
}

// Device is the device object
type Device struct {
	ID         string    `json:"_id,omitempty"`
	JID        string    `json:"jid,omitempty"`
	Phone      string    `json:"phone"`
	Name       string    `json:"name"`
	WebhookUrl string    `json:"webhook_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// storage provides the interface for account related operations
type storage interface {
	GetDeviceByID(ctx context.Context, id string) (Device, error)
	GetDeviceByPhone(ctx context.Context, phone string) (Device, error)
	GetDeviceByJID(ctx context.Context, id string) (Device, error)
	GetDevices(ctx context.Context, params httputils.GetQueryParams) (int64, []Device, error)
	InsertDevice(ctx context.Context, doc Device) (Device, error)
	UpdateDeviceName(ctx context.Context, id, deviceName string) error
	UpdateWebhook(ctx context.Context, id, webhook string) error
	UpdateJID(ctx context.Context, jid, id string) error
}

// Service prepares the interfaces related with this account service
type Service struct {
	storage storage
	log     *logger.Logger
}

// NewService creates a device service
func NewService(storage storage, log *logger.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

// GetDeviceByID extracts device data based on the ID
func (s *Service) GetDeviceByID(ctx context.Context, id string) (Device, error) {
	return s.storage.GetDeviceByID(ctx, id)
}

// GetDeviceByPhone extracts device data based on the phone
func (s *Service) GetDeviceByPhone(ctx context.Context, phone string) (Device, error) {
	// enriches with `+` symbol if missing
	if phone[0:1] != "+" {
		phone = fmt.Sprintf("+%s", phone)
	}
	return s.storage.GetDeviceByPhone(ctx, phone)
}

// GetDeviceByJID extracts device data based on the JID
func (s *Service) GetDeviceByJID(ctx context.Context, jid string) (Device, error) {
	return s.storage.GetDeviceByJID(ctx, jid)
}

// GetDevices fetches device data
func (s *Service) GetDevices(ctx context.Context, params httputils.GetQueryParams) (int64, []Device, error) {
	return s.storage.GetDevices(ctx, params)
}

// InsertDevice stores device data
func (s *Service) InsertDevice(ctx context.Context, doc Device) (Device, error) {
	return s.storage.InsertDevice(ctx, doc)
}

// UpdateJID updates device data
func (s *Service) UpdateJID(ctx context.Context, jid, id string) error {
	return s.storage.UpdateJID(ctx, jid, id)
}

// UpdateDeviceName updates device name
func (s *Service) UpdateDeviceName(ctx context.Context, id, deviceName string) error {
	return s.storage.UpdateDeviceName(ctx, id, deviceName)
}

// UpdateWebhook updates device webhook
func (s *Service) UpdateWebhook(ctx context.Context, id, webhook string) error {
	return s.storage.UpdateWebhook(ctx, id, webhook)
}

func (s *Service) Register(ctx context.Context, payload RegisterPayload) (Device, error) {
	var err error

	payload.Sanitize()

	// validates login data
	err = payload.Validate()
	if err != nil {
		return Device{}, err
	}

	// builds device object
	doc := Device{
		Phone:      payload.Phone,
		Name:       payload.Name,
		WebhookUrl: payload.WebhookUrl,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	// validates if this phone exists in the database
	// if error NOT found, means that this phone exist in the database
	_, err = s.GetDeviceByPhone(ctx, payload.Phone)
	if err == nil {
		return Device{}, fmt.Errorf("phone exists")
	}

	device, err := s.InsertDevice(ctx, doc)
	if err != nil {
		return Device{}, err
	}

	return device, nil
}

// Validate validates the input data
func (d *RegisterPayload) Validate() error {

	return nil
}

// Sanitize sanitizes the input data
func (d *RegisterPayload) Sanitize() {
	withPlusSymbol := true
	d.Phone = common.SanitizePhone(d.Phone, &withPlusSymbol)
}
