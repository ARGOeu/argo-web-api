package feeds

//FeedsTopo holds a list of topology feed parameters
type FeedsTopo struct {
	TopoType     string   `bson:"type" json:"type"`
	FeedURL      string   `bson:"feed_url" json:"feed_url"`
	Paginated    string   `bson:"paginated" json:"paginated"`
	FetchType    []string `bson:"fetch_type" json:"fetch_type"`
	UIDendpoints string   `bson:"uid_endpoints" json:"uid_endpoints"`
}
