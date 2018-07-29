package main

import (
	"context"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/marthjod/binquiry-experimental/pkg/dispatcher"
	pb "github.com/marthjod/binquiry-experimental/word"
	"github.com/marthjod/binquiry-experimental/wordtype"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

type server struct {
	dispatcher dispatcher.Dispatcher
}

func (s *server) Dispatch(ctx context.Context, in *pb.DispatchRequest) (*pb.DispatchResponse, error) {
	log.Info("received dispatch request")

	word, err := s.dispatcher.Dispatch(in.Word)
	return &pb.DispatchResponse{
		Word: word,
	}, err
}

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	nounParser := os.Getenv("NOUNPARSER")
	if len(nounParser) == 0 {
		log.Fatalln("NOUNPARSER not set")
	}

	parsers := map[wordtype.WordType]string{
		wordtype.WordType_Noun: nounParser,
	}

	srv := &server{
		dispatcher: dispatcher.NewDispatcher(parsers),
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: ", err)
	}
	s := grpc.NewServer()
	pb.RegisterDispatcherServer(s, srv)
	reflection.Register(s)
	log.Infof("listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}
