
.PHONY: build
build: build/dockerbox
	docker build -t gliderlabs/dockerbox .

build/dockerbox:
	mkdir -p build
	GOOS=linux go build -o ./build/dockerbox ./cmd/dockerbox/...
