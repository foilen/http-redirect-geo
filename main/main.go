package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/umahmood/haversine"
)

type urlAndLocation struct {
	URL   string
	Coord haversine.Coord
}

var rootConfig *RootConfiguration
var urlAndLocations []urlAndLocation
var nextRequestID *uint64 = new(uint64)
var geoIPDb *geoip2.Reader

func main() {

	// Get the configuration
	args := os.Args[1:]

	if len(args) != 1 {
		panic("You need to specify the config file to use")
	}

	var err error
	rootConfig, err = getRootConfiguration(args[0])
	if err != nil {
		panic(err)
	}

	// Open DBIP
	geoIPDb, err = geoip2.Open(rootConfig.DbIPFile)
	if err != nil {
		log.Fatal(err)
	}
	defer geoIPDb.Close()

	// Get the location of each redirections
	urlAndLocations = make([]urlAndLocation, len(rootConfig.RedirectionUrls))
	for i, v := range rootConfig.RedirectionUrls {
		log.Println("Getting the location for", v)
		hostname := strings.Split(v, "://")[1]
		addr, err := net.LookupHost(hostname)
		if err != nil {
			log.Fatal(err)
		}
		textIP := addr[0]

		latitude, longitude, err := dbIPResolve(geoIPDb, textIP)
		if err != nil {
			log.Fatal(err)
		}

		urlAndLocations[i].URL = v
		urlAndLocations[i].Coord = haversine.Coord{Lat: latitude, Lon: longitude}
		log.Println("Location of", v, "with ip", textIP, "is", latitude, longitude)
	}

	// Start the web service
	log.Println("Starting web service on port", rootConfig.Port)

	http.HandleFunc("/", handler)
	log.Print(http.ListenAndServe(fmt.Sprintf(":%v", rootConfig.Port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	requestID := "[" + strconv.FormatUint(atomic.AddUint64(nextRequestID, 1), 10) + "]"

	beginTime := time.Now()

	// Get the details
	path := r.URL.Path
	var textIP = r.Header.Get("X-Forwarded-For")
	if textIP == "" {
		textIP = strings.Split(r.RemoteAddr, ":")[0]
	}

	log.Println(requestID, "IP is", textIP, "; path is", path)

	// Get the location
	latitude, longitude, err := dbIPResolve(geoIPDb, textIP)
	var redirectTo string = ""
	var redirectToDistance float64 = -1
	if err == nil {

		log.Println(requestID, "IP's location is", latitude, ";", longitude)

		coordIP := haversine.Coord{Lat: latitude, Lon: longitude}
		redirectTo = urlAndLocations[0].URL + path

		// Get the distance to each redirection and keep the nearest
		for _, urlAndLocation := range urlAndLocations {

			_, distanceKm := haversine.Distance(urlAndLocation.Coord, coordIP)

			if redirectToDistance == -1 || redirectToDistance > distanceKm {
				redirectTo = urlAndLocation.URL
				redirectToDistance = distanceKm
			}

		}

	} else {

		log.Println(requestID, "[WARN] Got an error when getting location:", err)
		redirectTo = urlAndLocations[0].URL

	}

	// Log result
	endTime := time.Now()
	deltaTime := endTime.Sub(beginTime).Nanoseconds()
	redirectTo += path
	if redirectToDistance >= 0 {
		log.Println(requestID, "[OK] Redirecting to ", redirectTo, "which is", redirectToDistance, "km away in", deltaTime, "ns")
	} else {
		log.Println(requestID, "[OK] Redirecting to ", redirectTo, "which per default in", deltaTime, "ns")
	}
	http.Redirect(w, r, redirectTo, http.StatusFound)

}
