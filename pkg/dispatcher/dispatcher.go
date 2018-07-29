package dispatcher

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/marthjod/binquiry-experimental/noun"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	"github.com/marthjod/binquiry-experimental/word"
	"github.com/marthjod/binquiry-experimental/wordtype"
	"google.golang.org/grpc"
)

// Dispatcher knows how to pass off input to the appropriate parser, according to word type.
type Dispatcher interface {
	Dispatch([]byte) (*word.Word, error)
}

type dispatcher struct {
	Parsers map[wordtype.WordType]string
}

// NewDispatcher returns a new dispatcher.
func NewDispatcher(parsers map[wordtype.WordType]string) Dispatcher {
	return &dispatcher{
		Parsers: parsers,
	}
}

// Dispatch passes off input to the appropriate parser, according to word type.
func (d *dispatcher) Dispatch(in []byte) (*word.Word, error) {
	_, wordType, _, err := reader.Read(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	switch wordType {
	case wordtype.WordType_Noun:
		n, err := parseNoun(in, d.Parsers[wordtype.WordType_Noun])
		if err != nil {
			return &word.Word{}, err
		}
		return &word.Word{
			Type: wordtype.WordType_Noun,
			Noun: n,
		}, nil
	default:
		return nil, errors.New("not implemented yet")
	}
}

func parseNoun(in []byte, parser string) (*noun.Noun, error) {
	conn, err := grpc.Dial(parser, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := noun.NewNounParserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Parse(ctx, &noun.ParseRequest{Word: in})
	if err != nil {
		return nil, err
	}
	// log.Debugf("response: %s", r)

	return r.Noun, nil
}
