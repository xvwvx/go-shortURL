export PATH := $(GOPATH)/bin:$(PATH)
export GO15VENDOREXPERIMENT := 1
LDFLAGS := -s -w

all: build

build: app

app:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/shortURL_darwin_amd64 ./
	env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o ./build/shortURL_linux_386 ./
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/shortURL_linux_amd64 ./

temp:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/shortURL_linux_amd64 ./

	