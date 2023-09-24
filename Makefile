db_down:
	docker compose down 

db_up:
	docker compose up -d

tests:
	go test ./test -v

run:
	go run main.go