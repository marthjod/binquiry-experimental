package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/marthjod/binquiry-experimental/pkg/model/noun"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	xmlpath "gopkg.in/xmlpath.v2"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatalln("PORT not set")
	}

	logLevel := os.Getenv("LOGLEVEL")
	if len(logLevel) == 0 {
		logLevel = "info"
	}

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetLevel(lvl)

	http.HandleFunc("/parse", nounHandler)
	log.Printf("listening on :%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func nounHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("incoming request")
	header, _, xmlRoot, err := reader.Read(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	path := xmlpath.MustCompile("//tr/td[2]")
	word := noun.ParseNoun(header, path.Iter(xmlRoot))

	log.Debug(word)
	fmt.Fprintf(w, word.JSON())
}
