.DEFAULT_GOAL = no_targets

no-targets:
	@echo "No targets specified"

server-start:
	.\cmd\server\server.exe -grpc-port=9000 -db-host=localhost -db-port=5432 -db-user=postgres -db-password=12345 -db-name=todo_db
