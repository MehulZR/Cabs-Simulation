FROM golang:1.24 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 make build

FROM scratch

COPY --from=builder /usr/src/app/bin/processor /usr/src/app/bin/event-generator-sim /usr/src/app/bin/location-updator-sim ./

