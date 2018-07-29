//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../noun/noun.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../gender/gender.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../number/number.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../case/case.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../wordtype/wordtype.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../word/word.proto

package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/pkg/handler"
	"github.com/marthjod/binquiry-experimental/pkg/logging"
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

	dispatcher := os.Getenv("DISPATCHER")
	if len(dispatcher) == 0 {
		log.Fatalln("DISPATCHER not set")
	}

	hdlr := handler.NewBackendHandler(dispatcher)

	log.Infof("listening on :%s", port)
	log.Fatalln(http.ListenAndServe(":"+port, hdlr))
}
