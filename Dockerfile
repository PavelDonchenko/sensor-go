# Step 1: Modules caching
FROM golang:1.20-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.20-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/sensor ./cmd/sensor

# Step 3: Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates

EXPOSE 5000

# GOPATH for scratch images is /
COPY --from=builder /app/config.yaml /
COPY --from=builder /bin/sensor /sensor
CMD ["/sensor"]