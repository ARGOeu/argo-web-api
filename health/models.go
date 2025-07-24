package health

type Result struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type DataMongo struct {
	AutoCheckStatus    string `bson:"auto_check_status"`
	AutoCheckMsg       string `bson:"auto_check_message"`
	AutoCheckTimestamp string `bson:"auto_check_timestamp"`
	AckStatus          string `bson:"ack_status"`
	AckMsg             string `bson:"ack_message"`
	AckTimestamp       string `bson:"ack_timestamp"`
	AckTimeoutHours    int    `bson:"ack_timeout_hours"`
}

// errorMessage struct to hold the json error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
