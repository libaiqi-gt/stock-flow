# 1. 构建阶段
FROM golang:alpine AS builder

# 设置 Go 代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
# CGO_ENABLED=0 确保生成静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 2. 运行阶段 (最终镜像)
FROM alpine:latest
WORKDIR /root/

# 从构建阶段把二进制文件拷过来
COPY --from=builder /app/main .
# 如果你有配置文件(config.yaml)，记得取消下面这行的注释
# COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./main"]