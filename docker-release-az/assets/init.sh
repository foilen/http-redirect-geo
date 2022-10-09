#!/bin/bash
set -e

echo "SSH Start"
service ssh start

if [ -z "$CONFIG_FILE" ]; then
  CONFIG_FILE=/home/site/wwwroot/config.json
fi
echo "CONFIG_FILE : $CONFIG_FILE"

/usr/bin/http-redirect-geo $CONFIG_FILE
