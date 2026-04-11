FROM golang:1.25 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd/api -o /playlist-manager

# Deploy the application binary into a lean image
FROM ubuntu:24.04 AS build-release-stage


WORKDIR /

# Install CA certificates & install cURL + dotenvx
RUN apt-get -y update \
    && apt-get install -y ca-certificates \
    && apt-get -y install curl \
    && curl -sfS https://dotenvx.sh/install.sh | sh \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*


COPY --from=build-stage /playlist-manager /playlist-manager
COPY ./cmd/api/.env.prod /.env.prod

EXPOSE 8080

CMD ["dotenvx", "run", "-f", "/.env.prod", "--", "/playlist-manager"]
