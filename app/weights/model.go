package weights

//Weights holds a list of group weights of specific type and of specific group type
type Weights struct {
	ID         string   `bson:"id" json:"id"`
	Name       string   `bson:"name" json:"name"`
	WeightType string   `bson:"weight_type" json:"weight_type"`
	GroupType  string   `bson:"group_type" json:"group_type"`
	Groups     []Weight `bson:"groups" json:"groups"`
}

//Weight hols a mapping between group name and weight value
type Weight struct {
	Name  string  `bson:"name" json:"name"`
	Value float64 `bson:"value" json:"value"`
}

// SelfReference to hold links and id
type SelfReference struct {
	ID    string `json:"id" bson:"id,omitempty"`
	Links Links  `json:"links"`
}

// Links struct to hold links
type Links struct {
	Self string `json:"self"`
}
