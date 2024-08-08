.PHONY: run
run: run-server run-client

run-server:
	@echo "Run Server Script"
	cd server && go run cmd/api/main.go --config="./config/development"

run-client:
	@echo "Run Client Script"