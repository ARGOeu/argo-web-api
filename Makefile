PKGNAME=ar-web-api
SPECFILE=${PKGNAME}.spec
PKGVERSION = $(shell grep -s '^Version:' $(SPECFILE) | sed -e 's/Version: *//')

sources:
	mkdir -p ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/ar-web-api
	cp -rp go-api.conf app utils *.go ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/ar-web-api
	cp ${SPECFILE} ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/ar-web-api
	tar czf ${PKGNAME}-${PKGVERSION}.tar.gz ${PKGNAME}-${PKGVERSION}
	rm -fr ${PKGNAME}-${PKGVERSION}
