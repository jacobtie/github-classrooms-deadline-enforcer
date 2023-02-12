build:
	GOOS=linux GOARCH=amd64 go build -o main cmd/main.go
	zip main.zip main
	rm main


services:
	go run mock-service-api/main.go

test:
	IS_TEST=true go test ./...

test-gha:
	docker-compose -f docker-compose.test.yml up --build --force-recreate --abort-on-container-exit
