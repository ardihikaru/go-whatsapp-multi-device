package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	svc "github.com/ardihikaru/go-whatsapp-multi-device/internal/service/account"

	"github.com/ardihikaru/go-whatsapp-multi-device/pkg/utils/httputils"
)

type AccountsID string
type AccountsUserID string
type AccountsName string
type AccountsEmail string
type AccountsRole string
type AccountsContact string
type AccountsBackground string
type AccountsUsername string
type AccountsPassword string
type AccountsLastLogin string
type AccountsCreatedAt string
type AccountsUpdatedAt string

const (
	// AccountCollection defines the collection name
	AccountCollection = "accounts"

	// FnAccountsId defines the main identifier to query the user document
	FnAccountsId = AccountsID("_id")

	// FnAccountsUserId defines the main identifier to query the user document
	FnAccountsUserId = AccountsUserID("user_id")

	// FnAccountsEmail defines the email key in the user document
	FnAccountsEmail = AccountsEmail("email")

	// FnAccountsUsername defines the username key in the user document
	FnAccountsUsername = AccountsUsername("username")

	// FnAccountsPassword defines the password key in the user document
	FnAccountsPassword = AccountsPassword("password")

	// FnAccountsLastLogin defines the last login key in the user document
	FnAccountsLastLogin = AccountsLastLogin("last_login")

	// FnAccountsCreatedAt defines the createdAt key in the user document
	FnAccountsCreatedAt = AccountsCreatedAt("_c")

	// FnAccountsUpdatedAt defines the updatedAt key in the user document
	FnAccountsUpdatedAt = AccountsUpdatedAt("_m")
)

