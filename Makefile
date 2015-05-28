PKGNAME=argo-web-api
SPECFILE=${PKGNAME}.spec
PKGVERSION = $(shell grep -s '^Version:' $(SPECFILE) | sed -e 's/Version: *//')

sources:
	mkdir -p ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/argo-web-api
	cp -rp argo-web-api.conf app utils *.go ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/argo-web-api
	cp ${SPECFILE} ${PKGNAME}-${PKGVERSION}/src/github.com/argoeu/argo-web-api
	tar czf ${PKGNAME}-${PKGVERSION}.tar.gz ${PKGNAME}-${PKGVERSION}
	rm -fr ${PKGNAME}-${PKGVERSION}
