package status

import (
	"go.mongodb.org/mongo-driver/bson"
)

const statusGroupColName = "status_endpoint_groups"
const statusEndpointColName = "status_endpoints"

func queryGroups(input InputParams, reportID string) bson.M {
	filter := bson.M{
		"date_integer": bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":       reportID,
	}

	if len(input.group) > 0 {
		filter["endpoint_group"] = input.group
	}

	return filter
}

func queryEndpoints(input InputParams, reportID string) bson.M {

	// prepare the match filter
	filter := bson.M{
		"date_integer": bson.M{"$gte": input.startTime, "$lte": input.endTime},
		"report":       reportID,
	}

	if input.ID != "" {
		filter["info.ID"] = input.ID
	}

	return filter
}
