#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Go Web框架性能基准测试 - Docker环境   ${NC}"
echo -e "${BLUE}=========================================${NC}"

# 创建结果目录
mkdir -p benchmark_results

# 启动所有服务
echo -e "${GREEN}启动所有框架服务...${NC}"
docker-compose up -d --build nethttp gin hertz gozero kratos

# 等待所有服务完全启动
echo -e "${YELLOW}等待服务启动完成 (10s)...${NC}"
sleep 10

# 检查各服务是否正常运行
echo -e "${GREEN}检查服务状态...${NC}"
docker-compose ps

# 执行基准测试
echo -e "${GREEN}开始执行基准测试...${NC}"
docker-compose -f docker-compose.test.yml up --build

# 输出结果位置
echo -e "${GREEN}测试完成! 结果保存在 benchmark_results 目录${NC}"
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}   Go Web框架性能基准测试结束   ${NC}"
echo -e "${BLUE}=========================================${NC}" 