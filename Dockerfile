FROM golang:1.19-alpine AS build
ADD . /dnsproxy
ENV CGO_ENABLED=0
WORKDIR /dnsproxy
RUN GOOS=linux GOARCH=amd64 go build -o dnsproxy.bin ./cmd

FROM alpine:latest
COPY --from=build /dnsproxy/dnsproxy.bin /dnsproxy/dnsproxy.bin
ENTRYPOINT ["/dnsproxy/dnsproxy.bin"]
