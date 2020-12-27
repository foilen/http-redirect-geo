package main

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func dbIPResolve(db *geoip2.Reader, textIP string) (float64, float64, error) {

	ip := net.ParseIP(textIP)
	record, err := db.City(ip)
	if err != nil {
		return 0, 0, err
	}
	if record.Location.Latitude == 0 && record.Location.Longitude == 0 {
		return 0, 0, errors.New("IP location not found")
	}

	return record.Location.Latitude, record.Location.Longitude, nil

}
