package permissions

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

func Register(cfg *config.ServiceConfig, logger logging.Logger, engine *gin.Engine) {
	permConfig := ParseConfig(cfg.ExtraConfig, logger)
	if permConfig == nil {
		return
	}
	logger.Info("Dynamic Permissions Enabled for the endpoint :D")
	engine.Use(Middle(permConfig, logger))
}

func Middle(permConfig *PermissionConfig, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			logger.Error("Dynamic-Permissions: Authorization Header Field empty.")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			logger.Error("Dynamic-Permissions: TrimPrefix process does not work well.")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		role, err := GetRoleFromPayload(token)
		if err != nil {
			logger.Error("Dynamic-Permissions: Can't get the role from the token payload.")
			logger.Error(err.Error())
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		hasPermission := CheckRoleWithBackend(
			role,
			permConfig.PermissionCode,
			permConfig.Backend,
			permConfig.Method,
		)

		if !hasPermission {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Next()
	}
}
