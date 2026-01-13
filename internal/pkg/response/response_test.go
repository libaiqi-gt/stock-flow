package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestSuccess(t *testing.T) {
	r := setupRouter()
	r.GET("/success", func(c *gin.Context) {
		data := TestData{ID: 1, Name: "Test"}
		Success(c, data)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/success", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "操作成功", resp.Msg)
	
	// Data is interface{}, decoded as map[string]interface{} by json
	dataMap, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1.0, dataMap["id"]) // JSON numbers are float64
	assert.Equal(t, "Test", dataMap["name"])
	assert.NotEmpty(t, resp.Timestamp)
}

func TestError(t *testing.T) {
	r := setupRouter()
	r.GET("/error", func(c *gin.Context) {
		Error(c, 400, "Bad Request")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "Bad Request", resp.Msg)
	assert.Nil(t, resp.Data)
}

func TestErrorWithDetail(t *testing.T) {
	r := setupRouter()
	r.GET("/error_detail", func(c *gin.Context) {
		ErrorWithDetail(c, 500, "Server Error", "Database connection failed")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error_detail", nil)
	r.ServeHTTP(w, req)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, "Server Error", resp.Msg)
	assert.Equal(t, "Database connection failed", resp.ErrMsg)
}

func TestI18n(t *testing.T) {
	r := setupRouter()
	r.GET("/i18n", func(c *gin.Context) {
		Error(c, CodeUnauthorized, "")
	})

	// Test English
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/i18n", nil)
	req1.Header.Set("Accept-Language", "en-US")
	r.ServeHTTP(w1, req1)

	var resp1 Response
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.Equal(t, "Unauthorized, please login", resp1.Msg)

	// Test Chinese (Default)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/i18n", nil)
	r.ServeHTTP(w2, req2)

	var resp2 Response
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.Equal(t, "未授权，请登录", resp2.Msg)
}
