package store

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OpenSession to mongodb
func GetMongoClient(mongoCfg config.MongoConfig) *mongo.Client {
	conURL := "mongodb://" + mongoCfg.Host + ":" + strconv.Itoa(mongoCfg.Port)
	clientOptions := options.Client().ApplyURI(conURL)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
