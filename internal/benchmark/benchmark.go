package benchmark

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

// BenchmarkConfig 配置基准测试参数
type BenchmarkConfig struct {
	URL               string
	Connections       int
	Duration          time.Duration
	Threads           int
	Timeout           time.Duration
	Method            string
	Headers           map[string]string
	Body              string
	RateLimit         int
	ReportFormat      string
	ReportDestination string
}

// DefaultConfig 返回默认配置
func DefaultConfig() BenchmarkConfig {
	return BenchmarkConfig{
		Connections:  100,
		Duration:     10 * time.Second,
		Threads:      4,
		Timeout:      5 * time.Second,
		Method:       "GET",
		Headers:      map[string]string{"Content-Type": "application/json"},
		ReportFormat: "json",
	}
}

// RunHey 使用hey工具运行HTTP基准测试
func RunHey(config BenchmarkConfig) error {
	args := []string{
		"-n", fmt.Sprintf("%d", config.Connections),
		"-c", fmt.Sprintf("%d", config.Threads),
		"-t", config.Timeout.String(),
		"-m", config.Method,
	}

	for k, v := range config.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}

	if config.Body != "" {
		args = append(args, "-d", config.Body)
	}

	if config.ReportDestination != "" {
		args = append(args, "-o", config.ReportDestination)
	}

	args = append(args, config.URL)

	cmd := exec.Command("hey", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Running benchmark: hey %s", strings.Join(args, " "))
	return cmd.Run()
}

// RunWrk 使用wrk工具运行HTTP基准测试
func RunWrk(config BenchmarkConfig) error {
	args := []string{
		"-c", fmt.Sprintf("%d", config.Threads),
		"-t", fmt.Sprintf("%d", config.Threads),
		"-d", config.Duration.String(),
	}

	// wrk不支持直接设置头信息，需要通过脚本
	if len(config.Headers) > 0 || config.Method != "GET" || config.Body != "" {
		// 这里简化处理，实际使用时需要生成Lua脚本
		log.Println("Warning: wrk需要通过Lua脚本设置header、请求方法和body")
	}

	args = append(args, config.URL)

	cmd := exec.Command("wrk", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Running benchmark: wrk %s", strings.Join(args, " "))
	return cmd.Run()
}

// StartServer 启动服务并等待信号退出
func StartServer(name string, addr string, server *http.Server) {
	// 捕获信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// 启动服务器
	go func() {
		log.Printf("%s server started on %s", name, addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s server failed to start: %v", name, err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Printf("Shutting down %s server...", name)
}
