package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/platform/httpx"
)

const userIDContextKey = "user_id"

func RequireAuth(signer TokenSigner) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			httpx.Error(c, http.StatusUnauthorized, "missing Authorization header")
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			httpx.Error(c, http.StatusUnauthorized, "invalid Authorization header")
			c.Abort()
			return
		}

		rawToken := strings.TrimSpace(strings.TrimPrefix(header, prefix))
		if rawToken == "" {
			httpx.Error(c, http.StatusUnauthorized, "invalid Authorization header")
			c.Abort()
			return
		}

		claims, err := signer.Verify(rawToken)
		if err != nil || strings.TrimSpace(claims.Subject) == "" {
			httpx.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set(userIDContextKey, claims.Subject)
		c.Next()
	}
}

func UserIDFromContext(c *gin.Context) (string, bool) {
	value, ok := c.Get(userIDContextKey)
	if !ok {
		return "", false
	}

	userID, ok := value.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", false
	}

	return userID, true
}
