.DEFAULT_GOAL = no_targets

no-targets:
	@echo "No targets specified"

buf-gen:
	buf lint && buf generate

doc-gen:
	protoc --doc_out=./api/docs --doc_opt=html,doc.html api/todo/v1/*.proto

server-start:
	cd cmd\server &&\
	go build . &&\
	server.exe -grpc-port=9000 -db-host=localhost -db-port=5432 -db-user=postgres -db-password=12345 -db-name=todo_db

all: | buf-gen doc-gen server-start

ui-start:
	grpcui -plaintext localhost:9000
