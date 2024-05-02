.PHONY: swagger start new

swagger:
	swag init --parseDependency --parseInternal

start:
	go run main.go
new:
	swag init --parseDependency --parseInternal && go run main.go
connect-db:
	./scripts/connect-db.sh
docker-start:
	docker-compose up --build -d
docker-stop:
	docker-compose down
connect-datn-server:
	ssh chiennd@103.166.185.48 -p 2222