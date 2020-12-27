# About

Web server that redirect to different urls depending on the IP's geolocalisation.

You can use the free version of DBIP in MMDB format: https://db-ip.com/db/download/ip-to-city-lite .

## Logic

The logic is simple:
- Get the IPs of the redirection urls and get the location in lat/lon
- When a request comes in, get its IP address from:
    - X-Forwarded-For (to get the right one when behind a load-balancer)
    - or the remote address
- With the IP address, get the location in lat/lon and check which redirection is closest using the Haversine formula to determine the great-circle distance (https://en.wikipedia.org/wiki/Haversine_formula)

# Local Usage

## Compile

`./create-local-release.sh`

The file is then in `build/bin/http-redirect-geo`

## Config file content

```
cat > _config.json << _EOF
{
    "port" : 8888,
    "dbIpFile" : "dbip-city-lite.mmdb",
    "redirectionUrls" : [
        "https://tor.cdn.foilen.com",
        "https://fra.cdn.foilen.com"
    ]
}
_EOF
```

## Execute

To execute:
`./build/bin/http-redirect-geo _config.json`

# Create release

`./create-public-release.sh`

That will show the latest created version. Then, you can choose one and execute:
`./create-public-release.sh X.X.X`

# Use with debian

```bash
echo "deb https://dl.bintray.com/foilen/debian stable main" | sudo tee /etc/apt/sources.list.d/foilen.list
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 379CE192D401AB61
sudo apt update
sudo apt install http-redirect-geo
```

# Use with Docker

```bash
docker run --rm -ti foilen/foilen-http-redirect-geo:latest /usr/bin/http-redirect-geo
```
