.PHONY: swagger start

swagger:
	swag init --parseDependency --parseInternal

start:
	go run main.go
