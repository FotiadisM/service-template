FROM alpine:3.21
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY bin /app/bin

EXPOSE 8080

ENTRYPOINT ["/app/bin/app"]
