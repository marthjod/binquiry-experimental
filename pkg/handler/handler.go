package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-new/pkg/getter"
	"github.com/marthjod/binquiry-new/pkg/model/noun"
	pb "github.com/marthjod/binquiry-new/pkg/model/noun"
	"github.com/marthjod/binquiry-new/pkg/model/wordtype"
	"github.com/marthjod/binquiry-new/pkg/reader"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type handler struct {
	logger  *log.Entry
	parsers map[wordtype.WordType]string
}

func NewBackendHandler(parsers map[wordtype.WordType]string) *handler {
	return &handler{
		logger:  &log.Entry{},
		parsers: parsers,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	cid := generateCorrelationID()
	h.logger = log.WithFields(log.Fields{
		"cid":  cid,
		"task": "handler",
	})
	h.logger.Debug("starting")

	query := strings.TrimLeft(r.URL.Path, "/")
	h.logger.WithFields(log.Fields{
		"request": query,
	}).Debug()
	if len(query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	g := getter.NewGetter("http://dev.phpbin.ja.is/ajax_leit.php", &http.Client{Timeout: 5 * time.Second}, cid)
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
	type result struct {
		word wordtype.Word
		err  error
	}
	resultChan := make(chan result, len(res))

	for _, resp := range res {
		go func(r []byte) {
			w, err := dispatch(r, h.parsers)
			resultChan <- result{
				word: w,
				err:  err,
			}
		}(resp)
	}

	h.logger.Debugf("expecting %d result(s)", cap(res))
	for i := 0; i < cap(res); i++ {
		select {
		case result := <-resultChan:
			if result.err != nil {
				h.logger.WithFields(log.Fields{
					"error": result.err,
				}).Error()
				break
			}
			h.logger.WithFields(log.Fields{
				"result": result.word,
			}).Debug()
			words = append(words, result.word)
		}
	}

	fmt.Fprintf(w, "%s\n", words.JSON())
	end := time.Now()
	h.logger.WithField("duration", end.Sub(start)).Debug("done")
}

func dispatch(in []byte, parsers map[wordtype.WordType]string) (wordtype.Word, error) {
	_, wordType, _, err := reader.Read(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	switch wordType {
	case wordtype.Noun:
		return parseNoun(in, parsers[wordtype.Noun])
	default:
		return nil, errors.New("not implemented yet")
	}
}

func generateCorrelationID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error("unable to generate UUID")
		// rather than failing, run with degraded service
		return "no-uuid"
	}
	return uuid.String()
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
