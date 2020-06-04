<img src="https://jenkins.argo.grnet.gr/static/3c75a153/images/headshot.png" alt="Jenkins" width="25"/> [![Build Status](https://jenkins.argo.grnet.gr/job/argo-web-api_devel/badge/icon)](https://jenkins.argo.grnet.gr/job/argo-web-api_devel) ![Test Coverage](http://jenkins.argo.grnet.gr:9913/jenkins/c/http/jenkins.argo.grnet.gr/job/argo-web-api_devel)

<img src="https://s3.amazonaws.com/openshift-hub/production/quickstarts/26/travisci.png?1425058399" alt="Travis" width="25"/> [![Build Status](https://travis-ci.org/ARGOeu/argo-web-api.svg?branch=devel)](https://travis-ci.org/ARGOeu/argo-web-api)

<img src="http://cdn.slidesharecdn.com/profile-photo-Swagger-API-32x32.jpg?cb=1439244971" alt="swagger ui" width="25"/> [Swagger UI](https://api-doc.argo.grnet.gr/argo-web-api/)

# ARGO Web-API

ARGO is a flexible and scalable framework for monitoring status, availability and reliability of services provided by infrastructures with medium to high complexity. It can generate multiple reports using customer defined profiles (e.g. for SLA management, operations etc) and has built-in multi-tenant support in the core framework.

## Description 

The ARGO Web API provides the Serving Layer of ARGO. It is comprised of a high performance and scalable datastore and a multi-tenant REST HTTP API, which is used for retrieving the Status, Availability and Reliability reports and the actual raw metric results.


## Installation 

1. Install Golang and bzr library

2. Install godep tool

        go get github.com/tools/godep

3. Create a new work space and setup your environment:

        mkdir ~/go-workspace
        export GOPATH=~/go-workspace
        export PATH=$PATH:GOPATH

  You may add the `export` lines into the `~/.bashrc` or the `~/.bash_profile` file to have the `GOPATH` and `PATH` environment variables properly setup upon every login.

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

## Credits

The ARGO Messaging Service is developed by [GRNET](http://www.grnet.gr)

The work represented by this software was partially funded by:
- EGI Foundation
- EGI-ENGAGE project through the European Union (EU) Horizon 2020 program under Grant number 654142.
- EOSC-Hub project through the European Union (EU) Horizon 2020 program under Grant number 77753642.
- EUDAT2020 European Union’s H2020 Program under Contract No. 654065.

