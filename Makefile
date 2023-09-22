db_down:
	docker compose down 

db_up:
	docker compose up -d

tests:
	go test ./... -v