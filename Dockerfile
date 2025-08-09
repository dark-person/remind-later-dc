FROM golang:1.24.5 AS builder

# Set destination for COPY
ARG CGO_ENABLED=0
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build go executable
RUN go build

# Multi-stage build to reduce size ========================
FROM alpine:3.22.1
COPY --from=builder /app/remind-later-dc /remind-later-dc
ENTRYPOINT ["/remind-later-dc"]
