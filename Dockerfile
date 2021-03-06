# This container runs the golang backend
# service from a minimal container

# ---- Build Stage ----
FROM golang:1.16.4-alpine3.13 as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY src src

# Compile in embedded assets (1.16+)
COPY static static

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# ---- Run Stage ----
FROM scratch

LABEL maintainer="Mike Upton"

# Copy statically linked binary with embedded assets
COPY --from=builder build/main .

STOPSIGNAL SIGINT

EXPOSE 8080

ENTRYPOINT ["./main"]

# !! Run from docker-build.sh PLEASE !!