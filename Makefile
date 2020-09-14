fmt:
	go fmt *.go
make-migrate:
	go build -o migrate cmd/make-db/main.go
migrate: make-migrate
	./migrate
build:
	go build -o runserver cmd/fox/main.go
run: build
	./runserver

