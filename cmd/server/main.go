package main

import (
	"fmt"
	"github.com/robinrezwan/grpc-todo/pkg/cmd/server"
	"os"
)

func main() {
	if err := server.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
