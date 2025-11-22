FROM golang:1.24

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        make && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
COPY . .

EXPOSE 8080

CMD ["make", "start"]

