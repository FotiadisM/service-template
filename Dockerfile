FROM golang:1.24-alpine3.21 AS build

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/book-svc

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build --chown=1001:1001 app .

USER 1001:1001

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["./app"]
