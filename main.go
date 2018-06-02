package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	_ "github.com/mattn/go-sqlite3"

	"encoding/json"
)

type Page struct {
	Name     string
	DBStatus bool
}

type SearchResult struct {
	Title  string
	Author string
	Year   string
	ID     string
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "library.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := Page{Name: "Gopher"}

		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}
		p.DBStatus = db.Ping() == nil

		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
			results := []SearchResult{
				SearchResult{"Moby Dick", "Herman Melville", "1851", "222222"},
				SearchResult{"The Adventures of Huckleberry Finn", "Mark Twain", "1884", "444444"},
				SearchResult{"The Catcher in the Rye", "JD Salinger", "1951", "3333333"},
			}

			encoder := json.NewEncoder(w)
			if err := encoder.Encode(results); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

func search(query string) ([]SearchResult, error) {
	var resp *http.Response
	var err error

	//send query to classify api
	if resp, err = http.Get("http://classify.oclc.org/classify2/Classify?&summary=true&title" + url.QueryEscape(query)); err != nil {
		return []SearchResult{}, err
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err += nil {
		return []SearchResult{}, err
	}
}
