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
	pb "github.com/marthjod/binquiry-experimental/noun"
	"github.com/marthjod/binquiry-experimental/pkg/getter"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	"github.com/marthjod/binquiry-experimental/wordtype"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

// Handler represents an http.Handler (frontend) to a range of (backend) parsers.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type backendHandler struct {
	logger  *log.Entry
	parsers map[wordtype.WordType]string
}

// NewBackendHandler returns a pre-configured Handler.
func NewBackendHandler(parsers map[wordtype.WordType]string) Handler {
	return &backendHandler{
		logger:  &log.Entry{},
		parsers: parsers,
	}
}

func (h *backendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	for idx, resp := range res {
		h.logger.Debugf("dispatch #%d", idx)
		go func(r []byte) {
			w, err := dispatch(r, h.parsers)
			resultChan <- result{
				word: w,
				err:  err,
			}
		}(resp)
	}

	for i := 0; i < len(res); i++ {
		select {
		case result := <-resultChan:
			if result.err != nil {
				h.logger.WithFields(log.Fields{
					"error": result.err,
				}).Error()
				break
			}
			h.logger.WithFields(log.Fields{
				"result": result.word.CanonicalForm(),
			}).Debugf("received #%d", i)
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
	case wordtype.WordType_Noun:
		return parseNoun(in, parsers[wordtype.WordType_Noun])
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
	// log.Debugf("response: %s", r)

	return r.Noun, nil
}
