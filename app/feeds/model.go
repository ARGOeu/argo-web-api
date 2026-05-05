package feeds

// Topo holds a list of topology feed parameters
type Topo struct {
	TopoType                       string `bson:"type" json:"type"`
	FeedURL                        string `bson:"feed_url" json:"feed_url,omitempty"`
	FeedServiceGroups              string `bson:"feed_service_groups" json:"feed_service_groups,omitempty"`
	FeedServiceEndpoints           string `bson:"feed_service_endpoints" json:"feed_service_endpoints,omitempty"`
	FeedServiceEndpointsExtensions string `bson:"feed_service_endpoints_extensions" json:"feed_service_endpoints_extensions,omitempty"`

	Paginated    string    `bson:"paginated" json:"paginated,omitempty"`
	FetchType    *[]string `bson:"fetch_type" json:"fetch_type,omitempty"`
	UIDendpoints string    `bson:"uid_endpoints" json:"uid_endpoints,omitempty"`
}

// Weights holds a list of weight feed parameters
type Weights struct {
	// name-type of service that provides weights
	FeedType string `bson:"type" json:"type"`
	// url of the feed
	FeedURL string `bson:"feed_url" json:"feed_url"`
	// weight factor hepspec cpu, mem etch
	WeightType string `bson:"weight_type" json:"weight_type"`
	// group type that the weight affects
	GroupType string `bson:"group_type" json:"group_type"`
}

// Data feeds for combined reports
type Data struct {
	Tenants []string `bson:"tenants" json:"tenants"`
}

// Tenant Info containts quick info of tenant name and id
type TenantInfo struct {
	Name string `bson:"name"`
	ID   string `bson:"id"`
}
