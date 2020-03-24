# Detect system OS.
ifeq ($(OS),Windows_NT)
    detected_OS := Windows
else
    detected_OS := $(shell sh -c 'uname -s 2>/dev/null || echo not')
endif

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -race
GOGET=$(GOCMD) get


ifeq ($(detected_OS),Windows)
	BINARY_NAME=stock.exe
	BINARY_RACE_NAME=stock_race.exe
else
	BINARY_NAME=stock
	BINARY_RACE_NAME=stock_race
endif

flags="-X 'main.goversion=`go version`' -X 'main.buildstamp=`date --rfc-3339=seconds`' -X main.githash=`git describe --always --long --abbrev=14`"

all: test build
build:	
	$(GOBUILD) -ldflags ${flags} -x   -v -o $(BINARY_NAME)
test:
	$(GOTEST)  -v ./...
clean:
	$(GOCLEAN)
	rm -f engine.log
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f $(BINARY_RACE_NAME)
	rm -f $(BINARY_RACE_NAME).exe
run: build
	./$(BINARY_NAME)
race:
	$(GOBUILD)  -race -ldflags ${flags} -x   -v -o $(BINARY_RACE_NAME)
