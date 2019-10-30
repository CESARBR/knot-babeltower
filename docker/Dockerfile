# First stage: build the application
FROM golang:1.13.0-alpine AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/CESARBR/knot-babeltower
COPY . .

# Install build utilities and glide dependency manager
RUN apk add --update --no-cache --virtual .build-deps \
    make \
    git \
    curl
RUN curl https://glide.sh/get | sh

# Install project development tools and dependencies
RUN make tools
RUN make deps

# Build the application
RUN make bin

# Remove build dependencies
RUN apk del .build-deps

# Second stage: create the entrypoint to the application binary generated in the previous stage
FROM alpine:3.10
WORKDIR /root/

# Copy the configuration files from the build stage
RUN mkdir internal
COPY --from=builder /go/src/github.com/CESARBR/knot-babeltower/internal/ ./internal

# Copy the binary file from the build stage
COPY --from=builder /go/src/github.com/CESARBR/knot-babeltower/app-linux-amd64 app
ENTRYPOINT ["./app"]
