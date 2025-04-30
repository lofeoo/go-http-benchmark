package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/benchmark/webframework/internal/benchmark"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
}

func main() {
	// 禁用访问日志
	logx.Disable()
	// 或者仅显示错误日志
	// logx.SetLevel(logx.ErrorLevel)

	var config Config
	config.RestConf.Host = "0.0.0.0"
	config.RestConf.Port = 8083

	// 可选：在配置中也设置日志级别
	config.RestConf.Log.Mode = "console"
	config.RestConf.Log.Level = "error"

	server := rest.MustNewServer(config.RestConf)
	defer server.Stop()

	// 注册路由
	setupRoutes(server)

	log.Printf("Go-Zero server started on :8083")
	server.Start()
}

func setupRoutes(server *rest.Server) {
	// 获取所有测试场景
	testCases := benchmark.GetTestCases()

	// 为每个测试场景注册路由
	for _, tc := range testCases {
		switch tc.Name {
		case "hello":
			server.AddRoute(rest.Route{
				Method:  http.MethodGet,
				Path:    tc.Endpoint,
				Handler: helloHandler,
			})
		case "json":
			server.AddRoute(rest.Route{
				Method:  http.MethodGet,
				Path:    tc.Endpoint,
				Handler: jsonHandler,
			})
		case "params":
			server.AddRoute(rest.Route{
				Method:  http.MethodGet,
				Path:    "/params/:id",
				Handler: paramsHandler,
			})
		case "query":
			server.AddRoute(rest.Route{
				Method:  http.MethodGet,
				Path:    tc.Endpoint,
				Handler: queryHandler,
			})
		case "post":
			server.AddRoute(rest.Route{
				Method:  http.MethodPost,
				Path:    tc.Endpoint,
				Handler: postHandler,
			})
		}
	}
}

// helloHandler 处理 /hello 请求
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Hello, World!",
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// jsonHandler 处理 /json 请求
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	users := []benchmark.User{
		{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "2", Name: "李四", Email: "lisi@example.com"},
		{ID: "3", Name: "王五", Email: "wangwu@example.com"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    users,
	})
}

// paramsHandler 处理 /params/:id 请求
func paramsHandler(w http.ResponseWriter, r *http.Request) {
	// 从路径中提取参数
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	user := benchmark.User{
		ID:    id,
		Name:  "用户" + id,
		Email: "user" + id + "@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    user,
	})
}

// queryHandler 处理 /query 请求
func queryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Hello, " + name,
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// postHandler 处理 /post 请求
func postHandler(w http.ResponseWriter, r *http.Request) {
	var req benchmark.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request body",
			"code":    http.StatusBadRequest,
			"data":    nil,
		})
		return
	}

	response := benchmark.PostResponse{
		ID:      "post-123",
		Success: true,
		Message: "Created post: " + req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created",
		"code":    http.StatusCreated,
		"data":    response,
	})
}
