package ssl

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTlsHandlerRedirectsPlainHTTP(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TlsHandler("chat.example.com", 443))
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "http://service.internal/ping", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("expected redirect status, got %d", recorder.Code)
	}

	location := recorder.Header().Get("Location")
	if !strings.HasPrefix(location, "https://chat.example.com") {
		t.Fatalf("expected https redirect to configured host, got %q", location)
	}
}

func TestTlsHandlerHonorsForwardedProto(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TlsHandler("chat.example.com", 443))
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "http://service.internal/ping", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected request to pass through when proxy marks https, got %d", recorder.Code)
	}
}

func TestBuildSSLHostSkipsWildcardBindAddress(t *testing.T) {
	t.Parallel()

	if host := buildSSLHost("0.0.0.0", 8000); host != "" {
		t.Fatalf("expected empty ssl host for wildcard bind, got %q", host)
	}
}
