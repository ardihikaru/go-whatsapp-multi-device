package storage

import (
	"context"
	"fmt"

	"github.com/ardihikaru/go-modules/pkg/utils/httputils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/device"
)

const (
	// DeviceCollection defines the collection name
	DeviceCollection = "devices"

	// FnDevicesId defines the main identifier that acts as a Primary Key
	FnDevicesId = string("_id")

	// FnDevicesJID defines the JID of the registered device
	FnDevicesJID = string("jid")

	// FnDevicesPhone defines the phone number
	FnDevicesPhone = string("phone")

	// FnDevicesName defines the name of the device's owner
	FnDevicesName = string("name")

	// FnDevicesWebhookUrl defines the Webhook URL
	FnDevicesWebhookUrl = string("webhook_url")

	// FnDevicesCreatedAt defines the creation time
	FnDevicesCreatedAt = string("_c")

	// FnDevicesUpdatedAt defines the update time
	FnDevicesUpdatedAt = string("_m")
)

// DeviceDoc is the document prepared for the captured user information
type DeviceDoc struct {
	ID         primitive.ObjectID `bson:"_id"`
	JID        string             `bson:"jid"`
	Phone      string             `bson:"phone"`
	Name       string             `bson:"name"`
	WebhookUrl string             `bson:"webhook_url"`
	CreatedAt  primitive.DateTime `bson:"created_at"`
	UpdatedAt  primitive.DateTime `bson:"updated_at"`
}

// ToService converts the DeviceDoc struct into Device struct
func (u *DeviceDoc) ToService() svc.Device {
	return svc.Device{
		ID:         u.ID.Hex(),
		JID:        u.JID,
		Phone:      u.Phone,
		Name:       u.Name,
		WebhookUrl: u.WebhookUrl,
		CreatedAt:  u.CreatedAt.Time(),
		UpdatedAt:  u.UpdatedAt.Time(),
	}
}

// deviceToBsonObject converts the accountDoc struct into account struct from the service
func deviceToBsonObject(u svc.Device) (DeviceDoc, error) {
	return DeviceDoc{
		ID:         primitive.NewObjectID(),
		JID:        u.JID,
		Phone:      u.Phone,
		Name:       u.Name,
		WebhookUrl: u.WebhookUrl,
		CreatedAt:  primitive.NewDateTimeFromTime(u.CreatedAt),
		UpdatedAt:  primitive.NewDateTimeFromTime(u.UpdatedAt),
	}, nil
}

// GetDeviceByID fetch device data by ID
func (d *DataStoreMongo) GetDeviceByID(ctx context.Context, id string) (svc.Device, error) {
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	IdObject, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return svc.Device{}, err
	}

	// prepares the options
	var opts = options.FindOne()

	// prepares the filter
	filter := bson.D{primitive.E{Key: FnDevicesId, Value: IdObject}}

	// finds document by ID and convert the cursor result to bson object
	doc := DeviceDoc{}
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)
	err = collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		return svc.Device{}, fmt.Errorf("cannot find device: %w", err)
	}

	return doc.ToService(), nil
}

// GetDeviceByPhone fetch device data by phone
func (d *DataStoreMongo) GetDeviceByPhone(ctx context.Context, phone string) (svc.Device, error) {
	// prepares the options
	var opts = options.FindOne()

	// prepares the filter
	filter := bson.D{primitive.E{Key: FnDevicesPhone, Value: phone}}

	// finds document by ID and convert the cursor result to bson object
	doc := DeviceDoc{}
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)
	err := collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		return svc.Device{}, fmt.Errorf("cannot find device: %w", err)
	}

	return doc.ToService(), nil
}

// GetDeviceByJID fetch device data by JID
func (d *DataStoreMongo) GetDeviceByJID(ctx context.Context, jid string) (svc.Device, error) {
	// prepares the options
	var opts = options.FindOne()

	// prepares the filter
	filter := bson.D{primitive.E{Key: FnDevicesJID, Value: jid}}

	// finds document by ID and convert the cursor result to bson object
	doc := DeviceDoc{}
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)
	err := collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		return svc.Device{}, fmt.Errorf("cannot find device: %w", err)
	}

	return doc.ToService(), nil
}

