.PHONY: swagger start new

swagger:
	swag init --parseDependency --parseInternal

start:
	go run main.go
new:
	swag init --parseDependency --parseInternal && go run main.go
connect-db:
	./scripts/connect-db.sh
