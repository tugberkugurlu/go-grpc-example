package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/tugberkugurlu/go-grpc-example/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	spec.UnimplementedGreeterServer
}

func (*server) SayHello(ctx context.Context, in *spec.HelloRequest) (*spec.HelloReply, error) {
	log.Printf("received: %v", in.GetName())
	return &spec.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	spec.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}