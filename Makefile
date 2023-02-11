build:
	GOOS=linux GOARCH=amd64 go build -o main cmd/main.go
	zip main.zip main
	rm main
