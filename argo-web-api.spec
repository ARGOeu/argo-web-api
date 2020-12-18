#debuginfo not supported with Go
%global debug_package %{nil}

Name: argo-web-api
Summary: A/R API
Version: 1.9.2
Release: 1%{?dist}
License: ASL 2.0
Buildroot: %{_tmppath}/%{name}-buildroot
Group:     ARGO
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang
BuildRequires: bzr
BuildRequires: git
Requires: mongo-10gen
Requires: mongo-10gen-server
ExcludeArch: i386
Obsoletes: ar-web-api

%description
Installs the ARGO API.

%pre
/usr/bin/getent group argo-web-api || /usr/sbin/groupadd -r argo-web-api
/usr/bin/getent passwd argo-web-api || /usr/sbin/useradd -r -s /sbin/nologin -d /var/www/argo-web-api -g argo-web-api argo-web-api

%prep
%setup

%build
export GOPATH=$PWD
export PATH=$GOPATH/bin:$PATH
cd src/github.com/ARGOeu/argo-web-api/
export GIT_COMMIT=$(git rev-list -1 HEAD)
export BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
go install -ldflags "-X github.com/ARGOeu/argo-web-api/version.Commit=$GIT_COMMIT -X github.com/ARGOeu/argo-web-api/version.BuildTime=$BUILD_TIME"

%install
%{__rm} -rf %{buildroot}
install --directory %{buildroot}/var/www/argo-web-api
install --mode 755 bin/argo-web-api %{buildroot}/var/www/argo-web-api/argo-web-api

install --directory %{buildroot}/etc
install --mode 644 src/github.com/ARGOeu/argo-web-api/default.conf %{buildroot}/etc/argo-web-api.conf

install --directory %{buildroot}/etc/init
install --mode 644 src/github.com/ARGOeu/argo-web-api/argo-web-api.conf %{buildroot}/etc/init/

install --directory %{buildroot}/var/www/argo-web-api/certs

install --directory %{buildroot}/usr/lib/systemd/system
install --mode 644 src/github.com/ARGOeu/argo-web-api/argo-web-api.service %{buildroot}/usr/lib/systemd/system/

%clean
%{__rm} -rf %{buildroot}
export GOPATH=$PWD
cd src/github.com/ARGOeu/argo-web-api/
go clean

%files
%defattr(0644,argo-web-api,argo-web-api)
%attr(0755,argo-web-api,argo-web-api) /var/www/argo-web-api
%attr(0755,argo-web-api,argo-web-api) /var/www/argo-web-api/certs
%attr(0755,argo-web-api,argo-web-api) /var/www/argo-web-api/argo-web-api
%caps(cap_net_bind_service=+ep) /var/www/argo-web-api/argo-web-api
%config(noreplace) %attr(0644,argo-web-api,argo-web-api) /etc/argo-web-api.conf
%config(noreplace) %attr(0644,argo-web-api,argo-web-api) /etc/init/argo-web-api.conf
%attr(0644,root,root) /usr/lib/systemd/system/argo-web-api.service

