package main

import (
	"context"
	"flag"
	"github.com/robinrezwan/grpc-todo/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connct to server: %v", err)
	}

	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Call Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title:       "Make a coffee",
			Description: "With more sugar",
		},
	}

	res1, err := c.Create(ctx, &req1)

	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}

	res2, err := c.Read(ctx, &req2)

	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}

	log.Printf("Read result: <%+v>\n\n", res2)
}
