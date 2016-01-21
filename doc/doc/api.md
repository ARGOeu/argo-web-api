---
title: API documentation | ARGO
layout: api
page_title: API documentation 
font_title: 'fa fa-cogs'
description: This document describes the API service, using the HTTP application protocol. This API uses XML as the primary exchange format.
---

## Description

The ARGO Web API provides the Serving Layer of ARGO. It is comprised of a high
performance and scalable data store and a multi-tenant REST HTTP API, which is
used for retrieving the Status, Availability and Reliability reports and the
actual raw metric results.

## Installation

### RPM install

You will need a RHEL 6.x or similar (base installation) to proceed. As a first
step make sure that on your host an ntp client service is configured properly. 

#### Software Repositories

On your host the next step is to install (as root user) the ar-release package
via yum:

    yum install http://rpm.hellasgrid.gr/mash/centos6-arstats/x86_64/ar-release-1.0.0-3.21.el6.noarch.rpm

This package will configure on the host(s) the repository files under `/etc/yum.repos.d`.

Also install the EPEL repository. This can be done by installing the
epel-release package for the appropriate OS. For example:

    yum install http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm 

Then you can install the ARGO Web API using yum:

    yum install argo-web-api

### Manual build

1. Install Golang and bzr library

2. Install godep tool

        go get github.com/tools/godep

3. Create a new work space:

        mkdir ~/go-workspace
        export GOPATH=~/go-workspace

  You may add the last `export` line into the `~/.bashrc` or the
  `~/.bash_profile` file to have `GOPATH` environment variable properly setup
  upon  every login.

4. Get the latest version and all dependencies (Using Godep):

        godep update ...

5. To build the service use the following command:

        go build

6. To run the service use the following command:

        ./argo-web-api

7. To run the unit-tests with coverage results:

        gocov test ./... | gocov-xml > coverage.xml

8. To generate and serve godoc (@port 6060)

        godoc -http=:6060
        
## Configuration

The ARGO Web API uses TLS connections and it requires the existence of valid
X.509v3 certificate and the corresponding private key.

For a list of options use the following command:

    ./argo-web-api -h

## Examples

With the following configuration, the API binds to TCP port 443 on all the
available IPs and uses the certificate in `/etc/pki/tls/certs/cert.crt` and its
corresponding private key in `/etc/pki/tls/certs/priv.key`.

    [server]
    bindip = ""
    port = 443
    maxprocs = 4
    cache = false
    lrucache = 700000000
    gzip = true
    cert = /etc/pki/tls/certs/localhost.crt
    privkey = /etc/pki/tls/private/localhost.key
    reqsizelimit = 1073741824

    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "argo_core"

And the API can be started by issuing the command:

    service argo-web-api start

## Links and further reading

- [Swagger : API Demo](http://arpi.afroditi.hellasgrid.gr:8080/)

