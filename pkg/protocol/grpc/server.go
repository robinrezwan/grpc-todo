package grpc

import (
	"context"
	"github.com/robinrezwan/grpc-todo/pkg/gen/proto/go/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

// RunServer runs gRPC service to publish To-Do service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, gRPCPort string) error {
	// create server
	server := grpc.NewServer()

	v1.RegisterToDoServiceServer(server, v1API)

	// register server reflection
	reflection.Register(server)

	// shutdown gracefully
	ch := make(chan os.Signal, 1)

	// get interrupt signal
	signal.Notify(ch, os.Interrupt)

	go func() {
		for range ch {
			// shut down gRPC server
			log.Println("Shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("Starting gRPC server...")

	listen, err := net.Listen("tcp", ":"+gRPCPort)

	if err != nil {
		return err
	}

	return server.Serve(listen)
}
