package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benchmark/webframework/internal/benchmark"
)

type Config struct {
	Framework    string
	BaseURL      string
	Connections  int
	Threads      int
	Duration     time.Duration
	Tool         string
	ReportFormat string
}

var frameworks = map[string]string{
	"nethttp": "http://localhost:8080",
	"gin":     "http://localhost:8081",
	"hertz":   "http://localhost:8082",
	"gozero":  "http://localhost:8083",
	"kratos":  "http://localhost:8084",
}

func main() {
	config := parseFlags()

	// 获取测试场景
	testCases := benchmark.GetTestCases()

	// 在所有测试场景上循环
	for _, tc := range testCases {
		fmt.Printf("\n============ 测试场景: %s (%s) ============\n", tc.Name, tc.Description)

		// 构建完整URL
		url := fmt.Sprintf("%s%s", config.BaseURL, tc.Endpoint)
		if tc.Name == "params" {
			// 对于需要参数的API，添加示例参数
			url = fmt.Sprintf("%s/params/123", config.BaseURL)
		} else if tc.Name == "query" {
			// 对于需要查询参数的API，添加示例参数
			url = fmt.Sprintf("%s%s?name=测试用户", config.BaseURL, tc.Endpoint)
		}

		fmt.Printf("测试URL: %s\n", url)

		// 创建基准测试配置
		benchConfig := benchmark.BenchmarkConfig{
			URL:         url,
			Connections: config.Connections,
			Duration:    config.Duration,
			Threads:     config.Threads,
			Timeout:     5 * time.Second,
		}

		// 为POST请求添加请求体
		if tc.Name == "post" {
			benchConfig.Method = "POST"
			benchConfig.Body = `{"name":"测试内容","content":"这是一个测试POST请求的内容"}`
		}

		// 根据指定的工具执行基准测试
		var err error
		if config.Tool == "hey" {
			err = benchmark.RunHey(benchConfig)
		} else if config.Tool == "wrk" {
			err = benchmark.RunWrk(benchConfig)
		}

		if err != nil {
			log.Printf("执行 %s 基准测试失败: %v", tc.Name, err)
		}

		// 测试之间添加一点延迟，让系统冷却
		time.Sleep(2 * time.Second)
	}
}

func parseFlags() *Config {
	framework := flag.String("framework", "all", "要测试的框架 (nethttp, gin, hertz, gozero, kratos, all)")
	connections := flag.Int("n", 1000, "请求总数")
	concurrency := flag.Int("c", 100, "并发连接数")
	duration := flag.Duration("d", 10*time.Second, "测试持续时间")
	tool := flag.String("tool", "hey", "基准测试工具 (hey, wrk)")
	format := flag.String("format", "json", "报告格式 (json, csv)")

	flag.Parse()

	config := &Config{
		Framework:    *framework,
		Connections:  *connections,
		Threads:      *concurrency,
		Duration:     *duration,
		Tool:         *tool,
		ReportFormat: *format,
	}

	// 根据指定的框架设置基础URL
	if *framework != "all" {
		if baseURL, ok := frameworks[*framework]; ok {
			config.BaseURL = baseURL
		} else {
			fmt.Printf("未知的框架: %s\n", *framework)
			printFrameworkOptions()
			os.Exit(1)
		}
	} else {
		// 如果测试所有框架，默认从net/http开始
		config.BaseURL = frameworks["nethttp"]
		config.Framework = "nethttp"
	}

	return config
}

func printFrameworkOptions() {
	fmt.Println("可用的框架:")
	for k := range frameworks {
		fmt.Printf("  - %s\n", k)
	}
}
