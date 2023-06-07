package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	svc "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"

	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/httputils"
)

type UsersUserID string
type UsersName string
type UsersEmail string
type UsersRole string
type UsersContact string
type UsersBackground string

type UsersLastLogin string
type UsersCreatedAt string
type UsersUpdatedAt string

const (
	// UserCollection defines the collection name
	UserCollection = "users"

	// FnUsersUserId defines the main identifier to query the user document
	FnUsersUserId = UsersUserID("_id")

	// FnUsersRole defines the role key in the user document
	FnUsersRole = UsersRole("role")

	// FnUsersContact defines the contact key in the user document
	FnUsersContact = UsersContact("contact")

	// FnUsersBackground defines the background key in the user document
	FnUsersBackground = UsersBackground("background")

	// FnUsersEmail defines the email key in the user document
	FnUsersEmail = UsersEmail("email")

	// FnUsersName defines the name key in the user document
	FnUsersName = UsersName("name")
)

// UserDoc is the document prepared for the captured user information
type UserDoc struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Email      string             `bson:"email"`
	Role       string             `bson:"role"`
	Contact    string             `bson:"contact"`
	Background string             `bson:"background"`
	CreatedAt  primitive.DateTime `bson:"created_at"`
	UpdatedAt  primitive.DateTime `bson:"updated_at"`
}

// ToService converts the userDoc struct into User struct from the service
func (u *UserDoc) ToService() svc.User {
	return svc.User{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		Contact:    u.Contact,
		Background: u.Background,
		CreatedAt:  u.CreatedAt.Time(),
		UpdatedAt:  u.UpdatedAt.Time(),
	}
}

// UserToBsonObject converts the userDoc struct into User struct from the service
func UserToBsonObject(u svc.User) UserDoc {
	return UserDoc{
		ID:         primitive.NewObjectID(),
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		Contact:    u.Contact,
		Background: u.Background,
		CreatedAt:  primitive.NewDateTimeFromTime(u.CreatedAt),
		UpdatedAt:  primitive.NewDateTimeFromTime(u.UpdatedAt),
	}
}

// GetUserByID fetch user data by ID
func (d *DataStoreMongo) GetUserByID(ctx context.Context, userID primitive.ObjectID) (svc.User, error) {
	// prepares the options
	var opts *options.FindOneOptions = options.FindOne()

	// prepares the filter
	filter := bson.D{primitive.E{Key: string(FnUsersUserId), Value: userID}}

	// finds document by ID and convert the cursor result to bson object
	usrDoc := UserDoc{}
	usersCollection := d.Client.Database(d.DBName).Collection(UserCollection)
	err := usersCollection.FindOne(ctx, filter, opts).Decode(&usrDoc)
	if err != nil {
		return svc.User{}, fmt.Errorf("cannot find user: %w", err)
	}

	return usrDoc.ToService(), nil
}

// buildUserFilterOption builds an option for filtering document
func buildUserFilterOption(search string) bson.D {
	var filter bson.D
	if search != "" {
		regex := bson.D{{"$regex", primitive.Regex{Pattern: "^.*" + search + ".*$", Options: "i"}}}

		filter = bson.D{
			{"$or", bson.A{
				bson.M{string(FnUsersRole): regex},
				bson.M{string(FnUsersContact): regex},
				bson.M{string(FnUsersBackground): regex},
				bson.M{string(FnUsersEmail): regex},
				bson.M{string(FnUsersName): regex},
			}},
		}
	} else {
		filter = bson.D{}
	}

	return filter
}

// GetUsers fetch user data by custom query
func (d *DataStoreMongo) GetUsers(ctx context.Context, params httputils.GetQueryParams) (int64, []svc.User, error) {
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
	filter := buildUserFilterOption(params.Search)

	// gets cursor
	usersCollection := d.Client.Database(d.DBName).Collection(UserCollection)
	total, err := usersCollection.CountDocuments(ctx, bson.M{})
	cur, err := usersCollection.Find(ctx, filter, opts)

	if err != nil {
		return 0, nil, fmt.Errorf("cannot find users: %w", err)
	}
	defer cur.Close(ctx)

	res := make([]svc.User, 0)
	for cur.Next(ctx) {
		usrDoc := UserDoc{}

		err = cur.Decode(&usrDoc)
		if err != nil {
			return 0, nil, fmt.Errorf("cannot decode user: %w", err)
		}

		res = append(res, usrDoc.ToService())
	}

	return total, res, nil
}

// GetUserByEmail fetch user data by email
func (d *DataStoreMongo) GetUserByEmail(ctx context.Context, email string) (svc.User, error) {
	// prepares the options
	var opts *options.FindOneOptions = options.FindOne()

	// prepares the filter
	filter := bson.D{primitive.E{Key: string(FnUsersEmail), Value: email}}

	// finds document by ID and convert the cursor result to bson object
	usrDoc := UserDoc{}
	usersCollection := d.Client.Database(d.DBName).Collection(UserCollection)
	err := usersCollection.FindOne(ctx, filter, opts).Decode(&usrDoc)
	if err != nil {
		return svc.User{}, fmt.Errorf("cannot find user: %w", err)
	}

	return usrDoc.ToService(), nil
}

// InsertUser stores user data
func (d *DataStoreMongo) InsertUser(ctx context.Context, doc svc.User) (svc.User, error) {
	usersCollection := d.Client.Database(d.DBName).Collection(UserCollection)

	// build user document
	usrDoc := UserToBsonObject(doc)

	// finds document by ID and convert the cursor result to bson object
	insertResult, err := usersCollection.InsertOne(ctx, usrDoc)
	if err != nil {
		return doc, fmt.Errorf("cannot find user: %w", err)
	}

	// enrich with _id
	doc.ID = insertResult.InsertedID.(primitive.ObjectID)

	return doc, nil
}

// UpdateUser updates user data
func (d *DataStoreMongo) UpdateUser(ctx context.Context, doc svc.User, userId string) (svc.User, error) {
	usersCollection := d.Client.Database(d.DBName).Collection(UserCollection)

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return svc.User{}, err
	}

	// builds filter
	filter := bson.D{{Key: string(FnUsersUserId), Value: objID}}

	// prepares document to update
	docUpdated := bson.D{
		{Key: string(FnUsersName), Value: doc.Name},
		{Key: string(FnUsersRole), Value: doc.Role},
		{Key: string(FnUsersContact), Value: doc.Contact},
		{Key: string(FnUsersBackground), Value: doc.Background},
	}
	docBson := bson.D{
		{Key: "$set", Value: docUpdated},
	}

	// build user document
	usrDoc := UserToBsonObject(doc)

	// finds document by ID and convert the cursor result to bson object
	_, err = usersCollection.UpdateOne(ctx, filter, docBson)

	if err != nil {
		return svc.User{}, err
	}

	// enrich with _id
	usrDoc.ID = objID

	//return user, nil
	return usrDoc.ToService(), nil
}
