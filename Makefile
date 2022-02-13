.PHONY=run
run:
	go run main.go

.PHONY=build
build:
	go build -o proc-peepin main.go

.PHONY=build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o proc-peepin main.go
