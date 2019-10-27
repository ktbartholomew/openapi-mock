FROM golang:1.13 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

FROM alpine:latest

RUN apk update && apk add ca-certificates
COPY --from=builder /app/openapi-mock /usr/bin/openapi-mock

ENTRYPOINT ["/usr/bin/openapi-mock"]