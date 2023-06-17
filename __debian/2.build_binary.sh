#!/usr/bin/env bash

PKGDIR=certificatemanager-`cat current_pkg_release`_amd64

mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
mv control ${PKGDIR}/DEBIAN/
mv preinst ${PKGDIR}/DEBIAN/

echo "Building binary from source"
cd ../src
go build -o ../__debian/${PKGDIR}/opt/bin/cm .
strip ../__debian/${PKGDIR}/opt/bin/cm
chown 0:0 ../__debian/${PKGDIR}/opt/bin/cm

echo "Binary built. Now packaging..."
cd ../__debian/
dpkg-deb -b ${PKGDIR}
