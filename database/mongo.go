package database

import (
	"context"
	"go_tgbot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func ConnectMongoDB() {
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(config.SetConfig.MongoUrl))
	MongoDB = client.Database("tgbot_marika")
}
