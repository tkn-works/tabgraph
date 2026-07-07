BINARY = tabgraph

.PHONY: dev-api dev-web build build-all clean

dev-api:
	go run ./cmd/tabgraph

dev-web:
	cd web && NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev

build:
	cd web && npm run build
	rm -rf ui/out && cp -r web/out ui/out
	go build -o $(BINARY) ./cmd/tabgraph

clean:
	rm -rf ui/out web/out $(BINARY) dist
	cd web && rm -rf .next

build-all: build
	mkdir -p dist
	GOOS=darwin  GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64      ./cmd/tabgraph
	GOOS=darwin  GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64      ./cmd/tabgraph
	GOOS=linux   GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64       ./cmd/tabgraph
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe ./cmd/tabgraph
