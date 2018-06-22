//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../noun/noun.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../gender/gender.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../number/number.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../case/case.proto
//go:generate protoc -I ../.. --go_out=plugins=grpc:$GOPATH/src ../../wordtype/wordtype.proto

package main

import (
	"bytes"
	"context"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/noun"
	pb "github.com/marthjod/binquiry-experimental/noun"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	xmlpath "gopkg.in/xmlpath.v2"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Parse(ctx context.Context, in *pb.ParseRequest) (*pb.ParseResponse, error) {
	log.Info("received parse request")
	header, _, xmlRoot, err := reader.Read(bytes.NewReader(in.Word))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	path := xmlpath.MustCompile("//tr/td[2]")
	word := noun.ParseNoun(header, path.Iter(xmlRoot))
	log.Infof("returning %q", word.CanonicalForm())
	return &pb.ParseResponse{Noun: word}, nil
}

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}
	s := grpc.NewServer()
	pb.RegisterNounParserServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infof("listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
