package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/launcher/example/protos"
)

const (
	address        = "localhost:8088"
	defaultMessage = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := protos.NewExampleControllerClient(conn)

	// Contact the server and print out its response.
	message := defaultMessage
	if len(os.Args) > 1 {
		message = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, &protos.PingRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
