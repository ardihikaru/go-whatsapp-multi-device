// Package storage provides operations to deal directly with the database
package storage

import (
	"context"
	"github.com/satumedishub/sea-cucumber-api-service/pkg/utils/query"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	mopts "go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// mgClient represents the mongo client
	// it uses an official mongo supported driver
	mgClient *mongo.Client
)

const (
	// ID represents id without underscore
	ID = "id"
)

type Order int // for sorting documents

type E struct {
	Key   string
	Value interface{}
}

// DataStoreMongoConfig sets up the required parameters to build a Mongo Data Store
type DataStoreMongoConfig struct {
	// ConnectionString is the connection string
	ConnectionString string

	// DatabaseName is the name of the database
	DatabaseName string

	// ConnectionTimeout is the database timeout interval
	ConnectionTimeout time.Duration

	// HeartBeatInterval is the heartbeat interval
	HeartBeatInterval time.Duration

	// LocalThreshold is the local threshold
	LocalThreshold time.Duration

	// ServerSelectionTimeout is the server selection timeout
	ServerSelectionTimeout time.Duration

	// MaxPoolSize is the maximum size of the connection pool
	MaxPoolSize uint64
}

// DataStoreMongo describes the resource of the data store
type DataStoreMongo struct {
	DBName string
	Client *mongo.Client
}

// NewDataStoreMongo creates a MongoDB client with a valid database string
// This function is called ONCE on the service start and
//
//	then pass the DataStoreMongo object around to orchestrator, service layers and repository layers.
func NewDataStoreMongo(config DataStoreMongoConfig) (*DataStoreMongo, error) {
	//init master session
	var err error

	// prepares the mongo client options
	clientOptions := mopts.Client().ApplyURI(config.ConnectionString)
	clientOptions.SetConnectTimeout(config.ConnectionTimeout)
	clientOptions.SetHeartbeatInterval(config.HeartBeatInterval)
	clientOptions.SetLocalThreshold(config.LocalThreshold)
	clientOptions.SetMaxPoolSize(config.MaxPoolSize)
	clientOptions.SetServerSelectionTimeout(config.ServerSelectionTimeout)

	err = clientOptions.Validate()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	mgClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if mgClient == nil {
		return nil, errors.New("Unable to initialize the Mongo Client")
	}

	// from: https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial
	/*
		It is best practice to keep a client that is connected to MongoDB around so that the
		application can make use of connection  	 - you don't want to open and close a
		connection for each query. However, if your application no longer requires a connection,
		the connection can be closed with client.Disconnect() like so:
	*/
	err = mgClient.Ping(ctx, nil)
	if err != nil {
		mgClient = nil
		return nil, err
	}

	if mgClient == nil {
		return nil, errors.New("failed to open mongo-driver session")
	}
	db := &DataStoreMongo{
		DBName: config.DatabaseName,
		Client: mgClient,
	}

	return db, nil
}

// buildOrderOption builds an option for sorting document
func buildOrderOption(orderQuery, sort string) E {
	var order int
	if orderQuery == query.ASC {
		order = 1
	} else {
		order = -1
	}

	if sort == ID {
		sort = "_id"
	}

	return E{
		Key:   sort,
		Value: order,
	}
}
