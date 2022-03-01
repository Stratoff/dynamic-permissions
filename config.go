package permissions

import (
	"encoding/json"
	"fmt"

	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

const (
	NameSpace = "github.com/stratoff/dynamic-permissions"
)

type PermissionConfig struct {
	PermissionCode string `json:"permission"`
	Backend        string `json:"backend"`
	Method         string `json:"http_method"`
}

// Manticora Function based in IP-FILTER middleware (i think)
func ParseConfig(e config.ExtraConfig, logger logging.Logger) *PermissionConfig {
	logger.Info("<<<<<<< DynamicPermissions: Parse Config for Endpoint >>>>>>>>")
	tmp, ok := e[NameSpace]
	if !ok {
		return nil
	}

	data, err := json.Marshal(tmp)
	if err != nil {
		logger.Error(fmt.Sprintf("Marshal krakend-dynamicPermissions config error: %s", err.Error()))
		return nil
	}
	var permCfg PermissionConfig
	if err := json.Unmarshal(data, &permCfg); err != nil {
		logger.Error(fmt.Sprintf("Unmarshal krakend-dynamicPermissions config error: %s", err.Error()))
		return nil
	}

	// if !strings.HasPrefix(permCfg.Backend, "https://") {
	// 	return &permCfg
	// }

	return &permCfg
}
