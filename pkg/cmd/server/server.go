package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/robinrezwan/grpc-todo/pkg/protocol/grpc"
	"github.com/robinrezwan/grpc-todo/pkg/service/v1"
)

// Config is configuration for server
type Config struct {
	// gRPC server parameters section
	GRPCPort string

	// DB parameters section
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var config Config

	flag.StringVar(&config.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&config.DBHost, "db-host", "", "Database host")
	flag.StringVar(&config.DBPort, "db-port", "", "Database port")
	flag.StringVar(&config.DBUser, "db-user", "", "Database user")
	flag.StringVar(&config.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&config.DBName, "db-name", "", "Database schema")

	flag.Parse()

	if len(config.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: %s", config.GRPCPort)
	}

	dataSourceName := fmt.Sprintf("host=%s "+"port=%s "+"user=%s "+"password=%s "+"dbname=%s "+"sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, config.GRPCPort)
}
