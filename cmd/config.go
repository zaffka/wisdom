package cmd

import (
	"github.com/spf13/viper"
)

const (
	serviceName            = "wisdom"
	clientMicroServiceName = "client"
	serverMicroServiceName = "server"

	tcp4Network          = "tcp4"
	serverAddr           = "SERVER_ADDR"
	complexityPoW        = "COMPLEXITY_POW"
	defaultComplexityPoW = 10000000
)

func initConfigDefaults() {
	viper.SetDefault(serverAddr, "0.0.0.0:30333")
	viper.SetDefault(complexityPoW, defaultComplexityPoW)
}
