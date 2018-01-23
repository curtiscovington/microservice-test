package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/curtiscovington/curtiscovington.com/test/pb"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Test(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	log.Println("I'm here")
	return &pb.Result{Title: "test", Url: "tets", Snippet: "hsdf"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":36061")
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	g := grpc.NewServer()

	pb.RegisterTestServer(g, new(server))
	g.Serve(lis)
}
