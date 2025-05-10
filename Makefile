run-dev:
	go run cmd/server/main.go

race-detector:
	go run -race cmd/server/main.go

generate-unix-timestamp:
	date +%s

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

generate-random-key:
	openssl rand -base64 32

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-steps-up-one:
	go run cmd/migrate/main.go steps 1

migrate-steps-down-one:
	go run cmd/migrate/main.go steps -1

swagger:
	swag init --generalInfo cmd/server/main.go --output docs