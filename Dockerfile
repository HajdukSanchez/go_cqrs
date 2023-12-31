# Args variables
## Go version to use
ARG GO_VERSION=1.20.4

# ---------- COMPILATION PROCESS ----------
# Go version
FROM golang:${GO_VERSION}-alpine AS builder

# Tell go to use dependecies direct without using any proxy
RUN go env -w GOPROXY=direct
# Install Git to get packages
RUN apk add --no-cache git
# Install security certificates for app
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Define working dir
WORKDIR /src

# Copy files to install dependecites
COPY ./go.mod ./go.sum ./

# Get dependencies
RUN go mod download

# Copy our code
COPY cmd-feed cmd-feed
COPY cmd-search cmd-search
COPY cmd-pusher cmd-pusher
COPY internal internal
COPY service service

# Build app
RUN go install ./...

# ---------- EXECUTION PROCESS ----------
# Alpine version
FROM alpine:3.11

# Define workind directory
WORKDIR /usr/bin

# Copy compilation files
COPY --from=builder /go/bin .
