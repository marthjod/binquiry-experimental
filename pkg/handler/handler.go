package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/pkg/getter"
	"github.com/marthjod/binquiry-experimental/word"
	pb "github.com/marthjod/binquiry-experimental/word"

	uuid "github.com/satori/go.uuid"
)

// Handler represents an http.Handler (frontend) to a range of (backend) parsers.
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type backendHandler struct {
	logger            *log.Entry
	dispatcherAddress string
}

// NewBackendHandler returns a pre-configured Handler.
func NewBackendHandler(dispatcherAddress string) Handler {
	return &backendHandler{
		logger:            &log.Entry{},
		dispatcherAddress: dispatcherAddress,
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

	var words = word.Words{}
	type result struct {
		word *word.Word
		err  error
	}
	resultChan := make(chan result, len(res))

	for idx, resp := range res {
		h.logger.Debugf("dispatch #%d", idx)
		go func(r []byte) {
			w, err := h.dispatch(r)
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

func (h *backendHandler) dispatch(in []byte) (*word.Word, error) {
	conn, err := grpc.Dial(h.dispatcherAddress, grpc.WithInsecure())
	if err != nil {
		h.logger.WithFields(log.Fields{
			"error": err,
		}).Error("did not connect")
		return &word.Word{}, err
	}
	defer conn.Close()
	c := pb.NewDispatcherClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.Dispatch(ctx, &pb.DispatchRequest{Word: in})
	if err != nil {
		h.logger.WithFields(log.Fields{
			"error": err,
		}).Error("could not dispatch")
		return &word.Word{}, err
	}
	h.logger.WithFields(log.Fields{
		"response": r.Word.CanonicalForm(),
	}).Debug()
	return r.Word, nil
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
