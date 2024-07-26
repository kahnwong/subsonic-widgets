FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY templates templates
COPY *.go ./

RUN go build -ldflags "-w -s" -o /subsonic-widgets

FROM alpine:latest AS build-release-stage

WORKDIR /opt

COPY --from=build-stage /subsonic-widgets .

RUN chmod +x subsonic-widgets

ENTRYPOINT ["./subsonic-widgets"]
