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
var tplInfo = template.Must(template.ParseFiles("pages/info.html"))

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	serveMux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	serveMux.Handle("/static/", http.StripPrefix("/static/", fs))
	serveMux.HandleFunc("/assets/favicon.ico", faviconHandler)
	serveMux.HandleFunc("/", indexHandler)
	serveMux.HandleFunc("/search", searchHandler)
	serveMux.HandleFunc("/info", infoHandler)

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

	if searchKey == "" {
		err = tpl.Execute(w, parse.Search{SearchKey: ""})
		return
	}

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
		http.Error(w, "500 | Internal server error", http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "500 | Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&search.Results)

	if err != nil {
		http.Error(w, "500 | Internal server error", http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / pageSize)))

	if ok := !search.IsLastPage(); ok {
		search.NextPage++
	}

	err = tpl.Execute(w, search)

	if err != nil {
		http.Error(w, "500 | Internal server error", http.StatusInternalServerError)
		return
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	err := tplInfo.Execute(w, nil)

	if err != nil {
		log.Println(err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/favicon.ico")
}
