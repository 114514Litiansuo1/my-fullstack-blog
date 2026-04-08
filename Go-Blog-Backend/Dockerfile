# 阶段 1: 构建阶段
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# 编译可执行文件
RUN go build -o main .

# 阶段 2: 运行阶段 (最小化镜像)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# 记得拷贝你的图片上传目录(如果需要持久化)
RUN mkdir uploads

# 设置环境变量(从容器外部传入)
ENV password=""

EXPOSE 8080
CMD ["./main"]