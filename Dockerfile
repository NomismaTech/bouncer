FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/palantir/bouncer

COPY . .

RUN go build -o pkg/bouncer main/main.go

FROM golang:1.11.4-alpine

WORKDIR /work-dir

COPY --from=builder /go/src/github.com/palantir/bouncer/pkg/bouncer .

ENTRYPOINT [ "/work-dir/bouncer" ]
