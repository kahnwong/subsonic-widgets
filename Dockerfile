FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY templates templates
COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /subsonic-widgets

# hadolint ignore=DL3007
FROM alpine:latest AS build-release-stage

WORKDIR /opt/subsonic-widgets

COPY --from=build-stage /subsonic-widgets .

RUN chmod +x subsonic-widgets

ENTRYPOINT ["./subsonic-widgets"]
