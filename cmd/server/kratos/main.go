package main

import (
	"log"
	"net/http"

	"github.com/benchmark/webframework/internal/benchmark"
	"github.com/go-kratos/kratos/v2"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	// 创建HTTP服务
	httpSrv := khttp.NewServer(
		khttp.Address(":8084"),
		// 不添加中间件
	)

	// 注册路由
	setupRoutes(httpSrv)

	app := kratos.New(
		kratos.Name("kratos"),
		kratos.Server(httpSrv),
	)

	// 启动服务
	if err := app.Run(); err != nil {
		log.Fatalf("kratos server failed to start: %v", err)
	}
}

func setupRoutes(srv *khttp.Server) {
	// 获取所有测试场景
	testCases := benchmark.GetTestCases()

	router := srv.Route("/")

	// 为每个测试场景注册路由
	for _, tc := range testCases {
		switch tc.Name {
		case "hello":
			router.GET(tc.Endpoint, helloHandler)
		case "json":
			router.GET(tc.Endpoint, jsonHandler)
		case "params":
			router.GET("/params/{id}", paramsHandler)
		case "query":
			router.GET(tc.Endpoint, queryHandler)
		case "post":
			router.POST(tc.Endpoint, postHandler)
		}
	}
}

// helloHandler 处理 /hello 请求
func helloHandler(ctx khttp.Context) error {
	return ctx.Result(http.StatusOK, map[string]interface{}{
		"message": "Hello, World!",
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// jsonHandler 处理 /json 请求
func jsonHandler(ctx khttp.Context) error {
	users := []benchmark.User{
		{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "2", Name: "李四", Email: "lisi@example.com"},
		{ID: "3", Name: "王五", Email: "wangwu@example.com"},
	}

	return ctx.Result(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    users,
	})
}

// paramsHandler 处理 /params/:id 请求
func paramsHandler(ctx khttp.Context) error {
	id := ctx.Vars().Get("id")

	user := benchmark.User{
		ID:    id,
		Name:  "用户" + id,
		Email: "user" + id + "@example.com",
	}

	return ctx.Result(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    user,
	})
}

// queryHandler 处理 /query 请求
func queryHandler(ctx khttp.Context) error {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	return ctx.Result(http.StatusOK, map[string]interface{}{
		"message": "Hello, " + name,
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// postHandler 处理 /post 请求
func postHandler(ctx khttp.Context) error {
	var req benchmark.PostRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.Result(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
			"code":    http.StatusBadRequest,
			"data":    nil,
		})
	}

	response := benchmark.PostResponse{
		ID:      "post-123",
		Success: true,
		Message: "Created post: " + req.Name,
	}

	return ctx.Result(http.StatusCreated, map[string]interface{}{
		"message": "Post created",
		"code":    http.StatusCreated,
		"data":    response,
	})
}
