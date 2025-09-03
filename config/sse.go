package config

import "github.com/boyane126/go-common/config"

func init() {
	config.Add("sse", func() map[string]interface{} {
		return map[string]interface{}{
			"host":    config.Env("SSE_HOST", "0.0.0.0"),
			"port":    config.Env("SSE_PORT", "8085"),
			"channel": config.Env("SSE_REDIS_CHANNEL", "notifications"),
		}
	})
}
