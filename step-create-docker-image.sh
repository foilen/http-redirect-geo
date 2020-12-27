#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Prepare folder for docker image ]----
DOCKER_BUILD=$RUN_PATH/build/docker

rm -rf $DOCKER_BUILD
mkdir -p $DOCKER_BUILD/usr/bin/

cp -v build/bin/http-redirect-geo $DOCKER_BUILD/usr/bin/

cat > $DOCKER_BUILD/Dockerfile << _EOF
FROM ubuntu:20.04

COPY usr/ /usr/
RUN chmod 755 /usr/bin/http-redirect-geo

CMD /bin/bash
_EOF



echo ----[ Docker image folder content ]----
find $DOCKER_BUILD

echo ----[ Build docker image ]----
DOCKER_IMAGE=foilen-http-redirect-geo:$VERSION
docker build -t $DOCKER_IMAGE $DOCKER_BUILD
docker tag $DOCKER_IMAGE foilen/$DOCKER_IMAGE

rm -rf $DOCKER_BUILD
