run:
	go run ./... -t 6421080707:AAHHCXxe390fMNsiqUsAtmsj5Zd_njw1640
build:
	go build -v ./...

pgup:
	docker compose up -d

pgdown:
	docker compose down
	
createdb:
	docker exec -it tg-bath-bot-db-1 createdb --username=postgres --owner=postgres articles

dropdb:
	docker exec -it tg-bath-bot-db-1 dropdb articles

mgrup:
	migrate -path migration -database "postgresql://postgres:postgres@localhost:5432/articles?sslmode=disable" -verbose up

mgrdown:
	migrate -path migration -database "postgresql://postgres:postgres@localhost:5432/articles?sslmode=disable" -verbose down

.PHONY: run build pgup pgdown createdb dropdb mgrup mgrdown 