package openstack

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	host    = "http://169.254.169.254"
	timeout = 1000
)

func GetStringContents(url string) (string, error) {
	client := http.Client{
		Timeout: timeout * time.Millisecond,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func GetAvailabilityZone() (string, error) {
	availability_url := "/latest/meta-data/placement/availability-zone"
	zone, err := GetStringContents(host + availability_url)
	return zone, err
}

func IsKIC() bool {
	zone, err := GetAvailabilityZone()
	kicPrefix := "kr-central-1"
	if err == nil && strings.HasPrefix(zone, kicPrefix) {
		return true
	} else {
		return false
	}
}
