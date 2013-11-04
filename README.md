# EGI Availability & Reliability API

## Development

1. Install Golang
2. Create a new work space:

        mkdir ~/go-workspace
        
3. Install dependencies:

        go get labix.org/v2/mgo
        go get labix.org/v2/mgo/bson
        go get code.google.com/p/gcfg
	go get github.com/makistsan/go-lru-cache
	go get github.com/makistsan/go-api

	or

	cd src/api/main
	go get
        
4. Build the service:
	
	cd src/api/main
        go build

5. Test the service:
	cd src/
        go test api/*
        
6. Run the service:

        ~/go-workspace/src/api/main/egi-ari-rest-api
