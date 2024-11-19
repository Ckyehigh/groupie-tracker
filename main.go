package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"groupie-tracker/fetcher"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api, err := fetcher.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		artists, err := fetcher.FetchArtists(api.Artists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		err = tmpl.Execute(w, artists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/artists/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artists/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}

		api, err := fetcher.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		artists, err := fetcher.FetchArtists(api.Artists)
		if err != nil {
			log.Println("Error fetching artists:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if id < 1 || id > len(artists) {
			tmpl := template.Must(template.ParseFiles("templates/404.html"))
			w.WriteHeader(http.StatusNotFound)
			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Println("Error executing 404 template:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		artist, err := fetcher.FetchArtistByID(api.Artists, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/artists.html"))
		err = tmpl.Execute(w, artist)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/artists/locations/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artists/locations/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}

		api, err := fetcher.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		locations, err := fetcher.FetchLocations(api.Locations, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(locations)
	})

	http.HandleFunc("/artists/dates/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artists/dates/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}

		api, err := fetcher.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dates, err := fetcher.FetchDates(api.Dates, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Fetched dates:", dates)
		err = json.NewEncoder(w).Encode(dates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/artists/relations/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/artists/relations/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}

		api, err := fetcher.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		relations, err := fetcher.FetchRelations(api.Relation, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(relations)
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
