.PHONY: all build run clean

APP_NAME := firewall-list

all: build run

build:
	go build -o ./dist/$(APP_NAME)

run:
	./dist/$(APP_NAME)

clean:
	rm -f ./dist/$(APP_NAME)