FROM golang:1.22.1-alpine3.19 as build

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=build --chown=1001:1001 app .

USER 1001:1001

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT ["./app"]
