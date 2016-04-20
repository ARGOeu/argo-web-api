package authorization

import (
	"log"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"

	"gopkg.in/mgo.v2/bson"
)

// QRole holds roles resources relationships
type QRole struct {
	Resource string   `bson:"resource"`
	Roles    []string `bson:"roles"`
}

// HasResourceRoles returns if a resource has the roles given
func HasResourceRoles(cfg config.Config, resource string, roles []string) bool {

	session, err := mongo.OpenSession(cfg.MongoDB)
	defer mongo.CloseSession(session)

	if err != nil {
		panic(err)
	}

	var results []QRole

	query := bson.M{"resource": resource, "roles": bson.M{"$in": roles}}
	err = mongo.Find(session, cfg.MongoDB.Db, "roles", query, "resource", &results)

	if err != nil {
		log.Fatal(err)
	}

	if len(results) > 0 {
		return true
	}

	return false

}
