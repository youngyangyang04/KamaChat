package ssl

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"gochat/pkg/zlog"
	"strconv"
)

func TlsHandler(host string, port int) gin.HandlerFunc {
	sslHost := buildSSLHost(host, port)
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect:   true,
			SSLHost:       sslHost,
			IsDevelopment: false,
			SSLProxyHeaders: map[string]string{
				"X-Forwarded-Proto": "https",
			},
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if c.Writer.Status() == http.StatusMovedPermanently || c.Writer.Status() == http.StatusTemporaryRedirect {
			c.Abort()
			return
		}

		// If there was an error, do not continue.
		if err != nil {
			zlog.Fatal(err.Error())
			return
		}

		c.Next()
	}
}

func buildSSLHost(host string, port int) string {
	trimmedHost := strings.TrimSpace(host)
	if trimmedHost == "" || trimmedHost == "0.0.0.0" || trimmedHost == "::" {
		return ""
	}

	if port == 0 || port == 443 {
		return trimmedHost
	}

	return trimmedHost + ":" + strconv.Itoa(port)
}
