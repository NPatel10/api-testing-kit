FROM golang:1.25.0-bookworm AS build

WORKDIR /workspace/server

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

COPY server/go.mod ./
RUN go mod download

COPY server/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY --from=build /go/bin/goose /usr/local/bin/goose
COPY --from=build /out/api /usr/local/bin/api

EXPOSE 8080

CMD ["api"]
