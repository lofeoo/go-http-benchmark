#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Go Web框架性能基准测试开始   ${NC}"
echo -e "${BLUE}=========================================${NC}"

# 确保hey工具已安装
if ! command -v hey &> /dev/null; then
    echo -e "${YELLOW}未找到hey工具，正在安装...${NC}"
    go install github.com/rakyll/hey@latest
fi

# 定义测试参数
CONNECTIONS=2000   # 总请求数
CONCURRENCY=200    # 并发数
DURATION="30s"     # 测试时间
TOOL="hey"         # 使用的工具

# 测试结果目录
RESULTS_DIR="benchmark_results"
mkdir -p $RESULTS_DIR

# 获取当前时间作为测试标识
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 定义所有框架
FRAMEWORKS=("nethttp" "gin" "hertz" "gozero" "kratos")

# 定义所有测试场景
SCENARIOS=("hello" "json" "params" "query" "post")

# 检测运行环境 - Docker或本地
IN_DOCKER=false
if [ -f "/.dockerenv" ]; then
    IN_DOCKER=true
    echo -e "${YELLOW}检测到Docker环境${NC}"
    # 等待服务启动
    echo -e "${YELLOW}等待10秒让其他服务启动...${NC}"
    sleep 10
fi

# 创建测试摘要文件
SUMMARY_FILE="$RESULTS_DIR/summary_$TIMESTAMP.md"
echo "# Go Web框架性能测试结果摘要" > $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "测试时间: $(date)" >> $SUMMARY_FILE
echo "测试参数:" >> $SUMMARY_FILE
echo "- 总请求数: $CONNECTIONS" >> $SUMMARY_FILE
echo "- 并发数: $CONCURRENCY" >> $SUMMARY_FILE
echo "- 持续时间: $DURATION" >> $SUMMARY_FILE
echo "- 测试工具: $TOOL" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

echo "| 框架     | 场景     | 每秒请求数 (RPS) | 平均响应时间 (ms) | P99响应时间 (ms) |" >> $SUMMARY_FILE
echo "|---------|---------|----------------|-----------------|-----------------|" >> $SUMMARY_FILE

# 循环测试每个框架
for framework in "${FRAMEWORKS[@]}"; do
    echo -e "${GREEN}正在测试框架: $framework${NC}"
    
    # 设置主机和端口（Docker环境和本地环境不同）
    if [ "$IN_DOCKER" = true ]; then
        # Docker环境中，使用服务名作为主机名
        HOST=$framework
    else
        # 本地环境，使用localhost
        HOST="localhost"
    fi
    
    # 设置端口
    case $framework in
        "nethttp")
            PORT=8080
            ;;
        "gin")
            PORT=8081
            ;;
        "hertz")
            PORT=8082
            ;;
        "gozero")
            PORT=8083
            ;;
        "kratos")
            PORT=8084
            ;;
    esac
    
    # 检查服务是否已启动（仅在本地环境检查）
    if [ "$IN_DOCKER" = false ]; then
        if ! nc -z $HOST $PORT &>/dev/null; then
            echo -e "${YELLOW}警告: $framework 服务未在端口 $PORT 上运行，请确保服务已启动${NC}"
            continue
        fi
    fi
    
    # 循环测试每个场景
    for scenario in "${SCENARIOS[@]}"; do
        echo -e "  ${BLUE}场景: $scenario${NC}"
        
        # 构建测试URL
        BASE_URL="http://$HOST:$PORT"
        case $scenario in
            "hello")
                URL="$BASE_URL/hello"
                METHOD="GET"
                DATA=""
                ;;
            "json")
                URL="$BASE_URL/json"
                METHOD="GET"
                DATA=""
                ;;
            "params")
                URL="$BASE_URL/params/123"
                METHOD="GET"
                DATA=""
                ;;
            "query")
                URL="$BASE_URL/query?name=测试用户"
                METHOD="GET"
                DATA=""
                ;;
            "post")
                URL="$BASE_URL/post"
                METHOD="POST"
                # 使用普通单引号，避免转义问题
                DATA='{"name":"测试内容","content":"这是一个测试POST请求的内容"}'
                ;;
        esac
        
        # 执行测试
        RESULT_FILE="$RESULTS_DIR/${framework}_${scenario}_$TIMESTAMP.txt"
        
        echo "正在执行测试: $URL"
        if [ "$METHOD" == "GET" ]; then
            hey -n $CONNECTIONS -c $CONCURRENCY -z $DURATION $URL > $RESULT_FILE
        else
            hey -n $CONNECTIONS -c $CONCURRENCY -z $DURATION -m POST -d "$DATA" -T "application/json" $URL > $RESULT_FILE
        fi
        
        # 检查测试是否成功
        if [ $? -ne 0 ]; then
            echo -e "${YELLOW}测试失败，检查错误信息${NC}"
            cat $RESULT_FILE
        fi
        
        # 解析结果
        RPS=$(grep "Requests/sec:" $RESULT_FILE | awk '{print $2}')
        AVG_TIME=$(grep "Average" $RESULT_FILE | awk '{print $2}')
        P99_TIME=$(grep "99%" $RESULT_FILE | awk '{print $2}')
        
        # 添加结果到摘要
        echo "| $framework | $scenario | $RPS | $AVG_TIME | $P99_TIME |" >> $SUMMARY_FILE
        
        # 休息一下，让系统冷却
        sleep 2
    done
    
    echo "" >> $SUMMARY_FILE
done

echo -e "${GREEN}测试完成! 测试结果保存在: $SUMMARY_FILE${NC}"
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Go Web框架性能基准测试结束   ${NC}"
echo -e "${BLUE}=========================================${NC}" 