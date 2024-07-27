package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", indexHandler)

	log.Println("Begin listening on http://localhost:" + port)

	err := http.ListenAndServe(":"+port, serveMux)

	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("<h1>Welcome!</h1>"))

	// log.Println(r.Method, r.Response.StatusCode)

	if err != nil {
		log.Println(err)
	}
}
