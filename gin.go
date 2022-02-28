package permissions

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
	ginlura "github.com/luraproject/lura/router/gin"
)

func HandlerFactory(hf ginlura.HandlerFactory, logger logging.Logger) ginlura.HandlerFactory {
	return MiddlewareHandlerFunc(hf, logger)
}

func NewHasPermission(permConfig *PermissionConfig, authorization string) bool {

	token := strings.TrimPrefix(authorization, "Bearer ")
	if token == authorization {
		return false
	}

	role, err := GetRoleFromPayload(token)
	if err != nil {
		return false
	}

	hasPermission := CheckRoleWithBackend(
		role,
		permConfig.PermissionCode,
		permConfig.Backend,
		permConfig.Method,
	)
	return hasPermission
}

func MiddlewareHandlerFunc(hf ginlura.HandlerFactory, logger logging.Logger) ginlura.HandlerFactory {
	return func(cfg *config.EndpointConfig, prxy proxy.Proxy) gin.HandlerFunc {
		permConfig := ParseConfig(cfg.ExtraConfig, logger)
		logger.Info("<<<<<<<DynamicPermissions enabled for the endpoint >>>>>>>>", cfg.Endpoint)
		return func(c *gin.Context) {
			auth := c.Request.Header.Get("Authorization")
			hasPermission := NewHasPermission(permConfig, auth)
			if hasPermission {
				logger.Info("<<<<<<< DynamicPermissions: The user has the permission >>>>>>>>")
				proxyReq := ginlura.NewRequest(cfg.HeadersToPass)(c, cfg.QueryString)
				ctx, cancel := context.WithTimeout(c, cfg.Timeout)
				defer cancel()
				response, err := prxy(ctx, proxyReq)

				if err != nil {
					logger.Error("<<<<<<DynamicPermissions response error:", err.Error())
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}

				if response == nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				c.JSON(response.Metadata.StatusCode, response.Data)
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}
