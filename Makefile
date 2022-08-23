
run:
	go run ./main.go "ping google.com" "ping yandex.ru" "date" "exit 1" --timeout=1 --exit

build:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build

install: build
	cp ./parallel $(GOPATH)/bin/parallel