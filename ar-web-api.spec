Name: ar-web-api
Summary: A/R API
Version: 1.5.1
Release: 1%{?dist}
License: ASL 2.0
Buildroot: %{_tmppath}/%{name}-buildroot
Group:     EGI/SA4
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang
BuildRequires: bzr
BuildRequires: git
Requires: mongo-10gen
Requires: mongo-10gen-server
ExcludeArch: i386

%description
Installs the A/R API.

%prep
%setup

%build
export GOPATH=$PWD
cd src/github.com/argoeu/ar-web-api/
go get
go install

%install
%{__rm} -rf %{buildroot}
install --directory %{buildroot}/var/www/ar-web-api
install --mode 755 bin/ar-web-api %{buildroot}/var/www/ar-web-api/ar-web-api

install --directory %{buildroot}/etc/init
install --mode 644 src/github.com/argoeu/ar-web-api/ar-web-api.conf %{buildroot}/etc/init/

%clean
%{__rm} -rf %{buildroot}
export GOPATH=$PWD
cd src/github.com/argoeu/ar-web-api/
go clean

%files
%defattr(0644,root,root)
%attr(0750,root,root) /var/www/ar-web-api
%attr(0755,root,root) /var/www/ar-web-api/ar-web-api
%attr(0644,root,root) /etc/init/ar-web-api.conf

%changelog
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
