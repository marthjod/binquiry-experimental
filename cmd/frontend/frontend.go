package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-new/pkg/getter"
	"github.com/marthjod/binquiry-new/pkg/reader"
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
	// getterURL := os.Getenv("GETTER")
	// if len(getterURL) == 0 {
	// 	log.Fatalln("GETTER not set")
	// }

	// nounParser := os.Getenv("NOUN_PARSER")
	// if len(nounParser) == 0 {
	// 	log.Fatalln("NOUN_PARSER not set")
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimLeft(r.URL.Path, "/")
		log.Debugf("request: %q", query)
		g := getter.Getter{
			URLPrefix: "http://dev.phpbin.ja.is/ajax_leit.php",
			Client:    &http.Client{Timeout: 5 * time.Second},
		}
		res, err := g.GetWord(query)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, resp := range res {
			header, wordType, _, err := reader.Read(bytes.NewReader(resp))
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			fmt.Fprintf(w, "header: %q, word type: %q\n", header, wordType)
		}

		// TODO determine word type and pass on to right parser

	})
	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
