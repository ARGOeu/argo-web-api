package store

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// OpenSession to mongodb
func GetMongoClient(mongoCfg config.MongoConfig) *mongo.Client {
	conURL := "mongodb://" + mongoCfg.Host + ":" + strconv.Itoa(mongoCfg.Port)
	clientOptions := options.Client().ApplyURI(conURL)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
	opts := options.FindOne().
		SetSort(bson.D{{Key: "date_integer", Value: -1}}).
		SetProjection(bson.M{"date_integer": 1})

	result := DatedItem{}
	err := c.FindOne(context.TODO(), dateQuery, opts).Decode(&result)
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

// check if a replica set is enabled
func IsReplicaSet(client *mongo.Client) bool {
	result := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "isMaster", Value: 1}})
	resRaw, err := result.Raw()
	if err != nil {
		return false
	}
	_, err = resRaw.LookupErr("setName")
	return err == nil
}
