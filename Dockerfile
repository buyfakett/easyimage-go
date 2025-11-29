ARG ALPINE_VERSION=3.22
ARG GO_VERSION=1.25.4
ARG AUTHOR=buyfakett
ARG SERVER_NAME=easyimage-go

# 支持多平台构建
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# 后端
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG ALPINE_VERSION
ARG GO_VERSION
ARG SERVER_NAME
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache \
    build-base \
    libwebp-dev \
    libheif-dev \
    git

COPY . .

RUN go mod download

# 根据平台推导出 GOOS 和 GOARCH
RUN set -eux; \
    TARGETOS=${TARGETOS:-linux}; \
    TARGETARCH=${TARGETARCH:-amd64}; \
    echo "Building for TARGETOS=${TARGETOS} TARGETARCH=${TARGETARCH}"; \
    CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /app/${SERVER_NAME}

# 最小编译
FROM alpine:${ALPINE_VERSION} AS final

ARG SERVER_NAME

# 安装运行时依赖
RUN apk add --no-cache \
    tzdata \
    libwebp \
    libheif

ENV TERM=xterm-256color

COPY --from=builder /app/${SERVER_NAME} /app/${SERVER_NAME}

WORKDIR /app
EXPOSE 8888
ENTRYPOINT ["/app/easyimage-go"]
