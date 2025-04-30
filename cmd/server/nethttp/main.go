package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/benchmark/webframework/internal/benchmark"
)

func main() {
	// 创建路由处理
	mux := http.NewServeMux()
	setupRoutes(mux)

	// 创建服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 启动服务器
	benchmark.StartServer("net/http", ":8080", server)
}

func setupRoutes(mux *http.ServeMux) {
	// 获取所有测试场景
	testCases := benchmark.GetTestCases()

	// 为每个测试场景注册路由
	for _, tc := range testCases {
		switch tc.Name {
		case "hello":
			mux.HandleFunc(tc.Endpoint, helloHandler)
		case "json":
			mux.HandleFunc(tc.Endpoint, jsonHandler)
		case "params":
			// 标准库需要自己处理路径参数
			mux.HandleFunc("/params/", paramsHandler)
		case "query":
			mux.HandleFunc(tc.Endpoint, queryHandler)
		case "post":
			mux.HandleFunc(tc.Endpoint, postHandler)
		}
	}
}

// helloHandler 处理 /hello 请求
func helloHandler(w http.ResponseWriter, r *http.Request) {
	benchmark.CommonResponse(w, http.StatusOK, "Hello, World!", nil)
}

// jsonHandler 处理 /json 请求
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	users := []benchmark.User{
		{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "2", Name: "李四", Email: "lisi@example.com"},
		{ID: "3", Name: "王五", Email: "wangwu@example.com"},
	}
	benchmark.CommonResponse(w, http.StatusOK, "Success", users)
}

// paramsHandler 处理 /params/:id 请求
func paramsHandler(w http.ResponseWriter, r *http.Request) {
	// 手动解析路径参数
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		benchmark.CommonResponse(w, http.StatusBadRequest, "Missing ID parameter", nil)
		return
	}

	id := parts[2] // 获取ID参数
	user := benchmark.User{
		ID:    id,
		Name:  "用户" + id,
		Email: "user" + id + "@example.com",
	}

	benchmark.CommonResponse(w, http.StatusOK, "Success", user)
}

// queryHandler 处理 /query 请求
func queryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	benchmark.CommonResponse(w, http.StatusOK, "Hello, "+name, nil)
}

// postHandler 处理 /post 请求
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		benchmark.CommonResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req benchmark.PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		benchmark.CommonResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	response := benchmark.PostResponse{
		ID:      "post-123",
		Success: true,
		Message: "Created post: " + req.Name,
	}

	benchmark.CommonResponse(w, http.StatusCreated, "Post created", response)
}
