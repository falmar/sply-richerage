FROM golang:1.19-alpine as builder
WORKDIR /go-app
COPY go.mod go.sum ./
RUN go mod download
COPY /cmd /go-app/cmd
COPY /internal /go-app/internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/main ./cmd && \
    chmod +x ./bin/main

FROM alpine:3.17 as http
ENV PORT=80
ENV DEBUG=1
ENV LOG_LEVEL=DEBUG
COPY --from=builder /go-app/bin/main /main
ENTRYPOINT ["/main", "http"]
