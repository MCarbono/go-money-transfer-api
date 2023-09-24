db_down:
	docker compose down 

db_up:
	docker compose up -d

tests:
	go test ./test -v

run:
	go run main.go

build:
	docker-compose -f docker-compose.production.yml build

run_prod:
	docker-compose -f docker-compose.production.yml up -d

down:
	docker-compose -f docker-compose.production.yml down
