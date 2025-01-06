package handler

import (
	"github.com/labstack/echo/v4"
	"gitlab.snapp.ir/snappcloud/config-server/internal/engine"
	"net/http"
)

type Config struct {
	engine         engine.Engine
	internalConfig map[string]map[string]interface{}
}

func NewConfigHandler(engine engine.Engine, internalConfig map[string]map[string]interface{}) *Config {
	return &Config{engine, internalConfig}
}

func (h *Config) Get(c echo.Context) error {
	endpoint := c.Param("endpoint")

	config := h.internalConfig[endpoint]

	if configEngine, ok := h.engine.Engines[endpoint]; ok {
		engineConfig, err := configEngine.GetConfig()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		mergeMaps(config, engineConfig)
	}

	return c.JSON(http.StatusOK, config)
}

func mergeMaps(cfg, engineCfg map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for k, v := range engineCfg {
		merged[k] = v
	}

	for k, v := range cfg {
		merged[k] = v
	}

	return merged
}
