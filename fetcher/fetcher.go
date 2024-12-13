package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// ArtistInfo represents the complete information for an artist
type ArtistInfo struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string `json:"locations"`
	ConcertDates []string `json:"concertDates"`
	Relations    []string `json:"relations"`
}

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
	Locations    []string `json:"-"`
}

var artists []Artist

// Location represents the structure of a location data
type Location struct {
	Locations []string `json:"locations"` // Location list field
}

// Date represents the structure of a date data
type Date struct {
	Dates []string `json:"dates"`
}

// Relation represents the structure of a relation data
type Relation struct {
	ArtistName string `json:"artistName"`
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

// Search function to search for artists, locations, or members using Go routines and channels
func Search(data string) []map[string]string {
	var results []map[string]string
	_, err := strconv.Atoi(string(data[0]))
	if err != nil {
		// Search by artist name, location, or member
		var wg sync.WaitGroup
		resultsChan := make(chan map[string]string, len(artists))
		//idChan := make(chan map[int]int)

		for _, artist := range artists {
			wg.Add(1)
			go func(artist Artist) {
				defer wg.Done()
				if len(data) == 1 {
					if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(data)) {
						resultsChan <- map[string]string{"type": "Band", "id": strconv.Itoa(artist.ID), "name": artist.Name}
						//idChan[artist.ID]++
					}
					locations, _ := FetchLocations("https://groupietrackers.herokuapp.com/api/locations", artist.ID)
					for _, location := range locations {
						if strings.HasPrefix(strings.ToLower(location), strings.ToLower(data)) {
							resultsChan <- map[string]string{"type": "Show Location", "id": strconv.Itoa(artist.ID), "name": location + " - " + artist.Name}
							//idChan[artist.ID]++
						}
					}
					for _, member := range artist.Members {
						if strings.HasPrefix(strings.ToLower(member), strings.ToLower(data)) {
							resultsChan <- map[string]string{"type": "Member", "id": strconv.Itoa(artist.ID), "name": member + " - " + artist.Name}
							//idChan[artist.ID]++
						}
					}
				} else {

					if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(data)) {
						resultsChan <- map[string]string{"type": "Band", "id": strconv.Itoa(artist.ID), "name": artist.Name}
						//idChan[artist.ID]++
					}
					locations, _ := FetchLocations("https://groupietrackers.herokuapp.com/api/locations", artist.ID)
					for _, location := range locations {
						if strings.Contains(strings.ToLower(location), strings.ToLower(data)) {
							resultsChan <- map[string]string{"type": "Show Location", "id": strconv.Itoa(artist.ID), "name": location + " - " + artist.Name}
							//idChan[artist.ID]++
						}
					}
					for _, member := range artist.Members {
						if strings.Contains(strings.ToLower(member), strings.ToLower(data)) {
							resultsChan <- map[string]string{"type": "Member", "id": strconv.Itoa(artist.ID), "name": member + " - " + artist.Name}
							//idChan[artist.ID]++
						}
					}
				}
			}(artist)
		}

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		for result := range resultsChan {
			results = append(results, result)
		}
	} else {
		// Search by first album or creation date
		for _, artist := range artists {
			if strings.Contains(artist.FirstAlbum, data) {
				results = append(results, map[string]string{"type": "First Album", "id": strconv.Itoa(artist.ID), "name": artist.FirstAlbum + " - " + artist.Name})
			}
			if strings.Contains(strconv.Itoa(artist.CreationDate), data) {
				results = append(results, map[string]string{"type": "Creation Date", "id": strconv.Itoa(artist.ID), "name": strconv.Itoa(artist.CreationDate) + " - " + artist.Name})
			}
		}
	}
	return results
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
	for i, day := range datesData.Dates {
		datesData.Dates[i] = strings.TrimLeft(day, "*")
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

//when search string is only one letter only look for results starting with the letter **DONE

//if multiple results when "enter" is pressed reload the page with the search results only

//scroll bar on search suggestions with fixed size and or arrow key selection functionality **DONE
