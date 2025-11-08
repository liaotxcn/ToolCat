# Makefile for Weave project

# 设置变量
APP_NAME := weave
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*")
TEST_FLAGS := -v

# 默认目标
all: build

# 构建应用
build:
	@echo "构建应用..."
	go build -o $(APP_NAME) main.go

# 运行测试
test:
	@echo "运行测试..."
	go test $(TEST_FLAGS) ./...

# 运行应用
run:
	@echo "运行应用..."
	go run main.go

# 清理构建产物
clean:
	@echo "清理构建产物..."
	rm -f $(APP_NAME)
	rm -rf ./dist

# 安装依赖
install:
	@echo "安装依赖..."
	go mod download

# 更新依赖
update:
	@echo "更新依赖..."
	go get -u ./...
	go mod tidy

# 格式化代码
fmt:
	@echo "格式化代码..."
	gofmt -s -w $(GO_FILES)

# 检查代码
lint:
	@echo "检查代码..."
	golangci-lint run ./...

# 热重载开发（需要安装gin工具：go install github.com/codegangsta/gin）
watch:
	@echo "启动热重载开发服务器..."
	gin --appPort 8081 run main.go

# 帮助信息
help:
	@echo "可用命令："
	@echo "  make build     - 构建应用"
	@echo "  make test      - 运行测试"
	@echo "  make run       - 运行应用"
	@echo "  make clean     - 清理构建产物"
	@echo "  make install   - 安装依赖"
	@echo "  make update    - 更新依赖"
	@echo "  make fmt       - 格式化代码"
	@echo "  make lint      - 检查代码"
	@echo "  make watch     - 热重载开发"

.PHONY: all build test run clean install update fmt lint watch help