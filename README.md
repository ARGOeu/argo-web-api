# EGI Availability & Reliability API

## Development

1. Install Golang and bzr library
2. Create a new work space:

        mkdir ~/go-workspace
        export GOPATH=~/go-workspace

  You may add the last `export` line into the `~/.bashrc` or the `~/.bash_profile` file to have `GOPATH` environment variable properly setup upon every login. 

3. Get the latest version and all dependencies:

        go get github.com/ARGOeu/argo-web-api

4. To build the service use the following command:

        go build

5. To run the service use the following command:

        ./argo-web-api

  For a list of options use the following command:

        ./argo-web-api -h

6. To run the unit-tests:

        go test ./...

7. To generate and serve godoc (@port 6060)

        godoc -http=:6060
