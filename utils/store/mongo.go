package store

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"go.mongodb.org/mongo-driver/bson"
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

type DatedItem struct {
	DateInteger int `bson:"date_integer"`
}

func GetCloseDate(c *mongo.Collection, dt int) int {
	dateQuery := bson.M{"date_integer": bson.M{"$lte": dt}}
	result := DatedItem{}
	err := c.FindOne(context.TODO(), dateQuery).Decode(&result)
	if err != nil {
		return -1
	}
	return result.DateInteger
}

// GetReportID accepts a report name and returns the report's id
func GetReportID(col *mongo.Collection, report string) (string, error) {
	var result map[string]interface{}
	// reports are stored to the reports collection
	// query based on report name which is included in the info element
	query := bson.M{"info.name": report}
	// Execute the query and grab only the first result
	err := col.FindOne(context.TODO(), query).Decode(&result)

	if result != nil {
		return result["id"].(string), err
	}
	return "", err

}
