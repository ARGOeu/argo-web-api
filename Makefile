PKGNAME=ar-web-api
SPECFILE=${PKGNAME}.spec
PKGVERSION = $(shell grep -s '^Version:' $(SPECFILE) | sed -e 's/Version: *//')

sources:
	mkdir ${PKGNAME}-${PKGVERSION}
	cp -rp go-api.conf go-api.init src ${PKGNAME}-${PKGVERSION}
	cp ${SPECFILE} ${PKGNAME}-${PKGVERSION}
	tar czf ${PKGNAME}-${PKGVERSION}.tar.gz ${PKGNAME}-${PKGVERSION}
	rm -fr ${PKGNAME}-${PKGVERSION}
