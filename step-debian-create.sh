#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Create .deb ]----
DEB_FILE=http-redirect-geo_${VERSION}_amd64.deb
DEB_PATH=$RUN_PATH/build/debian_out/http-redirect-geo
rm -rf $DEB_PATH
mkdir -p $DEB_PATH $DEB_PATH/DEBIAN/ $DEB_PATH/usr/bin/

cat > $DEB_PATH/DEBIAN/control << _EOF
Package: http-redirect-geo
Version: $VERSION
Maintainer: Foilen
Architecture: amd64
Description: Web server that redirect to different urls depending on the IP's geolocalisation
_EOF

cat > $DEB_PATH/DEBIAN/postinst << _EOF
#!/bin/bash

set -e
_EOF
chmod +x $DEB_PATH/DEBIAN/postinst

cp -rv build/bin/* $DEB_PATH/usr/bin/

cd $DEB_PATH/..
dpkg-deb --no-uniform-compression --build http-redirect-geo
mv http-redirect-geo.deb $DEB_FILE
