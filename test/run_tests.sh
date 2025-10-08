#!/bin/bash

# 运行所有测试并生成覆盖率报告
run_tests_with_coverage() {
    echo "正在运行所有测试并生成覆盖率报告..."
    
    # 确保测试数据库环境变量已设置
    export TEST_DB_HOST=${TEST_DB_HOST:-localhost}
    export TEST_DB_PORT=${TEST_DB_PORT:-3306}
    export TEST_DB_USERNAME=${TEST_DB_USERNAME:-root}
    export TEST_DB_PASSWORD=${TEST_DB_PASSWORD:-123456}
    
    # 创建覆盖率目录
    mkdir -p coverage
    
    # 运行所有测试并生成覆盖率报告
    go test -v -coverprofile=coverage/coverage.out ./test/...
    
    # 检查测试是否成功
    if [ $? -ne 0 ]; then
        echo "测试失败！"
        exit 1
    fi
    
    # 生成HTML格式的覆盖率报告
    echo "正在生成HTML覆盖率报告..."
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    
    # 显示覆盖率摘要
    echo "\n覆盖率摘要："
    go tool cover -func=coverage/coverage.out | grep total
    
    echo "\n测试完成！覆盖率报告已保存到 coverage/coverage.html"
}

# 运行单个测试文件
run_single_test() {
    local test_file=$1
    echo "正在运行单个测试文件：$test_file"
    
    # 设置测试数据库环境变量
    export TEST_DB_HOST=${TEST_DB_HOST:-localhost}
    export TEST_DB_PORT=${TEST_DB_PORT:-3306}
    export TEST_DB_USERNAME=${TEST_DB_USERNAME:-root}
    export TEST_DB_PASSWORD=${TEST_DB_PASSWORD:-123456}
    
    # 运行指定的测试文件
    go test -v ./test/$test_file
    
    if [ $? -ne 0 ]; then
        echo "测试失败！"
        exit 1
    fi
}

# 显示帮助信息
display_help() {
    echo "用法：$0 [选项]"
    echo "选项："
    echo "  all        运行所有测试并生成覆盖率报告"
    echo "  single     运行单个测试文件（需要指定文件路径）"
    echo "  help       显示此帮助信息"
    echo ""
    echo "示例："
    echo "  $0 all                 # 运行所有测试并生成覆盖率报告"
    echo "  $0 single pkg/database_test.go  # 运行单个测试文件"
}

# 根据参数执行不同的功能
case "$1" in
    all)
        run_tests_with_coverage
        ;;
    single)
        if [ -n "$2" ]; then
            run_single_test "$2"
        else
            echo "错误：请指定要运行的测试文件路径"
            echo "示例：$0 single pkg/database_test.go"
            exit 1
        fi
        ;;
    help | *)
        display_help
        ;;
esac