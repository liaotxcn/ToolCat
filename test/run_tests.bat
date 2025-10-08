@echo off

REM 设置默认的测试数据库环境变量
if "%TEST_DB_HOST%"=="" set TEST_DB_HOST=localhost
if "%TEST_DB_PORT%"=="" set TEST_DB_PORT=3306
if "%TEST_DB_USERNAME%"=="" set TEST_DB_USERNAME=root
if "%TEST_DB_PASSWORD%"=="" set TEST_DB_PASSWORD=123456

REM 创建覆盖率目录
if not exist coverage mkdir coverage

REM 根据参数执行不同的功能
if "%1"=="all" (
    echo 正在运行所有测试并生成覆盖率报告...
    go test -v -coverprofile=coverage\coverage.out ./test/...
    
    if %errorlevel% neq 0 (
        echo 测试失败！
        exit /b 1
    )
    
    echo 正在生成HTML覆盖率报告...
    go tool cover -html=coverage\coverage.out -o coverage\coverage.html
    
    echo.
    echo 覆盖率摘要：
    go tool cover -func=coverage\coverage.out | findstr total
    
    echo.
    echo 测试完成！覆盖率报告已保存到 coverage\coverage.html
    echo 请在浏览器中打开 coverage\coverage.html 查看详细报告
    
) else if "%1"=="single" (
    if "%2"=="" (
        echo 错误：请指定要运行的测试文件路径
        echo 示例：%0 single pkg\database_test.go
        exit /b 1
    )
    
    echo 正在运行单个测试文件：%2
    go test -v ./test/%2
    
    if %errorlevel% neq 0 (
        echo 测试失败！
        exit /b 1
    )
    
) else (
    echo 用法：%0 [选项]
    echo 选项：
    echo   all        运行所有测试并生成覆盖率报告
    echo   single     运行单个测试文件（需要指定文件路径）
    echo   help       显示此帮助信息
    echo.
    echo 示例：
    echo   %0 all                 运行所有测试并生成覆盖率报告
    echo   %0 single pkg\database_test.go  运行单个测试文件
)

pause