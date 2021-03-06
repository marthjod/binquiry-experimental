//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../noun/noun.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../gender/gender.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../number/number.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../case/case.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../wordtype/wordtype.proto

package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/pkg/handler"
	"github.com/marthjod/binquiry-experimental/pkg/logging"
	"github.com/marthjod/binquiry-experimental/wordtype"
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
		wordtype.WordType_Noun: nounParser,
	}
	hdlr := handler.NewBackendHandler(parsers)

	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, hdlr))
}
