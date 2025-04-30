package benchmark

import (
	"encoding/json"
	"net/http"
)

// TestCase 定义了一个测试场景
type TestCase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
}

// ResponseData 简单的响应数据结构
type ResponseData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

// CommonResponse 返回通用的JSON响应
func CommonResponse(w http.ResponseWriter, code int, message string, data any) {
	resp := ResponseData{
		Message: message,
		Code:    code,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// GetTestCases 返回所有测试场景
func GetTestCases() []TestCase {
	return []TestCase{
		{
			Name:        "hello",
			Description: "简单的Hello World API",
			Endpoint:    "/hello",
		},
		{
			Name:        "json",
			Description: "返回JSON数据的API",
			Endpoint:    "/json",
		},
		{
			Name:        "params",
			Description: "带路径参数的API",
			Endpoint:    "/params/:id",
		},
		{
			Name:        "query",
			Description: "带查询参数的API",
			Endpoint:    "/query",
		},
		{
			Name:        "post",
			Description: "处理POST请求的API",
			Endpoint:    "/post",
		},
	}
}

// User 用户数据结构
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// PostRequest POST请求结构
type PostRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// PostResponse POST响应结构
type PostResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
