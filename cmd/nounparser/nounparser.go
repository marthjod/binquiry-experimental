package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/marthjod/binquiy-new/pkg/reader"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/parse", nounHandler)
	http.ListenAndServe(":"+port, nil)
}

func nounHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	header, wordType, xmlRoot, err := reader.Read(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
