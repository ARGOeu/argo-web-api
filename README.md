# EGI Availability & Reliability API

## Development

1. Install Golang
2. Create a new work space:

        mkdir ~/go-workspace
        
3. Install dependencies:

        go get github.com/dpapathanasiou/go-api
        go get labix.org/v2/mgo
        go get labix.org/v2/mgo/bson
        
4. Build the service:

        go build
        
5. Run the service:

        ./egi-ari-rest-api
