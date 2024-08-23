package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStore() *MongoStore {
	return &MongoStore{}
}

func (s *MongoStore) Init() {
}

func (s *MongoStore) Stop() {
	if err := s.client.Disconnect(context.Background()); err != nil {
		log.Error("Error while disconnecting from MongoDB", "error", err)
	}
	log.Info("🛑 MongoDB connection closed")
}

func (s *MongoStore) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	connStr := fmt.Sprintf(
		"%s:%s/%s",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return err
	}
	log.SetTimeFormat(time.Kitchen)
	log.Info("🗄️ DB Connected")
	s.client = client
	return nil
}

func (s *MongoStore) GetDatabaseName() *mongo.Database {
	s.db = s.client.Database(viper.GetString("database.dbname"))
	return s.db
}
