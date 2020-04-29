#!/bin/bash

set -x

project=$DEP_PROJECT
version=$DEP_VERSION
release=$DEP_RELEASE
# Verify deploy variables are set
[ -z "$project" ] && echo ">>> \$DEP_PROJECT variable is empty" && exit 1
[ -z "$version" ] && echo ">>> \$DEP_VERSION variable is empty" && exit 1
[ -z "$release" ] && echo ">>> \$DEP_RELEASE variable is empty" && exit 1

CUR_VERSION=`yum info $project | grep "Version" | awk '{print $3}'`
CUR_RELEASE=`yum info $project | grep "Release" | awk '{print $3}'`

[ -z "$CUR_VERSION" ] && echo "$project is not installed" && exit 1
[ -z "$CUR_RELEASE" ] && echo "$project is not installed" && exit 1

echo ">>> Update $project"
yum -y update $project-$version-$release
exit_status=$?
if [ $exit_status -eq 1 ]; then
    echo ">>> Deployment Failed"
    exit $exit_status
fi

echo ">>> Check new version $project"
VERSION_NEW=`yum info $project | grep "Version" | awk '{print $3}'`
RELEASE_NEW=`yum info $project | grep "Release" | awk '{print $3}'`

if [[ "$version" != "$VERSION_NEW"  ||  "$release" != "$RELEASE_NEW" ]];
then
    echo "Release failed"
    exit 1
fi

echo ">>> Restart Service"
systemctl daemon-reload
systemctl restart $project
echo ">>> Check Service"
systemctl is-active --quiet $project
exit_status=$?
if [ $exit_status -eq 1 ]; then
    echo ">>> Deployment Failed"
    exit $exit_status
fi
echo "=== Deployment Finished Successfully ==="
