package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRolesRaw, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
				"error": "forbidden: no roles assigned",
			})
			return
		}

		userRoles, ok := userRolesRaw.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
				"error": "forbidden: invalid roles format",
			})
			return
		}

		hasRole := false
		for _, ur := range userRoles {
			if slices.Contains(allowedRoles, ur) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{
				"error": "forbidden: insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
