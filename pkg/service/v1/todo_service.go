package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/robinrezwan/grpc-todo/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API being served
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	v1.UnimplementedToDoServiceServer
	db *sql.DB
}

// NewToDoServiceServer creates To-Do service
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: Service implements API version %s, but asked for %s", apiVersion, api)
		}
	}

	return nil
}

// connectDB returns SQL database connection from the pool
func (s *toDoServiceServer) connectDB(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to connect to database: "+err.Error())
	}

	return c, nil
}

// Create new To-Do
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check API version
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connectDB(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert To-Do entity data
	_, err = c.ExecContext(ctx, "INSERT INTO todo(title, description) VALUES($1, $2)",
		req.ToDo.Title, req.ToDo.Description)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to insert into DB: "+err.Error())
	}

	// TODO: get ID of created To-Do
	id := int64(0)

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read to-do task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get SQL connection from pool
	c, err := s.connectDB(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	// query To-Do by ID
	rows, err := c.QueryContext(ctx, "SELECT id, title, description FROM todo WHERE id=$1", req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to select from DB: "+err.Error())
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to retrieve data from DB: "+err.Error())
	}

	if !rows.Next() {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	// get To-Do data
	var todo v1.ToDo

	if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve field values from row: "+err.Error())
	}

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Reminder field has invalid format: "+err.Error())
	}

	if rows.Next() {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Found multiple ToDo rows with ID: %d", req.Id))
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &todo,
	}, nil
}
