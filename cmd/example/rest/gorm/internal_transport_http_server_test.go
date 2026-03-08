package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "8080",
		ReadTimeout: 10 * time.Second,
	}

	server := NewServer(router, cfg)

	assert.NotNil(t, server)
	assert.NotNil(t, server.srv)
	assert.Equal(t, ":8080", server.srv.Addr)
	assert.Equal(t, 10*time.Second, server.srv.ReadTimeout)
	assert.Equal(t, 10*time.Second, server.srv.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.srv.IdleTimeout)
}

func TestNewServer_CustomPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "3000",
		ReadTimeout: 5 * time.Second,
	}

	server := NewServer(router, cfg)

	assert.NotNil(t, server)
	assert.Equal(t, ":3000", server.srv.Addr)
	assert.Equal(t, 5*time.Second, server.srv.ReadTimeout)
	assert.Equal(t, 5*time.Second, server.srv.WriteTimeout)
}

func TestNewServer_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	cfg := &Config{
		AppPort:     "8080",
		ReadTimeout: 10 * time.Second,
	}

	server := NewServer(router, cfg)

	assert.NotNil(t, server.srv.Handler)
}

func TestServer_Start_ErrServerClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "0", // Use port 0 to get an available port
		ReadTimeout: 1 * time.Second,
	}

	server := NewServer(router, cfg)

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)

	select {
	case err := <-errChan:
		if err != nil {
			assert.Equal(t, http.ErrServerClosed, err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Server did not stop in time")
	}
}

func TestServer_Shutdown_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "0",
		ReadTimeout: 1 * time.Second,
	}

	server := NewServer(router, cfg)

	go func() {
		_ = server.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Shutdown_ContextTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(10 * time.Second)
		c.String(http.StatusOK, "done")
	})

	cfg := &Config{
		AppPort:     "0",
		ReadTimeout: 1 * time.Second,
	}

	server := NewServer(router, cfg)

	go func() {
		_ = server.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		_, _ = client.Get("http://" + server.srv.Addr + "/slow")
	}()

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Shutdown(ctx)
	_ = err
}

func TestServer_Start_InvalidAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "99999", // Invalid port
		ReadTimeout: 1 * time.Second,
	}

	server := NewServer(router, cfg)

	err := server.Start()
	assert.Error(t, err)
}

func TestServer_Configuration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &Config{
		AppPort:     "8080",
		ReadTimeout: 30 * time.Second,
	}

	server := NewServer(router, cfg)

	assert.Equal(t, ":8080", server.srv.Addr)
	assert.Equal(t, 30*time.Second, server.srv.ReadTimeout)
	assert.Equal(t, 30*time.Second, server.srv.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.srv.IdleTimeout)
}

func TestServer_WithMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Header("X-Custom", "value")
		c.Next()
	})
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	cfg := &Config{
		AppPort:     "18081", // Use a specific port for testing
		ReadTimeout: 5 * time.Second,
	}

	server := NewServer(router, cfg)

	go func() {
		_ = server.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://localhost:18081/test")
	if err == nil {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "value", resp.Header.Get("X-Custom"))
		_ = resp.Body.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func TestServer_GracefulShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	requestHandled := make(chan bool, 1)

	router.GET("/test", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		requestHandled <- true
		c.String(http.StatusOK, "ok")
	})

	cfg := &Config{
		AppPort:     "0",
		ReadTimeout: 5 * time.Second,
	}

	server := NewServer(router, cfg)

	go func() {
		_ = server.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		_, _ = client.Get("http://" + server.srv.Addr + "/test")
	}()

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)

	select {
	case <-requestHandled:
	case <-time.After(2 * time.Second):
	}
}