// reformatPhoneQuery removes + symbol on the first char (for phone number)
func reformatPhoneQuery(search string) string {
	if len(search) > 0 && search[0:1] == "+" {
		search = search[1:len(search)]
	}

	return search
}

// buildDeviceFilterOption builds an option for filtering document
func buildDeviceFilterOption(search string) bson.D {
	var filter bson.D

	// remove + symbol on the first char (for phone number)
	search = reformatPhoneQuery(search)

	if search != "" {
		regex := bson.D{{"$regex", primitive.Regex{Pattern: "^.*" + search + ".*$", Options: "i"}}}

		filter = bson.D{
			{"$or", bson.A{
				bson.M{FnDevicesJID: regex},
				bson.M{FnDevicesPhone: regex},
				bson.M{FnDevicesName: regex},
			}},
		}
	} else {
		filter = bson.D{}
	}

	return filter
}

// GetDevices fetch device data by custom query
func (d *DataStoreMongo) GetDevices(ctx context.Context, params httputils.GetQueryParams) (int64, []svc.Device, error) {
	// prepares the options
	var opts = options.Find()

	// set query parameters
	opts.SetLimit(params.Limit)
	opts.SetSkip(params.Offset)

	// sets order option
	if params.Sort == ID {
		opts.SetSort(buildOrderOption(params.Order, params.Sort))
	}

	// builds filter
	filter := buildDeviceFilterOption(params.Search)

	// gets cursor
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)
	total, err := collection.CountDocuments(ctx, bson.M{})
	cur, err := collection.Find(ctx, filter, opts)

	if err != nil {
		return 0, nil, fmt.Errorf("cannot find any device: %w", err)
	}
	defer cur.Close(ctx)

	res := make([]svc.Device, 0)
	for cur.Next(ctx) {
		doc := DeviceDoc{}

		err = cur.Decode(&doc)
		if err != nil {
			return 0, nil, fmt.Errorf("cannot decode device doc: %w", err)
		}

		res = append(res, doc.ToService())
	}

	return total, res, nil
}

// InsertDevice stores device data
func (d *DataStoreMongo) InsertDevice(ctx context.Context, doc svc.Device) (svc.Device, error) {
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)

	// build document
	accDoc, err := deviceToBsonObject(doc)
	if err != nil {
		return doc, fmt.Errorf("bson object convertion failed: %w", err)
	}

	// finds document by ID and convert the cursor result to bson object
	insertResult, err := collection.InsertOne(ctx, accDoc)
	if err != nil {
		return doc, fmt.Errorf("cannot find account: %w", err)
	}

	// enrich with _id
	doc.ID = insertResult.InsertedID.(primitive.ObjectID).Hex()

	return doc, nil
}

// UpdateJID updates device data
func (d *DataStoreMongo) UpdateJID(ctx context.Context, jid, id string) error {
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// builds filter
	filter := bson.D{{Key: FnDevicesId, Value: objID}}

	// prepares document to update
	doc := bson.D{
		{Key: FnDevicesJID, Value: jid},
	}
	docBson := bson.D{
		{Key: "$set", Value: doc},
	}

	// finds document by ID and executes update action
	_, err = collection.UpdateOne(ctx, filter, docBson)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDeviceName updates device name
func (d *DataStoreMongo) UpdateDeviceName(ctx context.Context, id, deviceName string) error {
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// builds filter
	filter := bson.D{{Key: FnDevicesId, Value: objID}}

	// prepares document to update
	doc := bson.D{
		{Key: FnDevicesName, Value: deviceName},
	}
	docBson := bson.D{
		{Key: "$set", Value: doc},
	}

	// finds document by ID and executes update action
	_, err = collection.UpdateOne(ctx, filter, docBson)
	if err != nil {
		return err
	}

	return nil
}

// UpdateWebhook updates webhook
func (d *DataStoreMongo) UpdateWebhook(ctx context.Context, id, webhook string) error {
	collection := d.Client.Database(d.DBName).Collection(DeviceCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// builds filter
	filter := bson.D{{Key: FnDevicesId, Value: objID}}

	// prepares document to update
	doc := bson.D{
		{Key: FnDevicesWebhookUrl, Value: webhook},
	}
	docBson := bson.D{
		{Key: "$set", Value: doc},
	}

	// finds document by ID and executes update action
	_, err = collection.UpdateOne(ctx, filter, docBson)
	if err != nil {
		return err
	}

	return nil
}
