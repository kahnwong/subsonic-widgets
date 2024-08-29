FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY templates templates
COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /subsonic-widgets

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian11:latest
COPY --from=build /subsonic-widgets /

CMD ["/subsonic-widgets"]
