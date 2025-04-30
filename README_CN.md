# Go Web框架性能基准测试

本项目提供了一个简单的性能测试套件，用于国内比较流行的Go语言Web框架：
- Hertz (CloudWeGo开源的高性能HTTP框架)
- Gin (流行的轻量级HTTP框架)
- Go-Zero (微服务框架)
- Kratos (B站开源的微服务框架)
- net/http (Go标准库)

本项目使用 Cursor 生成
本项目使用测试工具
- [hey](https://github.com/rakyll/hey)
- [wrk](https://github.com/wg/wrk) 作为基准测试工具。

## 使用方法
- 部署到docker，并执行测试
```docker compose up```
- 结果导出到本地
```docker cp testhertz-benchmark-1:/app/benchmark_results/summary_20250419_022006.md .```

## Cursor Promt
- 文档最后附Cursor Rules设置。
- 项目生成过程，输入以下提示词
>  **制定一个测试Hertz、Gin、Go-zero、Kratos、net/http各框架性能的文案**
- 输出以下内容
> 我来制定一个测试Go语言主流Web框架性能的文案和具体实现方案。 首先，让我探索下这个测试项目应该如何构建。
 Listed 0 items in current directory
我将创建一个完整的Go项目来测试Hertz、Gin、Go-zero、Kratos和标准库net/http的性能。这个项目将公平地比较这些框架在不同场景下的性能表现。 首先创建项目结构：
```
mkdir -p cmd/server/{hertz,gin,gozero,kratos,nethttp} internal/benchmark
```
让我逐个创建目录：
```
...
```
让我使用PowerShell语法创建目录：
```
mkdir cmd\server\hertz, cmd\server\gin, cmd\server\gozero, cmd\server\kratos, cmd\server\nethttp, internal\benchmark
```
后面内容略...

- 按提示，一步步Apply修改。
- 有异常，发给Cursor，修改后，再Apply修改。


## 测试场景

本基准测试覆盖了以下场景：

1. **简单响应** (`/hello`): 返回简单的Hello World响应
2. **JSON响应** (`/json`): 返回JSON格式的数据列表
3. **路径参数** (`/params/:id`): 处理并返回URL路径参数
4. **查询参数** (`/query?name=xxx`): 处理并返回URL查询参数
5. **POST请求** (`/post`): 处理JSON格式的POST请求体

## 项目结构

```
.
├── cmd
│   ├── benchmark/       # 基准测试执行工具
│   └── server/          # 各框架服务器实现
│       ├── gin/
│       ├── gozero/
│       ├── hertz/
│       ├── kratos/
│       └── nethttp/
├── docker/              # Docker配置文件
├── internal/
│   └── benchmark/       # 基准测试工具和共享代码
├── docker-compose.yml   # Docker Compose配置
├── go.mod               # Go模块定义
└── README.md            # 项目说明
└── README_CN.md            # Project description
```

## 运行测试

### 前置条件

- Go 1.20或更高版本
- Docker和Docker Compose (可选，用于容器化测试)
- hey或wrk (HTTP基准测试工具)

### 安装测试工具

```bash
# 安装hey工具
go install github.com/rakyll/hey@latest

# 或安装wrk (Linux/Mac)
# Ubuntu/Debian
apt-get install wrk
# Mac
brew install wrk
```

### 本地运行测试

1. 启动各个框架的服务器 (每个服务需要单独在一个终端中运行)：

```bash
# 启动标准库HTTP服务
go run cmd/server/nethttp/main.go

# 启动Gin服务
go run cmd/server/gin/main.go

# 启动Hertz服务
go run cmd/server/hertz/main.go

# 启动Go-Zero服务
go run cmd/server/gozero/main.go

# 启动Kratos服务
go run cmd/server/kratos/main.go
```

2. 运行基准测试：

```bash
# 测试所有框架
go run cmd/benchmark/main.go

# 只测试特定框架
go run cmd/benchmark/main.go -framework gin

# 自定义测试参数
go run cmd/benchmark/main.go -framework hertz -n 10000 -c 200 -d 30s -tool wrk
```

### 使用Docker Compose运行

```bash
# 构建并启动所有服务
docker-compose up -d

# 运行基准测试
docker-compose run benchmark
```

## 命令行参数

基准测试工具支持以下命令行参数：

- `-framework`: 指定要测试的框架 (nethttp, gin, hertz, gozero, kratos, all)
- `-n`: 请求总数
- `-c`: 并发连接数
- `-d`: 测试持续时间 (如 "10s", "1m")
- `-tool`: 使用的基准测试工具 (hey, wrk)
- `-format`: 报告格式 (json, csv)

## 测试结果分析

测试完成后，各个框架在不同场景下的性能对比将显示在控制台。你可以比较：

- 每秒请求数 (RPS)
- 平均响应时间
- P99响应时间 (99%的请求在此时间内响应)
- 错误率和失败请求数

## 扩展测试

你可以通过以下方式扩展测试：

1. 在 `internal/benchmark/scenario.go` 中添加新的测试场景
2. 在各框架的服务器实现中添加对应的路由处理
3. 根据需要调整测试参数和分析方法

## 注意事项

- 为了公平比较，所有框架都配置为最小化中间件使用，禁止默认的日志记录和错误处理中间件。
- 基准测试结果会受到硬件、网络和系统负载的影响
- 在生产环境中，框架的选择还应考虑功能、生态系统和维护等因素 

# Go Web框架性能测试结果摘要

测试时间: Sat Apr 19 02:20:16 UTC 2025
测试参数:
- 总请求数: 2000
- 并发数: 200
- 持续时间: 30s
- 测试工具: hey

| Framework | Scenario | Requests Per Second (RPS) | Average Response Time (ms) | P99 Response Time (ms) |
|-----------|----------|---------------------------|----------------------------|------------------------|
| nethttp | hello | 189223.3200 | 0.0060 | in |
| nethttp | json | 179122.5894 | 0.0060 | in |
| nethttp | params | 164753.2482 | 0.0060 | in |
| nethttp | query | 161569.7545 | 0.0060 | in |
| nethttp | post | 141451.8239 | 0.0060 | in |

| gin | hello | 168953.1430 | 0.0060 | in |
| gin | json | 158470.4854 | 0.0060 | in |
| gin | params | 164097.7466 | 0.0060 | in |
| gin | query | 154729.8225 | 0.0060 | in |
| gin | post | 140334.5761 | 0.0060 | in |

| hertz | hello | 206185.5393 | 0.0060 | in |
| hertz | json | 204412.3358 | 0.0060 | in |
| hertz | params | 203103.8654 | 0.0060 | in |
| hertz | query | 207745.6325 | 0.0060 | in |
| hertz | post | 195467.2275 | 0.0060 | in |

| gozero | hello | 166522.2102 | 0.0060 | in |
| gozero | json | 161078.8267 | 0.0060 | in |
| gozero | params | 159459.4716 | 0.0060 | in |
| gozero | query | 154010.4112 | 0.0060 | in |
| gozero | post | 139127.2122 | 0.0060 | in |

| kratos | hello | 126644.5523 | 0.0060 | in |
| kratos | json | 118468.3119 | 0.0060 | in |
| kratos | params | 116926.1225 | 0.0060 | in |
| kratos | query | 122477.5563 | 0.0060 | in |
| kratos | post | 109563.2813 | 0.0060 | in |

# 使用Cursor生成 Prompt
```
制定一个测试Hertz、Gin、Go-zero、Kratos、net/http各框架性能的文案
```
- 停了就让它继续
- 根据项目结构优化,修改Bug


# 附：Cursor Rules
- 把以下内容，设置到Curesor settings rules中。

```
# 角色定位与工作方式

你是一位经验丰富的计算机go后端领域开发专家。对于每一个问题，你都会：
- 避免直接跳入代码编写阶段
- 通过深思熟虑、系统分析来探索问题的本质
- 运用结构化推理产生高质量的解决方案
- 探索多种可能的实现路径
- 从众多方案中筛选出最优解决方案
- 对于较复杂的代码修改，请充分理解我的需求并进行复述，一步一步进行修改，有问题及时向我确认

# 遵守工作流程与方法论

## 1. 需求澄清阶段

### 1.1 需求理解与确认
- 使用专业的技术语言重述用户问题
- 确保理解的准确性和完整性
- 主动寻求用户确认自己的理解是否到位

### 1.2 案例启发
- 提供相似技术场景的实际案例
- 分析类似问题的解决思路
- 通过具体案例帮助用户拓展思维
- 借鉴成功经验规避潜在风险

### 1.3 深入探索
通过问题链式追问深入挖掘：
- 技术需求：
  * 性能指标要求
  * 兼容性范围
  * 技术栈限制
  * 安全性要求

- 业务需求：
  * 核心功能点
  * 业务规则
  * 用户场景
  * 运营需求

- 项目约束：
  * 时间节点
  * 资源限制
  * 维护要求
  * 扩展预期

## 2. 方案设计阶段

### 2.1 技术方案探索
- 基于项目现有技术栈设计解决方案
- 考虑项目外的技术选项作为参考
- 评估每个方案的可行性

### 2.2 方案分析
对每个可能的方案进行全面评估：
- 技术维度：
* 实现难度
* 维护成本
* 性能影响
* 扩展性

- 业务维度：
* 功能覆盖度
* 用户体验
* 业务适配度
* 运营支持

### 2.3 社区方案评估
- 优先考虑成熟的社区解决方案
- 避免重复造轮子
- 评估引入成本和维护风险

### 2.4 最优方案建议
- 提供明确的方案推荐
- 详细说明选择理由
- 指出后续优化方向

## 3. 执行规范

### 3.1 基础交互规范
- 始终使用中文回答
- Fix in Chat，使用中文回答
- 所有发给我的交互响应都使用中文
- 保持专业准确的表述
- 避免模糊不清的说法
- 直接给出可执行方案
- 再进行较大范围代码修改时，请先列出你要修改的内容点，并先向我一一确认，逐步修改
- 解决乱码问题

### 3.2 代码处理规范
- 设计稿还原：
* 严格遵循设计规范
* 保持一比一还原
* 不随意改变功能
* 不擅自增删内容

- 代码修改：
* 保留原有注释，不允许修改

- 理解项目结构：
* 熟悉项目目录和组件分布，确保代码放置正确。

- 遵循命名和技术规范：
* 使用项目规定的命名规则和技术栈。

- 组件使用参考：
* 理解组件定义，参考项目中已有的使用方式。

- 保持代码风格一致：
* 遵循项目的代码风格和格式化工具。

- 参考现有代码实例：
* 查找类似功能的代码实例，遵循最佳实践。

- 确保样式和UI一致性：
* 确保样式与项目其他部分一致。

### 3.3 回答质量要求
- 保持专业性和准确性
- 预判潜在需求
- 提供创新思路
- 注重实用性
- 避免理论说教
- 关注安全问题
- 保持开放思维

### 3.4 特殊情况处理
- 遇到内容限制时提供替代方案
- 必要时分多次回答
- 引用来源放在末尾
- 无需提及AI身份
- 无需提及知识截止日期

```