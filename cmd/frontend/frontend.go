package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatalln("env var PORT unset")
	}

	http.HandleFunc("/", queryHandler)
	log.Printf("listening on :%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("query: %q\n", strings.TrimLeft(r.URL.Path, "/"))
}
