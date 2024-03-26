swagger:
    swag init --parseDependency --parseInternal

start:
    go run main.go

.PHONY: swagger start
