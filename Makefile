.PHONY: postgres createdb dropdb sqlc test server redis

postgres:
	docker run -d --name clove \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=justlivin@56 \
		-e POSTGRES_DB=clovedb \
		-p 5432:5432 postgres:17-alpine

createdb:
	docker exec -it clove createdb --username=root --owner=root clovedb

dropdb:
	docker exec -it clove dropdb clovedb

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

redis:
	 docker run --name redis -p 6379:6379 -d redis:8-alpine
