package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServeStaticSPA_EmptyPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	err := os.MkdirAll(assetsDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	assetContent := "body { color: red; }"
	err = os.WriteFile(filepath.Join(assetsDir, "main.css"), []byte(assetContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SPA")
}

func TestServeStaticSPA_WithPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	err := os.MkdirAll(assetsDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>App</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "/app", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "App")
}

func TestServeStaticSPA_AssetsRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	err := os.MkdirAll(assetsDir, 0755)
	assert.NoError(t, err)

	assetContent := "console.log('app');"
	err = os.WriteFile(filepath.Join(assetsDir, "app.js"), []byte(assetContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assets/app.js", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "console.log")
}

func TestServeStaticSPA_IndexHtmlRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>Index</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/index.html", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusMovedPermanently {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Index")
}

func TestServeStaticSPA_ClientSideRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA Router</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "SPA Router")
}

func TestServeStaticSPA_UnknownFile404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/unknown.js", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServeStaticSPA_NonGetMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServeStaticSPA_PrefixNormalization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>App</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "/app/", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServeStaticSPA_WithPrefixAssets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	assetsDir := filepath.Join(tmpDir, "assets")
	err := os.MkdirAll(assetsDir, 0755)
	assert.NoError(t, err)

	assetContent := "/* styles */"
	err = os.WriteFile(filepath.Join(assetsDir, "styles.css"), []byte(assetContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "/app", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/assets/styles.css", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "styles")
}

func TestServeStaticSPA_MissingIndexHtml(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServeStaticSPA_DeepClientRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>Deep Route</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dashboard/users/123/settings", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Deep Route")
}

func TestServeStaticSPA_ImageFile404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logo.png", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServeStaticSPA_WithPrefixClientRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>App Route</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "/app", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/profile", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "App Route")
}

func TestServeStaticSPA_WithPrefixUnknownFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	indexContent := "<html><body>SPA</body></html>"
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	assert.NoError(t, err)

	ServeStaticSPA(router, "/app", tmpDir)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/app/missing.js", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
