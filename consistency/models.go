package consistency

type Result struct {
	Status             string `json:"status"`
	Timestamp          string `json:"timestamp"`
	Message            string `json:"message"`
	AutoCheckStatus    string `json:"auto_check_status,omitempty"`
	AutoCheckMsg       string `json:"auto_check_mesage,omitempty"`
	AutoCheckTimestamp string `json:"auto_check_timestamp,omitempty"`
	AckStatus          string `json:"ack_status,omitempty"`
	AckMsg             string `json:"ack_message,omitempty"`
	AckTimestamp       string `json:"ack_timestamp,omitempty"`
	AckTimeoutHours    int    `json:"ack_timeout_hours,omitempty"`
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

type AckMongo struct {
	AckStatus       string `bson:"ack_status" json:"status"`
	AckMsg          string `bson:"ack_message" json:"message"`
	AckTimestamp    string `bson:"ack_timestamp" json:"timestamp"`
	AckTimeoutHours int    `bson:"ack_timeout_hours" json:"timeout_hours"`
}

type AutoCheckMongo struct {
	AutoCheckStatus    string `bson:"auto_check_status" json:"status"`
	AutoCheckMsg       string `bson:"auto_check_message" json:"message"`
	AutoCheckTimestamp string `bson:"auto_check_timestamp" json:"timestamp"`
}

// errorMessage struct to hold the json error response
type errorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
