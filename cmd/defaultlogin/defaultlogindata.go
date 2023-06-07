package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/satumedishub/sea-cucumber-api-service/internal/app"
	"github.com/satumedishub/sea-cucumber-api-service/internal/config"
	"github.com/satumedishub/sea-cucumber-api-service/internal/logger"
	accSvc "github.com/satumedishub/sea-cucumber-api-service/internal/service/account"
	userSvc "github.com/satumedishub/sea-cucumber-api-service/internal/service/user"
	s "github.com/satumedishub/sea-cucumber-api-service/internal/storage"
)

// main run the main application
func main() {
	// loads configuration
	cfg, err := config.Get()
	if err != nil {
		app.FatalOnError(err, "error loading configuration")
	}

	// configures logger
	log, err := logger.New(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		app.FatalOnError(err, "failed to prepare the logger")
	}

	// initializes persistent store
	db := app.InitializeDB(cfg, log)

	// counts data
	accountCount := countAccountData(db, cfg.DBName, log)

	// shows the build version
	log.Info("assessment summary. ",
		zap.Int("accountCount", accountCount),
	)

	// if no data found, create one
	if accountCount == 0 {
		createAccountData(db, log)
	} else {
		log.Info("nothing todo. default account has been created -> user/pass: administrator / administrator")
	}

}

// countAccountData counts total records from the accounts collection
func countAccountData(db *s.DataStoreMongo, dbName string, log *logger.Logger) int {
	accountCollection := db.Client.Database(dbName).Collection(s.AccountCollection)

	accountCount, err := accountCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		fatalOnError(err, "failed to count account data", log)
	}

	return int(accountCount)
}

// createAccountData creates a new account
func createAccountData(db *s.DataStoreMongo, log *logger.Logger) {
	var err error

	// first, create a new user
	userId := createUserData(db, log)

	// hash the password
	hashedPassword, err := accSvc.HashPassword("administrator")
	if err != nil {
		fatalOnError(err, "failed to hash the provided password", log)
	}

	accDoc := accSvc.Account{
		UserId:   userId.Hex(),
		Username: "administrator",
		Password: hashedPassword,
	}

	_, err = db.InsertAccount(context.Background(), accDoc)
	if err != nil {
		fatalOnError(err, "failed to create a new account data", log)
	}

	log.Info("new account has been created -> user/pass: administrator / administrator")
}

// createUserData creates a new user
func createUserData(db *s.DataStoreMongo, log *logger.Logger) primitive.ObjectID {
	userDoc := userSvc.User{
		Name:       "administrator",
		Email:      "administrator@email.com",
		Role:       "SUPER_ADMIN",
		Contact:    "085608560856",
		Background: "a background",
	}

	user, err := db.InsertUser(context.Background(), userDoc)
	if err != nil {
		fatalOnError(err, "failed to create a new account data", log)
	}

	return user.ID
}

// fatalOnError stops the application if any fatal error occurs
func fatalOnError(err error, msg string, log *logger.Logger) {
	if err != nil {
		log.Error("fatal error.", zap.Error(err))
		zap.S().Fatalf("%s:%s", msg, err)
	}
}
