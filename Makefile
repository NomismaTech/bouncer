run:
	docker run \
		-it \
		--name go \
		-v `pwd`:/go/src/github.com/palantir/bouncer \
		--workdir /go/src/github.com/palantir/bouncer \
		--rm \
		nmiyake/go:go-darwin-linux-no-cgo-1.11.4-t144 \
		/bin/bash

build:
	go build -o pkg/bouncer main/main.go
	chmod +x pkg/bouncer
