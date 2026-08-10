package middleware

import (
	"net/http"
	"strings"

	"github.com/clementscontractors/equipment/internal/authutil"
	"github.com/clementscontractors/equipment/internal/dto"
	"github.com/clementscontractors/equipment/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID = "userID"
	ContextRole   = "role"
	ContextEmail  = "email"
)

// CORS sets CORS headers for the SPA origin.
func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Auth validates Bearer JWT access tokens.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			abortErr(c, http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token")
			return
		}
		claims, err := authutil.ParseAccessToken(secret, strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			abortErr(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Set(ContextEmail, claims.Email)
		c.Next()
	}
}

// RequireRoles allows only listed roles.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(ContextRole)
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			abortErr(c, http.StatusForbidden, "FORBIDDEN", "insufficient role")
			return
		}
		c.Next()
	}
}

// RequireWrite allows operator+ for assignment writes; manager+ for other mutations.
func RequireManagerPlus() gin.HandlerFunc {
	return RequireRoles(models.RoleAdmin, models.RoleManager)
}

// RequireOperatorPlus allows operator, manager, admin.
func RequireOperatorPlus() gin.HandlerFunc {
	return RequireRoles(models.RoleAdmin, models.RoleManager, models.RoleOperator)
}

// RequireAdmin allows admin only.
func RequireAdmin() gin.HandlerFunc {
	return RequireRoles(models.RoleAdmin)
}

func abortErr(c *gin.Context, status int, code, msg string) {
	c.AbortWithStatusJSON(status, dto.ErrorResponse{Error: dto.APIError{Code: code, Message: msg}})
}
