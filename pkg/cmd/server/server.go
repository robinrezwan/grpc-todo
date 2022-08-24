package server

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
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
	// get configuration
	var cfg Config

	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DBPort, "db-port", "", "Database port")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBName, "db-name", "", "Database name")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: %s", cfg.GRPCPort)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	ctx := context.Background()

	db, err := sql.Open("pgx", dbUrl)

	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
