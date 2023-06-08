package session

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"

	"go.mau.fi/whatsmeow/types"
	"go.uber.org/zap"

	"github.com/ardihikaru/go-modules/pkg/logger"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
)

type MessagePayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

// Service prepares the interfaces related with this auth service
type Service struct {
	deviceSvc    *svc.Service
	log          *logger.Logger
	whatsAppBot  *botHook.WaManager
	BotClients   *botHook.BotClientList
	httpClient   *http.Client
	webhookUrl   string
	qrCodeDir    string
	echoMsg      bool
	wHookEnabled bool
}

// NewService creates a new auth service
func NewService(deviceSvc *svc.Service, log *logger.Logger,
	whatsAppBot *botHook.WaManager, httpClient *http.Client, webhookUrl, qrCodeDir string,
	echoMsg, wHookEnabled bool, bcList *botHook.BotClientList) *Service {

	return &Service{
		deviceSvc:    deviceSvc,
		log:          log,
		whatsAppBot:  whatsAppBot,
		httpClient:   httpClient,
		webhookUrl:   webhookUrl,
		qrCodeDir:    qrCodeDir,
		echoMsg:      echoMsg,
		wHookEnabled: wHookEnabled,
		BotClients:   bcList,
	}
}

// New creates a new session or reconnects an existing session
func (s *Service) New(ctx context.Context, phone string) error {
	// validates if phone exists in the database
	device, err := s.deviceSvc.GetDeviceByPhone(ctx, phone)
	if err != nil {
		return err
	}

	// lock it with a null value
	(*s.BotClients)[phone] = nil

	// run in background process
	go s.Process(phone, device)

	return nil
}

// Process processes the request as new session or reconnect the existing session
func (s *Service) Process(phone string, device svc.Device) {
	var err error
	var bot *botHook.WaBot
	var thisJID string

	if device.JID == "" {
		// creates new bot client
		s.log.Info("creating a new whatsapp session")
		bot, err = botHook.NewWhatsappClient(s.httpClient, s.webhookUrl, s.whatsAppBot.Container, s.log,
			phone, s.qrCodeDir, s.echoMsg, s.wHookEnabled)
		if err != nil {
			s.log.Warn("error create whatsapp client")
			delete(*s.BotClients, phone)
			return
		}

		thisJID = bot.Client.Store.ID.String()

		// updates JID and webhook from the device document
		err = s.deviceSvc.UpdateJID(context.Background(), thisJID, device.ID)
		if err != nil {
			s.log.Warn("failed to update JID information")
			delete(*s.BotClients, phone)
			return
		}
		s.log.Warn("finished updating the JID information")
	} else {
		// opens an existing session
		s.log.Info(fmt.Sprintf("reconnecting an existing whatsapp session with JID -> %s", device.JID))

		bot, err = botHook.LoginExistingWASession(s.httpClient, s.webhookUrl, s.whatsAppBot.Container, s.log,
			device.JID, phone, s.echoMsg, s.wHookEnabled)
		if err != nil {
			s.log.Warn(fmt.Sprintf("error create whatsapp client with an existing JID -> %s", device.JID))
			delete(*s.BotClients, phone)
			return
		}

		thisJID = device.JID
	}

	// registers event handler
	bot.Register()

	// add to client list
	(*s.BotClients)[phone] = bot

	// prints JID
	s.log.Info(fmt.Sprintf("captured JID -> %s", thisJID))
}

// Disconnect close the existing session
func (s *Service) Disconnect(phone string) string {
	var msg string

	// if key exists, disconnect and remove the key first
	if _, ok := (*s.BotClients)[phone]; ok {
		// get session client and disconnect it
		(*s.BotClients)[phone].Client.Disconnect()

		// removes from the map
		delete(*s.BotClients, phone)

		msg = fmt.Sprintf("session has been disconnected")
		s.log.Info(fmt.Sprintf("session [%s] has been disconnected", phone))
	} else {
		msg = fmt.Sprintf("session does not exists yet. do nothing")
		s.log.Info(fmt.Sprintf("session [%s] does not exists yet. do nothing", phone))
	}

	return msg
}

// SendTextMessage sends a text message
func (s *Service) SendTextMessage(payload MessagePayload) error {
	payload.sanitize()
	err := payload.validate()
	if err != nil {
		return err
	}

	// if device in From (=phone) does not exists, rejects
	if _, ok := (*s.BotClients)[payload.From]; !ok {
		return fmt.Errorf("no active session found for this device")
	} else {
		// if session is not ready yet, rejects
		if (*s.BotClients)[payload.From] == nil {
			return fmt.Errorf("session for this device is not ready yet")
		}

		// validates phone number and get the recipient
		recipient, err := (*s.BotClients)[payload.From].ValidateAndGetRecipient(payload.To, true)
		if err != nil {
			s.log.Error(fmt.Sprintf("phone [%s] got validation error(s)", payload.To), zap.Error(err))
			return fmt.Errorf("phone got validation error(s)")
		}

		// starts sending the message in a background
		go s.sendTextMessageInBackground(recipient, payload)
	}

	return nil
}

// sendTextMessageInBackground sends a text message in a background
func (s *Service) sendTextMessageInBackground(recipient *types.JID, payload MessagePayload) {
	err := (*s.BotClients)[payload.From].SendMsg(*recipient, payload.Message)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to send the message to [%s]", payload.To), zap.Error(err))
	}
}

// validate validates message payload
func (p *MessagePayload) validate() error {
	return nil
}

// sanitize sanitizes message payload
func (p *MessagePayload) sanitize() {
	// removes `+` symbol if exists
	if p.From[0:1] == "+" {
		_, i := utf8.DecodeRuneInString(p.From)
		p.From = p.From[i:]
	}

	// removes `+` symbol if exists
	if p.To[0:1] == "+" {
		_, i := utf8.DecodeRuneInString(p.To)
		p.To = p.To[i:]
	}
}
