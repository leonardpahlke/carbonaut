/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"carbonaut.cloud/pkg/schema"
)

// IPAPI is an geolocation provider that uses the ip-api.com API

// Tools to discover the geolocation of a machine
// this information is need to to discover the energy mix of the grid
// and further to calculate the carbon intensity of the grid
type Provider struct{}

// The Geolocation look up by ip uses the API from https://ip-api.com/docs which throttles the requests to 45 per minute
// requires a valid IPv4 address
func (Provider) CollectGeolocation(ip string) (*schema.Geolocation, error) {
	var r Geolocation
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, fmt.Errorf("error while getting location: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}
	r, err = UnmarshalGeolocation(body)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling location: %w", err)
	}
	if r.Status == "fail" {
		return nil, fmt.Errorf("error while getting location: %s, message: %s", r.Status, r.Message)
	}
	return &schema.Geolocation{
		Country:     r.Country,
		CountryCode: r.CountryCode,
		Region:      r.Region,
		RegionName:  r.RegionName,
		City:        r.City,
		Zip:         r.Zip,
		Lat:         r.Lat,
		Lon:         r.Lon,
		IP:          ip,
	}, nil
}

func UnmarshalGeolocation(data []byte) (Geolocation, error) {
	var r Geolocation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Geolocation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Geolocation struct {
	// IP address
	Query string `json:"query"`
	// Status of the request. Either "success" or "fail"
	Status string `json:"status"`
	// Error message if status is "fail"
	Message string `json:"message"`
	// Country name (e.g. "United States")
	Country string `json:"country"`
	// Country code (e.g. "US")
	CountryCode string `json:"countryCode"`
	// Region code (e.g. "CA")
	Region string `json:"region"`
	// Region name (e.g. "California")
	RegionName string `json:"regionName"`
	// City (e.g. "Mountain View")
	City string `json:"city"`
	// Zip code
	Zip string `json:"zip"`
	// Latitude
	Lat float64 `json:"lat"`
	// Longitude
	Lon float64 `json:"lon"`
	// Timezone (e.g. "America/Los_Angeles")
	Timezone string `json:"timezone"`
	// Internet Service Provider
	ISP string `json:"isp"`
	// Organization name
	Org string `json:"org"`
	// AS number and organization, separated by space (RIR). Empty for IP blocks not being announced in BGP tables.
	As string `json:"as"`
}
