package session

import (
	"context"
	"fmt"
	"go.mau.fi/whatsmeow/types"
	"go.uber.org/zap"
	"math/rand"
	"net/http"

	"github.com/ardihikaru/go-modules/pkg/logger"
	botHook "github.com/ardihikaru/go-modules/pkg/whatsappbot/wawebhook"
	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
)

// Service prepares the interfaces related with this auth service
type Service struct {
	deviceSvc    *svc.Service
	log          *logger.Logger
	whatsAppBot  *botHook.WaManager
	BotClients   *botHook.BotClientList
	httpClient   *http.Client
	imageDir     string
	qrCodeDir    string
	echoMsg      bool
	wHookEnabled bool
	qrToTerminal bool
}

// NewService creates a new auth service
func NewService(deviceSvc *svc.Service, log *logger.Logger,
	whatsAppBot *botHook.WaManager, httpClient *http.Client, imageDir, qrCodeDir string,
	echoMsg, wHookEnabled, qrToTerminal bool, bcList *botHook.BotClientList) *Service {

	return &Service{
		deviceSvc:    deviceSvc,
		log:          log,
		whatsAppBot:  whatsAppBot,
		httpClient:   httpClient,
		imageDir:     imageDir,
		qrCodeDir:    qrCodeDir,
		echoMsg:      echoMsg,
		wHookEnabled: wHookEnabled,
		qrToTerminal: qrToTerminal,
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
	go s.Process(ctx, phone, device)

	return nil
}

// Process processes the request as new session or reconnect the existing session
func (s *Service) Process(ctx context.Context, phone string, device svc.Device) {
	var err error
	var bot *botHook.WaBot
	var thisJID string

	// fetches device document
	deviceDoc, err := s.deviceSvc.GetDeviceByPhone(ctx, phone)
	if err != nil {
		s.log.Warn(fmt.Sprintf("failed tp fetch device document by phone [%s]", phone))
		delete(*s.BotClients, phone)
		return
	}

	if device.JID == "" {
		// creates new bot client
		s.log.Info("creating a new whatsapp session")
		bot, err = botHook.NewWhatsappClient(s.httpClient, deviceDoc.WebhookUrl, s.imageDir, s.whatsAppBot.Container, s.log,
			phone, s.qrCodeDir, s.echoMsg, s.wHookEnabled, s.qrToTerminal)
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

		bot, err = botHook.LoginExistingWASession(s.httpClient, deviceDoc.WebhookUrl, s.imageDir, s.whatsAppBot.Container,
			s.log, device.JID, phone, s.echoMsg, s.wHookEnabled)
		if err != nil {
			s.log.Warn(fmt.Sprintf("error create whatsapp client with an existing JID -> %s", device.JID),
				zap.Error(err))
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
func (s *Service) SendTextMessage(payload botHook.MessagePayload) error {
	payload.Sanitize()
	err := payload.Validate()
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

// SendImageMessage sends ann image-based message
func (s *Service) SendImageMessage(payload botHook.MessagePayload) error {
	payload.Sanitize()
	err := payload.Validate()
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
		go s.sendImageMessageInBackground(recipient, payload)
	}

	return nil
}

// sendTextMessageInBackground sends a text message in a background
func (s *Service) sendTextMessageInBackground(recipient *types.JID, payload botHook.MessagePayload) {
	err := (*s.BotClients)[payload.From].SendMsg(*recipient, payload.Message)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to send the message to [%s]", payload.To), zap.Error(err))
	}
}

// sendImageMessageInBackground sends an image-based message in a background
func (s *Service) sendImageMessageInBackground(recipient *types.JID, payload botHook.MessagePayload) {
	var err error

	// builds image full path
	imgPath := fmt.Sprintf("%s/%s", s.imageDir, payload.ImageFileName)

	// uploads to whatsapp server
	imgInBytes, uploaded, err := (*s.BotClients)[payload.From].UploadImgToWhatsapp(imgPath)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to upload file (=%s) to Whatsapp server", payload.ImageFileName), zap.Error(err))
		return
	}

	// prepares image information
	contentType := http.DetectContentType(*imgInBytes)
	fileLength := uint64(len(*imgInBytes))

	// sends image message to whatsapp
	err = (*s.BotClients)[payload.From].SendImgMsg(*recipient, uploaded, payload.ImageCaption, contentType, fileLength)
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to send the image message to [%s]", payload.To), zap.Error(err))
	}
}

// IsOnWhatsapp verify if this phone number on Whatsapp or not
func (s *Service) IsOnWhatsapp(phone string) (bool, error) {
	var err error

	// picks one random active client session
	clientPhone := s.getRandomPhoneAsClient()
	if clientPhone == nil {
		return false, fmt.Errorf("no active session to be used")
	}

	phones := buildValidatedPhone(phone)
	onWhatsapp, err := (*s.BotClients)[*clientPhone].Client.IsOnWhatsApp(phones)
	if err != nil {
		s.log.Error("failed to check on the Whatsapp Server", zap.Error(err))
		return false, err
	}
	s.log.Debug(fmt.Sprintf("%v", onWhatsapp))

	return onWhatsapp[0].IsIn, nil
}

func (s *Service) getRandomPhoneAsClient() *string {
	// returns nil if no active session found
	totalSession := len(*s.BotClients)
	if totalSession == 0 {
		return nil
	}

	randIter := rand.Intn(totalSession-0) + 0
	for phone, _ := range *s.BotClients {
		if randIter == 0 {
			return &phone
		}

		randIter -= 1
	}

	return nil
}

func buildValidatedPhone(phone string) []string {
	phones := make([]string, 1)

	// enriches with `+` symbol if missing
	if phone[0:1] != "+" {
		phone = fmt.Sprintf("+%s", phone)
	}

	phones[0] = phone

	return phones
}
