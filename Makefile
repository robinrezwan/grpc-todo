.DEFAULT_GOAL = no_targets

no-targets:
	@echo "No targets specified"

buf-gen:
	buf generate

server-start:
	cd cmd\server &&\
	go build . &&\
	server.exe -grpc-port=9000 -db-host=localhost -db-port=5432 -db-user=postgres -db-password=12345 -db-name=todo_db

all: | buf-gen server-start
