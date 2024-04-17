package main

import (
	"context"
	"flag"
	"github.com/tugberkugurlu/go-grpc-example/spec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")

	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "helloworld.Greeter", "method": "SayHello"}],
		  "waitForReady": true,
"timeout": "1.5s",
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".0s",
			  "MaxBackoff": ".0s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": ["DEADLINE_EXCEEDED"]
		  }
		}]}`
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := spec.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		// Contact the server and print out its response.
		// ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
		// defer cancel()
		r, err := c.SayHello(context.Background(), &spec.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())

		time.Sleep(1 * time.Second)
	}
}
