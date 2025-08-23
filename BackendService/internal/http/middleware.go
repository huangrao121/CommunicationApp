package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey struct{}

func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Writer.Header().Set("request_id", rid)

		l := slog.Default().With(
			"method", c.Request.Method,
			"path", c.FullPath(),
			"remote_ip", c.ClientIP(),
		)
		ctx := context.WithValue(c.Request.Context(), ctxKey{}, l)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		l.Info("request completed", "status", c.Writer.Status(), "latency_ms", time.Since(start).Milliseconds(), "bytes", c.Writer.Size())
	}
}

func FromCtx(ctx context.Context) *slog.Logger {
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return slog.Default()
}

// SecureHeaders 添加安全相关的HTTP头
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		// 启用XSS过滤
		c.Header("X-XSS-Protection", "1; mode=block")
		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")
		// 强制使用HTTPS (如果在生产环境)
		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		// 禁止在请求参数中包含敏感信息的引用
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// CSRFProtection CSRF保护中间件
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对非GET请求进行CSRF保护
		if c.Request.Method != "GET" {
			// 获取请求头中的CSRF令牌
			requestToken := c.GetHeader("X-CSRF-Token")
			if requestToken == "" {
				// 尝试从表单获取
				requestToken = c.PostForm("_csrf")
			}

			// 从Cookie中获取CSRF令牌
			csrfCookie, err := c.Cookie("csrf_token")
			if err != nil || csrfCookie == "" || csrfCookie != requestToken {
				c.AbortWithStatusJSON(http.StatusForbidden, "CSRF token validation failed")
				return
			}
		}

		// 对所有请求生成并设置新的CSRF令牌
		token := generateCSRFToken()
		c.SetCookie("csrf_token", token, int(12*time.Hour.Seconds()), "/", "", false, true)
		c.Set("csrf_token", token)
		c.Header("X-CSRF-Token", token)

		c.Next()
	}
}

// CSRFToken 生成CSRF令牌
func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
