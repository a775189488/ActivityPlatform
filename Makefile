.PHONY: build clean unittest

build:
		mkdir ./release
		go build -o ./release/server cmd/server/main.go
		go build -o ./release/tool cmd/tool/tool.go

clean:
		rm -rf ./release


unittest:
		go test -v ./internal