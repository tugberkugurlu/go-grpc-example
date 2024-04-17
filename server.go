package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tugberkugurlu/go-grpc-example/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/stats"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	spec.UnimplementedGreeterServer
	eo *EvenOdd
}

func (s *server) SayHello(ctx context.Context, in *spec.HelloRequest) (*spec.HelloReply, error) {
	name := in.GetName()
	val := s.eo.Next()
	fmt.Println(val)
	if val%2 != 0 {
		time.Sleep(2 * time.Second)
	} else {
		time.Sleep(1 * time.Second)
	}
	log.Printf("received: %v", name)
	return &spec.HelloReply{Message: "Hello " + name}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	eo := NewEvenOdd()
	// s := grpc.NewServer(grpc.StatsHandler(serverStatsHandler{}))
	s := grpc.NewServer()
	reflection.Register(s)
	spec.RegisterGreeterServer(s, &server{
		eo: eo,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type EvenOdd struct {
	mu       sync.Mutex
	nextEven bool
}

func NewEvenOdd() *EvenOdd {
	return &EvenOdd{
		nextEven: false,
	}
}

func (eo *EvenOdd) Next() int {
	eo.mu.Lock()
	defer eo.mu.Unlock()
	if eo.nextEven {
		eo.nextEven = false
		return 2
	} else {
		eo.nextEven = true
		return 1
	}
}

type serverStatsHandler struct {
}

func (s serverStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	log.Println("--------------------")
	log.Println("TagRPC")
	log.Println(fmt.Sprintf("%v", reflect.TypeOf(info)))
	log.Println("--------------------")
	return ctx
}

func (s serverStatsHandler) HandleRPC(ctx context.Context, stats stats.RPCStats) {
	log.Println("--------------------")
	log.Println("HandleRPC")
	log.Println(fmt.Sprintf("%v", reflect.TypeOf(stats)))
	log.Println("--------------------")
}

func (s serverStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	log.Println("--------------------")
	log.Println("TagConn")
	log.Println("--------------------")
	return ctx
}

func (s serverStatsHandler) HandleConn(ctx context.Context, stats stats.ConnStats) {
	log.Println("--------------------")
	log.Println("HandleConn")
	log.Println(fmt.Sprintf("%v", reflect.TypeOf(stats)))
	log.Println("--------------------")
}
