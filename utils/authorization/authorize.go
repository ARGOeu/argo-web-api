package authorization

import (
	"context"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"go.mongodb.org/mongo-driver/bson"
)

// QRole holds roles resources relationships
type QRole struct {
	Resource string   `bson:"resource"`
	Roles    []string `bson:"roles"`
}

// HasResourceRoles returns if a resource has the roles given
func HasResourceRoles(cfg config.Config, resource string, roles []string) bool {

	rolesCol := cfg.MongoClient.Database(cfg.MongoDB.Db).Collection("roles")

	query := bson.M{"resource": resource, "roles": bson.M{"$in": roles}}
	queryResult := rolesCol.FindOne(context.TODO(), query)
	return (queryResult.Err() == nil)
}
