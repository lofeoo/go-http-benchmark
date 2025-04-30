package main

import (
	"net/http"

	"github.com/benchmark/webframework/internal/benchmark"
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	router := gin.New()

	// 仅使用恢复中间件，不使用日志中间件以避免影响性能
	router.Use(gin.Recovery())

	// 设置路由
	setupRoutes(router)

	// 创建服务器
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// 启动服务器
	benchmark.StartServer("Gin", ":8081", server)
}

func setupRoutes(router *gin.Engine) {
	// 获取所有测试场景
	testCases := benchmark.GetTestCases()

	// 为每个测试场景注册路由
	for _, tc := range testCases {
		switch tc.Name {
		case "hello":
			router.GET(tc.Endpoint, helloHandler)
		case "json":
			router.GET(tc.Endpoint, jsonHandler)
		case "params":
			router.GET("/params/:id", paramsHandler)
		case "query":
			router.GET(tc.Endpoint, queryHandler)
		case "post":
			router.POST(tc.Endpoint, postHandler)
		}
	}
}

// helloHandler 处理 /hello 请求
func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// jsonHandler 处理 /json 请求
func jsonHandler(c *gin.Context) {
	users := []benchmark.User{
		{ID: "1", Name: "张三", Email: "zhangsan@example.com"},
		{ID: "2", Name: "李四", Email: "lisi@example.com"},
		{ID: "3", Name: "王五", Email: "wangwu@example.com"},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    users,
	})
}

// paramsHandler 处理 /params/:id 请求
func paramsHandler(c *gin.Context) {
	id := c.Param("id")

	user := benchmark.User{
		ID:    id,
		Name:  "用户" + id,
		Email: "user" + id + "@example.com",
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"code":    http.StatusOK,
		"data":    user,
	})
}

// queryHandler 处理 /query 请求
func queryHandler(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest")

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, " + name,
		"code":    http.StatusOK,
		"data":    nil,
	})
}

// postHandler 处理 /post 请求
func postHandler(c *gin.Context) {
	var req benchmark.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created",
		"code":    http.StatusCreated,
		"data":    response,
	})
}
