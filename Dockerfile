# 构建阶段
FROM golang:1.23.4-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置国内Go模块代理以加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 复制go.mod文件
COPY go.mod ./

# 复制go.sum文件
COPY go.sum ./

# 下载依赖（使用代理加速）
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -o toolcat main.go

# 运行阶段
FROM alpine:3.18

# 添加非root用户
RUN addgroup -S toolcat && adduser -S toolcat -G toolcat

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/toolcat .

# 复制配置文件和web目录
COPY --chown=toolcat:toolcat config/ ./config/
COPY --chown=toolcat:toolcat web/ ./web/
COPY --chown=toolcat:toolcat plugins/ ./plugins/

# 创建日志目录
RUN mkdir -p /app/logs && chown toolcat:toolcat /app/logs

# 切换到非root用户
USER toolcat

# 暴露应用端口
EXPOSE 8081

# 设置启动命令
CMD ["./toolcat"]