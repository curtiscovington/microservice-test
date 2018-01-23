package main

import (
	"log"
	"net"

	test "github.com/curtiscovington/curtiscovington.com/test/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	testClient test.TestClient
}

func (s *server) Test(ctx context.Context, req *test.Request) (*test.Result, error) {
	log.Println("I'm here")
	return s.testClient.Test(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":36060")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := new(server)

	conn, err := grpc.Dial("localhost:36061", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	s.testClient = test.NewTestClient(conn)

	g := grpc.NewServer()
	test.RegisterTestServer(g, s)
	g.Serve(lis)

}
