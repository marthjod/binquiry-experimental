package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-new/pkg/handler"
	"github.com/marthjod/binquiry-new/pkg/logging"
	"github.com/marthjod/binquiry-new/pkg/model/wordtype"
)

func main() {
	logLevel := os.Getenv("LOGLEVEL")
	if len(logLevel) == 0 {
		logLevel = "info"
	}

	logging.MustSetLoglevel(logLevel)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatalln("PORT not set")
	}

	nounParser := os.Getenv("NOUNPARSER")
	if len(nounParser) == 0 {
		log.Fatalln("NOUNPARSER not set")
	}

	parsers := map[wordtype.WordType]string{
		wordtype.Noun: nounParser,
	}
	hdlr := handler.NewBackendHandler(parsers)

	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, hdlr))
}
