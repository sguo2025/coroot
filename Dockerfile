# ------------------------------
# 构建阶段
# ------------------------------
FROM golang:1.23-alpine AS backend-builder

# 安装构建依赖
RUN apk add --no-cache git build-base liblz4-dev bash ca-certificates

# 设置 Go 代理与缓存目录
ENV GOPROXY=https://goproxy.cn,direct
ENV GOCACHE=/go/cache
ENV GOPATH=/go

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod .
COPY go.sum .

# 生成 vendor 目录（离线模式）
RUN go mod tidy && go mod vendor

# 复制源码
COPY . .

ARG VERSION=unknown
# 使用 vendor 目录构建，避免每次下载依赖
RUN go build -mod=vendor -ldflags "-X main.version=$VERSION" -o coroot .

# ------------------------------
# 运行阶段
# ------------------------------
FROM registry.access.redhat.com/ubi9/ubi

ARG VERSION=unknown
LABEL name="coroot" \
      vendor="Coroot, Inc." \
      maintainer="Coroot, Inc." \
      version=${VERSION} \
      release="1" \
      summary="Coroot Community Edition." \
      description="Coroot Community Edition container image."

COPY LICENSE /licenses/LICENSE

# 复制可执行文件
COPY --from=backend-builder /app/coroot /usr/bin/coroot

# 创建数据目录
RUN mkdir /data && chown 65534:65534 /data

USER 65534:65534
VOLUME /data
EXPOSE 8080

ENTRYPOINT ["/usr/bin/coroot"]
      