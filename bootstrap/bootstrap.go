package bootstrap

import (
	btsConfig "gobase/config"

	"github.com/boyane126/go-common/config"
)

func Bootstrap() {
	btsConfig.Initialize()
	config.InitConfig("")

	SetupLogger()
}
