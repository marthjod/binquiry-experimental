package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-new/pkg/getter"
	"github.com/marthjod/binquiry-new/pkg/model/noun"
	pb "github.com/marthjod/binquiry-new/pkg/model/noun"
	"github.com/marthjod/binquiry-new/pkg/model/wordtype"
	"github.com/marthjod/binquiry-new/pkg/reader"
	"google.golang.org/grpc"
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
				// Set up a connection to the server.
				conn, err := grpc.Dial(nounParser, grpc.WithInsecure())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
				defer conn.Close()
				c := pb.NewNounParserClient(conn)

				// Contact the server and print out its response.
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := c.Parse(ctx, &pb.ParseRequest{Word: resp})
				if err != nil {
					log.Fatalln(err)
				}
				log.Debugf("Response: %s", r.Json)
				w, err := noun.FromJSON(r.Json)
				if err != nil {
					log.Fatalln(err)
				}

				// req, err := http.NewRequest(http.MethodGet, nounParser, bytes.NewReader(resp))
				// if err != nil {
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
				// c := &http.Client{}
				// resp, err := c.Do(req)
				// if err != nil {
				// 	log.Error(err)
				// }

				// b, err := ioutil.ReadAll(resp.Body)
				// if err != nil {
				// 	log.Error(err)
				// }
				// defer resp.Body.Close()
				// w, err := noun.FromJSON(b)
				// if err != nil {
				// 	log.Error(err)
				// }
				words = append(words, w)
			}

		}

		fmt.Fprintf(w, "%s\n", words.JSON())

	})
	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
