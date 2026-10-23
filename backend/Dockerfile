FROM golang:1.27.3-alpine@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4
ARG APK_MIRROR=https://mirrors.aliyun.com/alpine
ARG GO_BUILD_PARALLELISM=8
ARG GO_BUILD_MAX_PROCS=8
ARG GO_BUILD_MEMORY_LIMIT=6GiB

WORKDIR /app

# 安装必要的工具
RUN sed -i "s#https://dl-cdn.alpinelinux.org/alpine#${APK_MIRROR}#g" /etc/apk/repositories && \
    apk add --no-cache git

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN VERSION_VALUE="$(./scripts/resolve-version.sh)" && \
    GOMAXPROCS="${GO_BUILD_MAX_PROCS}" GOMEMLIMIT="${GO_BUILD_MEMORY_LIMIT}" \
    go build -p "${GO_BUILD_PARALLELISM}" -ldflags="-s -w -X main.Version=${VERSION_VALUE}" -o main ./cmd/server/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
