# ------------------------------
# 构建阶段（缓存化）
# ------------------------------
FROM golang:1.23-alpine AS backend-builder

# 安装构建依赖
RUN apk add --no-cache git build-base liblz4-dev bash ca-certificates

# 设置 Go 代理与缓存目录
ENV GOPROXY=https://goproxy.cn,direct
ENV GOCACHE=/go/cache
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod .
COPY go.sum .

# 下载依赖（会写入缓存目录）
RUN go mod tidy && go mod download

# 复制源码
COPY . .

ARG VERSION=unknown
# 构建可执行文件
RUN go build -mod=readonly -ldflags "-X main.version=$VERSION" -o coroot .

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
COPY --from=backend-builder /app/coroot /usr/bin/coroot

RUN mkdir /data && chown 65534:65534 /data

USER 65534:65534
VOLUME /data
EXPOSE 8080

ENTRYPOINT ["/usr/bin/coroot"]
      