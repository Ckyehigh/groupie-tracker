package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// API represents the structure of the main API response
type API struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

// Artist represents the structure of an artist
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Image        string   `json:"image"`
}

// Location represents the structure of a location data
type Location struct {
	Locations []string `json:"locations"` // Location list field
}

// Date represents the structure of a date data
type Date struct {
	Dates []string `json:"dates"`
	// Add other fields if needed
}

// Relation represents the structure of a relation data
type Relation struct {
	ArtistName string `json:"artistName"`
	// Add other fields if needed
}

// FetchAPI fetches the API information
func FetchAPI() (*API, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var api API
	err = json.Unmarshal(body, &api)
	if err != nil {
		return nil, err
	}

	return &api, nil
}

// FetchArtists fetches the list of artists
func FetchArtists(url string) ([]Artist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func FetchArtistByID(url string, id int) (*Artist, error) {
	resp, err := http.Get(url + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artist Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return nil, err
	}

	return &artist, nil
}

// FetchLocations fetches the locations for a given artist ID
func FetchLocations(url string, id int) ([]string, error) {
	resp, err := http.Get(url + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var locationData Location
	err = json.Unmarshal(body, &locationData)
	if err != nil {
		return nil, err
	}

	return locationData.Locations, nil
}

// FetchDates fetches all available dates
func FetchDates(url string, id int) ([]string, error) {
	resp, err := http.Get(url + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var datesData Date
	err = json.Unmarshal(body, &datesData)
	if err != nil {
		return nil, err
	}

	return datesData.Dates, nil
}

// FetchRelations fetches the list of relations for a specific artist
func FetchRelations(url string, id int) ([]Relation, error) {
	resp, err := http.Get(url + "/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var relations []Relation
	err = json.Unmarshal(body, &relations)
	if err != nil {
		return nil, err
	}

	return relations, nil
}
