run-dev:
	APP_ENV=dev go run cmd/server/main.go

generate-unix-timestamp:
	date +%s