package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RecoveryMiddleware())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCORSMiddleware_AllowsOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	router.ServeHTTP(w, req)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.OPTIONS("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "should not reach")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:8080")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://localhost:8080", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_HeadersSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	router.ServeHTTP(w, req)

	headers := w.Header()
	assert.NotEmpty(t, headers.Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, headers.Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, headers.Get("Access-Control-Allow-Headers"))
}

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggerMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestLoggerMiddleware_StatusCapture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggerMiddleware())

	router.GET("/status404", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/status404", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLoggerMiddleware_PathCapture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggerMiddleware())

	router.GET("/custom/path", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/custom/path", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJSONError_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	JSONError(c, http.StatusBadRequest, "bad request error")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]any
	_ = jsonDecode(w.Body.Bytes(), &response)
	assert.Equal(t, "bad request error", response["error"])
}

func TestJSONError_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	JSONError(c, http.StatusInternalServerError, "internal error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]any
	_ = jsonDecode(w.Body.Bytes(), &response)
	assert.Equal(t, "internal error", response["error"])
}

func TestJSONError_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	JSONError(c, http.StatusConflict, "conflict error")

	assert.Equal(t, http.StatusConflict, w.Code)
	var response map[string]any
	_ = jsonDecode(w.Body.Bytes(), &response)
	assert.Equal(t, "conflict error", response["error"])
}

func TestNotFoundHandler_Method(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Method: "GET"}

	NotFoundHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNotFoundHandler_Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFoundHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]any
	err := jsonDecode(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "resource not found", response["error"])
}

func TestCORSMiddleware_NonPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "created")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_BlocksUnknownOriginPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware())

	router.OPTIONS("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "should not reach")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://evil.example")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestSecurityHeadersMiddleware_SetsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeadersMiddleware())
	router.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
	assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
}

func TestLoggerMiddleware_MethodCapture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(LoggerMiddleware())

	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusCreated, "created")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func jsonDecode(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(v)
}
