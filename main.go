package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"portal_to_the_news/config"
	"portal_to_the_news/parse"
	"strconv"
)

var tpl = template.Must(template.ParseFiles("pages/index.html"))

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	serveMux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	serveMux.Handle("/static/", http.StripPrefix("/static/", fs))
	serveMux.HandleFunc("/", indexHandler)
	serveMux.HandleFunc("/search", searchHandler)

	log.Println("Begin listening on http://localhost:" + port)

	err := http.ListenAndServe(":"+port, serveMux)

	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500 | Internal server error"))
		return
	}

	params := u.Query()

	searchKey := params.Get("q")
	page := params.Get("page")

	if page == "" {
		page = "1"
	}

	search := &parse.Search{}
	search.SearchKey = searchKey

	next, err := strconv.Atoi(page)

	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "500 | Internal server error", http.StatusInternalServerError)
		return
	}

	search.NextPage = next
	pageSize := 20
	endpoint := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en",
		url.QueryEscape(search.SearchKey),
		pageSize,
		search.NextPage,
		config.API_KEY,
	)

	response, err := http.Get(endpoint)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&search.Results)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / pageSize)))
	err = tpl.Execute(w, search)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
