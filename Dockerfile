FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY *.go ./ /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /playlist-manager

# Run tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

# Install CA certificates
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*

WORKDIR /

COPY --from=build-stage /playlist-manager /playlist-manager

EXPOSE 8080

CMD ["/playlist-manager"]
