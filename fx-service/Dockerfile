# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git protobuf-dev

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install github.com/twitchtv/twirp/protoc-gen-twirp@latest

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf code
RUN mkdir -p rpc/fxservice && \
    protoc --proto_path=api/proto \
    --go_out=rpc/fxservice \
    --twirp_out=rpc/fxservice \
    --go_opt=paths=source_relative \
    --twirp_opt=paths=source_relative \
    api/proto/fxservice.proto && \
    go mod tidy

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/bin/server ./cmd/server

# Runtime stage for server
FROM alpine:latest AS server

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/bin/server .

EXPOSE 4000

CMD ["./server"]

