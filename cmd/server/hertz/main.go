package main

import (
	"context"
	"log"

	"github.com/benchmark/webframework/internal/benchmark"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	// 禁用默认日志
	hlog.SetSilentMode(true)
	// 创建Hertz服务器实例，禁用默认日志
	h := server.New(
		server.WithHostPorts(":8082"),
	)

	// 设置路由
	setupRoutes(h)

	// 启动服务器
	log.Println("Hertz server started on :8082")
	h.Spin()
}

func setupRoutes(h *server.Hertz) {
	// 获取所有测试场景
	testCases := benchmark.GetTestCases()

	// 为每个测试场景注册路由
	for _, tc := range testCases {
		switch tc.Name {
		case "hello":
			h.GET(tc.Endpoint, helloHandler)
		case "json":
			h.GET(tc.Endpoint, jsonHandler)
		case "params":
			h.GET("/params/:id", paramsHandler)
		case "query":
			h.GET(tc.Endpoint, queryHandler)
		case "post":
			h.POST(tc.Endpoint, postHandler)
		}
	}
}

// helloHandler 处理 /hello 请求
func helloHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{
		"message": "Hello, World!",
		"code":    consts.StatusOK,
		"data":    nil,
	})
}

// jsonHandler 处理 /json 请求
func jsonHandler(ctx context.Context, c *app.RequestContext) {
	users := []benchmark.User{
		{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "2", Name: "李四", Email: "lisi@example.com"},
		{ID: "3", Name: "王五", Email: "wangwu@example.com"},
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "Success",
		"code":    consts.StatusOK,
		"data":    users,
	})
}

// paramsHandler 处理 /params/:id 请求
func paramsHandler(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")

	user := benchmark.User{
		ID:    id,
		Name:  "用户" + id,
		Email: "user" + id + "@example.com",
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "Success",
		"code":    consts.StatusOK,
		"data":    user,
	})
}

// queryHandler 处理 /query 请求
func queryHandler(ctx context.Context, c *app.RequestContext) {
	name := c.DefaultQuery("name", "Guest")

	c.JSON(consts.StatusOK, utils.H{
		"message": "Hello, " + name,
		"code":    consts.StatusOK,
		"data":    nil,
	})
}

// postHandler 处理 /post 请求
func postHandler(ctx context.Context, c *app.RequestContext) {
	var req benchmark.PostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"message": "Invalid request body",
			"code":    consts.StatusBadRequest,
			"data":    nil,
		})
		return
	}

	response := benchmark.PostResponse{
		ID:      "post-123",
		Success: true,
		Message: "Created post: " + req.Name,
	}

	c.JSON(consts.StatusCreated, utils.H{
		"message": "Post created",
		"code":    consts.StatusCreated,
		"data":    response,
	})
}
