all: clean build run

build:
	go build -o bin/weatherapp .
run: 
	bin/weatherapp
clean:
	go mod tidy
	rm bin/*
