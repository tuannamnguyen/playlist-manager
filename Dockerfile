FROM golang:1.22 AS build-stage

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/api -o /playlist-manager

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

# Install CA certificates
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

WORKDIR /

COPY --from=build-stage /playlist-manager /playlist-manager
COPY ./cmd/api/.env.vault /.env.vault

EXPOSE 8080

CMD ["/playlist-manager"]
