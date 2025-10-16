# 构建阶段：使用官方 golang 镜像编译
FROM golang:1.23-bullseye AS backend-builder
RUN apt-get update && apt-get install -y --no-install-recommends liblz4-dev build-essential pkg-config ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /tmp/src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=unknown
RUN go build -mod=readonly -ldflags "-X main.version=$VERSION" -o /tmp/coroot .

# 运行阶段：使用 Debian slim 作为最小运行时镜像
FROM debian:bookworm-slim

ARG VERSION=unknown
LABEL name="coroot" \
      vendor="Coroot, Inc." \
      maintainer="Coroot, Inc." \
      version=${VERSION} \
      release="1" \
      summary="Coroot Community Edition." \
      description="Coroot Community Edition container image."

COPY LICENSE /licenses/LICENSE

# 从构建阶段拷贝已编译的二进制
COPY --from=backend-builder /tmp/coroot /usr/bin/coroot

# 安装 lz4 运行时库，避免二进制运行时报错
RUN apt-get update && apt-get install -y --no-install-recommends liblz4-1 ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /data && chown 65534:65534 /data

USER 65534:65534
VOLUME /data
EXPOSE 8080

ENTRYPOINT ["/usr/bin/coroot"]