%changelog
* Wed Dec 16 2020 Konstantinos Kagkelidis <kaggis@gmail.com> 1.9.2-1%{dist}
- Release of argo-web-api version 1.9.2
* Wed Jul 08 2020 Konstantinos Kagkelidis <kaggis@gmail.com> 1.9.1-1%{dist}
- Release of argo-web-api version 1.9.1
* Thu Mar 26 2020 Konstantinos Kagkelidis <kaggis@gmail.com> 1.9.0-1%{dist}
- Release of argo-web-api version 1.9.0
* Thu Nov 14 2019 Konstantinos Kagkelidis <kaggis@gmail.com> 1.8.1-1%{dist}
- Release of argo-web-api version 1.8.1
* Mon Nov 4 2019 Konstantinos Kagkelidis <kaggis@gmail.com> 1.8.0-1%{dist}
- Release of argo-web-api version 1.8.0
- Refactor building method to include commit hash and build time information into the binary
* Fri Mar 22 2019 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.9-1%{dist}
- ARGO-1455 - Migrate to golang/dep tool
- ARGO-1438 Implement tenant general status
- ARGO-1680 serve endpoint a/r results
* Wed Nov 7 2018 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.8-2%{dist}
- ARGO-1435 Fix configuration_profile json field in tenant status call
- ARGO-1433 Add tenant status roles to init db script
- ARGO-1268 Serve topology statistics per report
- ARGO-451 Close status timelines with latest daily result
* Tue Sep 18 2018 Angelos Tsalapatis <agelos.tsal@gmail.com> 1.7.7-1%{dist}
- ARGO-1390 API CALL - Update recomputation
- ARGO-1389 API CALL - Delete Recomputation
- ARGO-1395 Operations profile name field should be unique
- ARGO-1396 Metric profile name field should be unique
- ARGO-1394 Aggregation profile name field should be unique
* Wed Sep 12 2018 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.6-1%{dist}
- ARGO-1298 Show/Update tenant's argo-engine status
* Wed Sep 12 2018 Angelos Tsalapatis <agelos.tsal@gmail.com> 1.7.5-1%{dist}
- ARGO-1381 Api call update report name field not unique
- ARGO-1388 Api call update tenant name field
- ARGO-1345 update Tenant model to handle field roles
- ARGO-1391 Wrong response for empty factors list
- ARGO-1381 Refactor error messages in argo-web-api thresholds package
* Tue Sep 4 2018 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.4-1%{dist}
- ARGO-545 Add api call for latest non-ok entries
* Tue Sep 4 2018 Angelos Tsalapatis <agelos.tsal@gmail.com> 1.7.3-1%{dist}
- ARGO-1380 Refactor error messages in argo-web-api tenants package
- ARGO-1337 Refactor error messages in argo-web-api factors package
- ARGO-445 Recomputation details error
- ARGO-1379 Refactor error messages in the reports package
* Tue Aug 21 2018 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.2-1%{dist}
- ARGO-1351 Refactor error messages in the aggregation profiles package
- ARGO-1349 Refactor error messages in the metric profiles package
- ARGO-1346 Refactor error messages in the opperations package
- ARGO-1275 Refactor Report resource schema
- ARGO-1260 Implement CRUD on threshold profiles resource
- ARGO-1099 Add read-only super-admin
- ARGO-894 Fix error handling for internal server errors
- ARGO-835 Set timeout to hbase related requests
- ARGO-776 Show latest status results if no data are present for the be…
- ARGO-794 Hbase zkquorum config fix. Hbase query minor fixes
- ARGO-723 Add hbase support to argo-web-api
* Mon Dec 12 2016 Konstantinos Kagkelidis <kaggis@gmail.com> 1.7.1-1%{dist}
- Set log output to stdout
- ARGO-606 Add WEB API username to logging
* Mon Oct 24 2016 Themis Zamani <themiszamani@gmail.com> - 1.6.5-2%{?dist}
- New RPM package release.
* Wed Oct 12 2016 Themis Zamani <themiszamani@gmail.com> - 1.6.5-1%{?dist}
- New RPM package release.
* Sat Sep 24 2016 Themis Zamani <themiszamani@gmail.com> - 1.6.4-1%{?dist}
- Update to the latest GOLang version (1.7.1).
* Wed Feb 3 2016 Christos Kanellopoulos <skanct@gmail.com> - 1.6.3-1%{?dist}
- ARGO-292 Use godep tool for 3rd party pkg management
- ARGO-291 use mongodb 3.x in travis
- ARGO-284 Implement factors in APIv2
- Improved documentation
* Thu Dec 3 2015 Avraam Tsantekidis <avraamt@lab.grid.auth.gr> - 1.6.2-1%{?dist}
- ARGO-259 results swagger definitions
- ARGO-257 Swagger documentation for status requests
- ARGO-289 Fix swagger tag misplacement
- ARGO-279 Add test coverage metrics
- Fix routes to not require trailing slash
- ARGO-260 Swagger definition for tenants
- Add check for existence of time range
- ARGO-261 Swagger definitions for report calls
- ARGO-264 Swagger yaml definitions for operations profiles
- ARGO-263 Swagger definition for aggregation profiles
- ARGO-262 Swagger yaml definitions for metric profiles
* Thu Nov 12 2015 Avraam Tsantekidis <avraamt@lab.grid.auth.gr> - 1.6.1-1%{?dist}
- ARGO-256 fixes for status reponses
- ARGO-245 Reference and tag results using report uuid
* Wed Oct 14 2015 Paschalis Korosoglou <pkoro@grid.auth.gr> - 1.6.0-3%{?dist}
- Adds service configuration file
* Mon Aug 10 2015 Paschalis Korosoglou <pkoro@grid.auth.gr> -  1.6.0-2%{?dist}
- Correction in cases imports
* Thu May 28 2015 Pavlos Daoglou <pdaog@grid.auth.gr> - 1.6.0-1%{?dist}
- ARGO-104 Update github import urls to be consistent with the repo name changes
* Wed May 6 2015 Konstantinos Kagkelidis <kaggis@gmail.com> - 1.5.1-5%{?dist}
- Fix Av.profile update/delete responses. Add Check for valid object ids
* Thu Jan 15 2015 Konstantinos Kagkelidis <kaggis@gmail.com> - 1.5.1-2%{?dist}
- Add prev timestamp support at the beginning of status timelines
* Fri Dec 19 2014 Konstantinos Kagkelidis <kaggis@gmail.com> - 1.5.1-1%{?dist}
- Add support for endpoint/service/sites aggregations
* Wed Dec 17 2014 Konstantinos Kagkelidis <kaggis@gmail.com> - 1.5.0-1%{?dist}
- Add support for Status Results and raw data results
* Mon Jun 16 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.4.0-1%{?dist}
- Added support for custom factor retrival per site
* Tue Jun 3 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.3.0-1%{?dist}
- Major code refactoring, proper error handling, http headers
* Mon May 12 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.5-1%{?dist}
- POEM profile retrieval support
* Wed May 7 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.4-2%{?dist}
- VO and SF resutls fix
* Wed Apr 30 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.4-1%{?dist}
- Various changes and bug fixes
* Thu Apr 24 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.3-1%{?dist}
- Added support for service flavor result querying
* Wed Apr 16 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.2-1%{?dist}
- Fixed sites result querying
* Wed Apr 09 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.2-1%{?dist}
- Added CRUD support for Availability Profiles
* Tue Mar 25 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.1-2%{?dist}
- Fixed recalculation history bug
* Thu Mar 20 2014 Nikos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.1-1%{?dist}
- Changes in results querying to reflect new database schema
* Wed Mar 19 2014 Nikolaos Triantafyllidis <ntrianta@grid.auth.gr> - 1.2.0-1%{?dist}
- Support for VOs. Changes in grouping
* Tue Mar 4 2014 Paschalis Korosoglou <pkoro@grid.auth.gr> - 1.1.1-1%{?dist}
- Suport for https
* Thu Feb 6 2014 Paschalis Korosoglou <pkoro@grid.auth.gr> - 1.1.0-2%{?dist}
- Fix in spec file
* Thu Feb 6 2014 Paschalis Korosoglou <pkoro@grid.auth.gr> - 1.1.0-1%{?dist}
- Fix in Av computation
* Thu Nov 7 2013 Paschalis Korosoglou <pkoro@grid.auth.gr> - 1.0.17-2%{?dist}
- Initial koji import
