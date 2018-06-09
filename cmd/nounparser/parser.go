package main

import (
	"bytes"
	"context"
	"log"
	"net"

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
	header, _, xmlRoot, err := reader.Read(bytes.NewReader(in.Word))
	if err != nil {
		return nil, err
	}
	path := xmlpath.MustCompile("//tr/td[2]")
	word := noun.ParseNoun(header, path.Iter(xmlRoot))
	return &pb.ParseResponse{Noun: word}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNounParserServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
