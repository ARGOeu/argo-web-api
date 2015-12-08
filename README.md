[![Build Status](https://travis-ci.org/ARGOeu/argo-web-api.svg?branch=devel)](https://travis-ci.org/ARGOeu/argo-web-api)
# EGI Availability & Reliability API

## Development

1. Install Golang and bzr library

2. Install godep tool

        go get github.com/tools/godep

3. Create a new work space:

        mkdir ~/go-workspace
        export GOPATH=~/go-workspace

  You may add the last `export` line into the `~/.bashrc` or the `~/.bash_profile` file to have `GOPATH` environment variable properly setup upon every login.

4. Get the latest version and all dependencies (Using Godep):

        godep update ...

5. To build the service use the following command:

        go build

6. To run the service use the following command:

        ./argo-web-api

  For a list of options use the following command:

        ./argo-web-api -h

7. To run the unit-tests with coverage results:

        gocov test ./... | gocov-xml > coverage.xml

8. To generate and serve godoc (@port 6060)

        godoc -http=:6060
