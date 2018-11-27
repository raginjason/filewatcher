FROM golang:1.11 as builder

WORKDIR /build

# Get dependencies setup first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy rest of application in place and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...
RUN go test ./...

# Bare minimum container
FROM scratch
WORKDIR /app
COPY --from=builder /build/main .

CMD ["/app/main"]
