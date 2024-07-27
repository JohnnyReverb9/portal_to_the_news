package main

import (
	"html/template"
	"log"
	"net/http"
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
