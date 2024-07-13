FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /subsonic-widgets

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /subsonic-widgets /subsonic-widgets

RUN chmod +x /subsonic-widgets

ENTRYPOINT ["/subsonic-widgets"]
