

# Prerequisites

- You will need a RHEL 6.x or similar OS (base installation) to proceed. Note that the following instructions have been tested against CentOS 6.x OSes. 
- Make sure that on your host an ntp client service is configured properly. 
- Configure the OS firewall to accept incoming `tcp` connections to port `443`.

The first step is to install (as root user) the `epel` and `argo` release packages via yum:

    # yum install http://dl.fedoraproject.org/pub/epel/6/x86_64/epel-release-6-8.noarch.rpm
    # yum install http://rpm.hellasgrid.gr/mash/centos6-arstats/i386/ar-release-1.0.0-3.el6.noarch.rpm

These packages will configure on the host the necessary repository files under `/etc/yum.repos.d`.

# Installation

Install the ARGO web API service with the following command:

    # yum install argo-web-api

# Configuration

Edit the `/etc/argo-web-api.conf` configuration file and 

- set the values of the `mongo_host` and `mongo_port` variables to point to a running mongo service
- set the value of the `db` parameter in the `mongo` section to the name of the core database holding the tenant's information

Optionally, you may further configure service parameters (i.e. interface to bind to, port number etc). 

# Services

Start the service using the following command:

    # start argo-web-api

To stop the service use the following command:

    # stop argo-web-api

The check is the service is started (running) or not use the following command:

    # status argo-web-api

