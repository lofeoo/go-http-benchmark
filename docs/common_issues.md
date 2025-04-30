# Go Web框架性能测试常见问题和解决方案

## 环境设置问题

### 安装依赖报错

**问题**: 安装框架依赖时出现错误，如`go get`命令失败。

**解决方案**:
1. 确保Go版本≥1.20
2. 检查网络连接，必要时设置代理：
   ```bash
   go env -w GOPROXY=https://goproxy.cn,direct
   ```
3. 清除Go模块缓存后重试：
   ```bash
   go clean -modcache
   go mod tidy
   ```

### 端口占用冲突

**问题**: 启动服务器时报端口已被占用错误。

**解决方案**:
1. 使用以下命令查找占用端口的进程：
   ```bash
   # Windows
   netstat -ano | findstr :<PORT>
   
   # Linux/Mac
   lsof -i :<PORT>
   ```
2. 终止对应进程或修改服务器配置使用不同端口

## 测试工具问题

### hey工具安装失败

**问题**: 无法安装hey工具或运行时报错。

**解决方案**:
1. 确保正确安装：
   ```bash
   go install github.com/rakyll/hey@latest
   ```
2. 确保`$GOPATH/bin`在系统PATH中：
   ```bash
   # Bash
   export PATH=$PATH:$(go env GOPATH)/bin
   
   # Windows PowerShell
   $env:Path += ";$(go env GOPATH)\bin"
   ```

### wrk工具兼容性问题

**问题**: wrk工具在Windows上不可用。

**解决方案**:
1. 在Windows上，使用hey工具代替
2. 或使用WSL (Windows Subsystem for Linux) 来运行wrk
3. 或使用Docker容器运行wrk

## 框架特定问题

### Hertz框架依赖问题

**问题**: Hertz依赖解析或导入失败。

**解决方案**:
1. 确保安装完整的Hertz依赖：
   ```bash
   go get github.com/cloudwego/hertz/cmd/hz@latest
   go get github.com/cloudwego/hertz@latest
   ```
2. 可能需要安装thrift相关依赖：
   ```bash
   go get github.com/apache/thrift@latest
   ```

### Go-Zero框架启动问题

**问题**: Go-Zero服务启动时配置错误。

**解决方案**:
1. 确保配置文件格式正确
2. 简化配置，移除不必要的中间件和功能
3. 检查依赖版本兼容性：
   ```bash
   go mod tidy
   go mod verify
   ```

### Kratos框架路由注册问题

**问题**: Kratos的HTTP处理函数注册失败或格式不匹配。

**解决方案**:
1. 确保使用正确的HTTP处理函数签名
2. 使用正确的路由注册方法：
   ```go
   // 正确的注册方式
   srv.Route("/api", func(r *mux.Router) {
       r.HandleFunc("/hello", HelloHandler)
   })
   ```

## 性能测试问题

### 测试结果不稳定

**问题**: 多次运行相同测试，结果波动很大。

**解决方案**:
1. 减少系统后台任务和其他应用程序
2. 增加预热时间，丢弃初始测试结果
3. 多次运行测试取平均值
4. 确保CPU频率稳定 (禁用动态频率调整)

### 内存泄漏

**问题**: 长时间运行测试后，服务内存使用持续增长。

**解决方案**:
1. 使用pprof进行内存分析：
   ```go
   import _ "net/http/pprof"
   ```
2. 检查是否正确关闭资源（文件、连接等）
3. 为每个框架实现添加适当的超时和限制

## 结果分析问题

### 数据解析错误

**问题**: 无法从测试工具输出中解析性能数据。

**解决方案**:
1. 检查测试工具版本，确保输出格式符合预期
2. 修改结果解析脚本以适应当前输出格式
3. 使用`-o`参数将结果输出为JSON或CSV格式便于解析

### 图表生成失败

**问题**: 测试结果可视化或图表生成失败。

**解决方案**:
1. 检查数据格式是否符合图表工具要求
2. 尝试使用不同的可视化工具（如Gnuplot、Python matplotlib或在线工具）
3. 确保安装了正确的依赖库

## Docker相关问题

### 容器间通信问题

**问题**: 使用Docker Compose时，基准测试容器无法连接到服务容器。

**解决方案**:
1. 确保所有容器在同一网络中
2. 使用服务名而非localhost作为主机名
3. 等待服务完全启动后再开始测试（可增加依赖检查）

### 资源限制问题

**问题**: Docker容器性能受限。

**解决方案**:
1. 检查并调整容器资源限制：
   ```yaml
   # docker-compose.yml
   services:
     myservice:
       deploy:
         resources:
           limits:
             cpus: '2'
             memory: 2G
   ```
2. 在主机上直接运行测试以获得更准确的结果

## 报告生成问题

### Markdown表格格式错误

**问题**: 生成的Markdown报告表格格式错乱。

**解决方案**:
1. 确保所有行具有相同的列数
2. 检查特殊字符是否被正确转义
3. 使用Markdown预览工具验证语法 