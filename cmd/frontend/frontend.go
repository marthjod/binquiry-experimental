package main

import (
	"bytes"
	"context"
	"errors"
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
		wordChan := make(chan wordtype.Word, len(res))
		errChan := make(chan error, len(res))

		for _, resp := range res {
			go func(r []byte) {
				w, err := dispatch(r, nounParser)
				if err != nil {
					errChan <- err
					return
				}
				log.Debugf("sending to word chan: %s", w)
				wordChan <- w
			}(resp)
		}

		log.Debugf("will select %d time(s)", len(res))
		for i := 0; i < len(res); i++ {
			select {
			case word := <-wordChan:
				log.Debugf("word chan received %s", word)
				words = append(words, word)
			case err := <-errChan:
				log.Debug("got error on error chan")
				log.Error(err)
			}
		}
		fmt.Fprintf(w, "%s\n", words.JSON())

	})
	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func dispatch(in []byte, nounParser string) (wordtype.Word, error) {
	_, wordType, _, err := reader.Read(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	switch wordType {
	case wordtype.Noun:
		return parseNoun(in, nounParser)
	default:
		return nil, errors.New("not implemented yet")
	}
}

func parseNoun(in []byte, parser string) (wordtype.Word, error) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(parser, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewNounParserClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Parse(ctx, &pb.ParseRequest{Word: in})
	if err != nil {
		return nil, err
	}
	// log.Debugf("Response: %s", r.Json)
	w, err := noun.FromJSON(r.Json)
	if err != nil {
		return nil, err
	}

	return w, nil
}