// AccountDoc is the document prepared for the captured user information
type AccountDoc struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	LastLogin primitive.DateTime `bson:"last_login"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

// ToService converts the userDoc struct into User struct from the service
func (u *AccountDoc) ToService() svc.Account {
	return svc.Account{
		ID:        u.ID,
		UserId:    u.UserId.Hex(),
		Username:  u.Username,
		Password:  u.Password,
		LastLogin: u.LastLogin.Time(),
		CreatedAt: u.CreatedAt.Time(),
		UpdatedAt: u.UpdatedAt.Time(),
	}
}

// AccToBsonObject converts the accountDoc struct into account struct from the service
func AccToBsonObject(u svc.Account) (AccountDoc, error) {
	objID, err := primitive.ObjectIDFromHex(u.UserId)
	if err != nil {
		return AccountDoc{}, err
	}

	return AccountDoc{
		ID:        primitive.NewObjectID(),
		UserId:    objID,
		Username:  u.Username,
		Password:  u.Password,
		LastLogin: primitive.NewDateTimeFromTime(u.LastLogin),
		CreatedAt: primitive.NewDateTimeFromTime(u.CreatedAt),
		UpdatedAt: primitive.NewDateTimeFromTime(u.UpdatedAt),
	}, nil
}

// GetAccountByID fetch account data by ID
func (d *DataStoreMongo) GetAccountByID(ctx context.Context, accountID primitive.ObjectID,
	ignorePasswd bool) (svc.Account, error) {
	// prepares the options
	var opts *options.FindOneOptions = options.FindOne()

	// prepares the projection
	if ignorePasswd {
		projection := bson.D{
			{Key: string(FnAccountsPassword), Value: 0},
		}
		opts.SetProjection(projection)
	}

	// prepares the filter
	filter := bson.D{primitive.E{Key: string(FnAccountsId), Value: accountID}}

	// finds document by ID and convert the cursor result to bson object
	accountDoc := AccountDoc{}
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)
	err := accountsCollection.FindOne(ctx, filter, opts).Decode(&accountDoc)
	if err != nil {
		return svc.Account{}, fmt.Errorf("cannot find user: %w", err)
	}

	return accountDoc.ToService(), nil
}

// buildAccountFilterOption builds an option for filtering document
func buildAccountFilterOption(search string) bson.D {
	var filter bson.D
	if search != "" {
		regex := bson.D{{"$regex", primitive.Regex{Pattern: "^.*" + search + ".*$", Options: "i"}}}

		filter = bson.D{
			{"$or", bson.A{
				bson.M{string(FnAccountsEmail): regex},
				bson.M{string(FnAccountsUsername): regex},
			}},
		}
	} else {
		filter = bson.D{}
	}

	return filter
}

// GetAccounts fetch account data by custom query
func (d *DataStoreMongo) GetAccounts(ctx context.Context, params httputils.GetQueryParams, ignorePasswd bool) (int64, []svc.Account, error) {
	// prepares the options
	var opts *options.FindOptions = options.Find()

	// set query parameters
	opts.SetLimit(params.Limit)
	opts.SetSkip(params.Offset)

	// sets order option
	if params.Sort == ID {
		opts.SetSort(buildOrderOption(params.Order, params.Sort))
	}

	// builds filter
	filter := buildAccountFilterOption(params.Search)

	// prepares the projection
	if ignorePasswd {
		projection := bson.D{
			{Key: string(FnAccountsPassword), Value: 0},
		}
		opts.SetProjection(projection)
	}

	// gets cursor
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)
	total, err := accountsCollection.CountDocuments(ctx, bson.M{})
	cur, err := accountsCollection.Find(ctx, filter, opts)

	if err != nil {
		return 0, nil, fmt.Errorf("cannot find Accounts: %w", err)
	}
	defer cur.Close(ctx)

	res := make([]svc.Account, 0)
	for cur.Next(ctx) {
		accDoc := AccountDoc{}

		err = cur.Decode(&accDoc)
		if err != nil {
			return 0, nil, fmt.Errorf("cannot decode user: %w", err)
		}

		res = append(res, accDoc.ToService())
	}

	return total, res, nil
}

// GetAccountByUsername fetch account data by username
func (d *DataStoreMongo) GetAccountByUsername(ctx context.Context, username string, ignorePasswd bool) (svc.Account, error) {
	// prepares the options
	var opts *options.FindOneOptions = options.FindOne()

	// prepares the projection
	if ignorePasswd {
		projection := bson.D{
			{Key: string(FnAccountsPassword), Value: 0},
		}
		opts.SetProjection(projection)
	}

	// prepares the filter
	filter := bson.D{primitive.E{Key: string(FnAccountsUsername), Value: username}}

	// finds document by ID and convert the cursor result to bson object
	accDoc := AccountDoc{}
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)
	err := accountsCollection.FindOne(ctx, filter, opts).Decode(&accDoc)
	if err != nil {
		return svc.Account{}, fmt.Errorf("cannot find account: %w", err)
	}

	return accDoc.ToService(), nil
}

// GetAccountByEmail fetch account data by email
func (d *DataStoreMongo) GetAccountByEmail(ctx context.Context, email string, ignorePasswd bool) (svc.Account, error) {
	// prepares the options
	var opts *options.FindOneOptions = options.FindOne()

	// prepares the projection
	if ignorePasswd {
		projection := bson.D{
			{Key: string(FnAccountsPassword), Value: 0},
		}
		opts.SetProjection(projection)
	}

	// prepares the filter
	filter := bson.D{primitive.E{Key: string(FnAccountsEmail), Value: email}}

	// finds document by ID and convert the cursor result to bson object
	accDoc := AccountDoc{}
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)
	err := accountsCollection.FindOne(ctx, filter, opts).Decode(&accDoc)
	if err != nil {
		return svc.Account{}, fmt.Errorf("cannot find account: %w", err)
	}

	return accDoc.ToService(), nil
}

// InsertAccount stores account data
func (d *DataStoreMongo) InsertAccount(ctx context.Context, doc svc.Account) (svc.Account, error) {
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)

	// build account document
	accDoc, err := AccToBsonObject(doc)
	if err != nil {
		return doc, fmt.Errorf("failed to convert userId: %w", err)
	}

	// finds document by ID and convert the cursor result to bson object
	insertResult, err := accountsCollection.InsertOne(ctx, accDoc)
	if err != nil {
		return doc, fmt.Errorf("cannot find account: %w", err)
	}

	// enrich with _id
	doc.ID = insertResult.InsertedID.(primitive.ObjectID)

	return doc, nil
}

// UpdatePassword updates password
func (d *DataStoreMongo) UpdatePassword(ctx context.Context, doc svc.Account, accountId string) (svc.Account, error) {
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)

	objID, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return svc.Account{}, err
	}

	// builds filter
	filter := bson.D{{Key: string(FnAccountsId), Value: objID}}

	// prepares document to update
	docUpdated := bson.D{
		{Key: string(FnAccountsPassword), Value: doc.Password},
	}
	docBson := bson.D{
		{Key: "$set", Value: docUpdated},
	}

	// build account document
	accDoc, err := AccToBsonObject(doc)
	if err != nil {
		return svc.Account{}, fmt.Errorf("failed to convert userId: %w", err)
	}

	// finds document by ID and convert the cursor result to bson object
	_, err = accountsCollection.UpdateOne(ctx, filter, docBson)

	if err != nil {
		return svc.Account{}, err
	}

	// enrich with _id
	accDoc.ID = objID

	return accDoc.ToService(), nil
}

// UpdateLastLogin updates last login information
func (d *DataStoreMongo) UpdateLastLogin(ctx context.Context, doc svc.Account, accountId string) (svc.Account, error) {
	accountsCollection := d.Client.Database(d.DBName).Collection(AccountCollection)

	objID, err := primitive.ObjectIDFromHex(accountId)
	if err != nil {
		return svc.Account{}, err
	}

	// builds filter
	filter := bson.D{{Key: string(FnAccountsId), Value: objID}}

	// prepares document to update
	docUpdated := bson.D{
		{Key: string(FnAccountsLastLogin), Value: doc.LastLogin},
	}
	docBson := bson.D{
		{Key: "$set", Value: docUpdated},
	}

	// finds document by ID and convert the cursor result to bson object
	_, err = accountsCollection.UpdateOne(ctx, filter, docBson)

	if err != nil {
		return svc.Account{}, err
	}

	return doc, nil
}
