FROM golang:latest AS build
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine:latest
RUN apk add --no-cache mysql-client

WORKDIR /app
COPY --from=build /app/build/backup_app /app/backup_app
RUN mkdir -p /app/backups && chmod 777 /app/backups

ENTRYPOINT ["/app/backup_app"]
CMD ["--help"]
