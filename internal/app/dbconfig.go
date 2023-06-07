package app

import (
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/config"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/logger"
	"github.com/ardihikaru/go-whatsapp-multi-device/internal/storage"
)

func InitializeDB(cfg *config.Config, log *logger.Logger) *storage.DataStoreMongo {
	// initializes persistent store
	db, err := storage.NewDataStoreMongo(MakeDataStoreConfig(cfg))
	if err != nil {
		FatalOnError(err, "failed to connect to db")
	}

	return db
}

// MakeDataStoreConfig builds database config
// the storage object will keep this information for the further usages
func MakeDataStoreConfig(cfg *config.Config) storage.DataStoreMongoConfig {
	return storage.DataStoreMongoConfig{
		ConnectionString:       cfg.DbConnURI,
		DatabaseName:           cfg.DBName,
		ConnectionTimeout:      cfg.DbConnTimeout,
		HeartBeatInterval:      cfg.DbHeartBeatInterval,
		LocalThreshold:         cfg.DbLocalThreshold,
		ServerSelectionTimeout: cfg.DbServerSelTimeout,
		MaxPoolSize:            cfg.DbMaxPoolSize,
	}
}
