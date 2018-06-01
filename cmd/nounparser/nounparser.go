package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marthjod/binquiry-new/pkg/reader"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/parse", nounHandler)
	log.Printf("listening on :%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func nounHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming request")
	header, wordType, _, err := reader.Read(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "header: %q, word type: %q\n", header, wordType)
}
