package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
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

	// log.Println("Search Query is: ", searchKey)
	// log.Println("Results page is: ", page)
}
