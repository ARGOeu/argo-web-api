# EGI Availability & Reliability API

## Development

1. Install Golang
2. Create a new work space:

        mkdir ~/go-workspace
        
3. Install dependencies:

        go get labix.org/v2/mgo
        go get labix.org/v2/mgo/bson
        go get code.google.com/p/gcfg
        go get github.com/argoeu/go-lru-cache
	or

        go get
        
4. Build the service:
	
        go build

5. Run the service:

        ./ar-web-api
