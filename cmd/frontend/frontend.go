package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-new/pkg/getter"
	"github.com/marthjod/binquiry-new/pkg/model/noun"
	"github.com/marthjod/binquiry-new/pkg/model/wordtype"
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

	nounParser := os.Getenv("NOUNPARSER")
	if len(nounParser) == 0 {
		log.Fatalln("NOUNPARSER not set")
	}

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

		var words = wordtype.Words{}

		for _, resp := range res {
			_, wordType, _, err := reader.Read(bytes.NewReader(resp))
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				return
			}

			switch wordType {
			case wordtype.Noun:
				req, err := http.NewRequest(http.MethodGet, nounParser, bytes.NewReader(resp))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				c := &http.Client{}
				resp, err := c.Do(req)
				if err != nil {
					log.Error(err)
				}

				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err)
				}
				defer resp.Body.Close()
				w, err := noun.FromJSON(b)
				if err != nil {
					log.Error(err)
				}
				words = append(words, w)
			}

		}

		log.Infof("%s", words.JSON())

	})
	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
