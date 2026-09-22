# syntax=docker/dockerfile:1
ARG GOTOOLCHAIN=go1.26.8
FROM golang:1.26-bookworm AS gobuild

WORKDIR /src
COPY gosrc/go.mod gosrc/go.sum ./
RUN go mod download github.com/sagernet/sing-box github.com/cloudflare/cloudflared || true
COPY gosrc/ ./
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -tags with_quic -trimpath -ldflags "-s -w" -o /out/niccore ./niccore && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w" -o /out/niclink ./niclink

FROM node:20-alpine

RUN apk add --no-cache ca-certificates curl openssl unzip \
  && update-ca-certificates

WORKDIR /app
COPY package.json index.js README.md LICENSE ./
COPY public/ ./public/
COPY --from=gobuild /out/niccore /out/niclink /app/.bin/

RUN mkdir -p /app/.bin/.run && chown -R node:node /app \
  && chmod +x /app/.bin/niccore /app/.bin/niclink

USER node

ENV PORT=3000 \
    BIN_DIR=/app/.bin \
    BIN_TTL_SEC=0 \
    AT_LINK_MODE=temp \
    WS_PATH=/link \
    VLESS_PORT=18000 \
    LOG_LEVEL=warn \
    GOGC=20 \
    UV_THREADPOOL_SIZE=2 \
    NODE_OPTIONS="--max-old-space-size=96"

EXPOSE 3000

CMD ["node", "index.js"]
