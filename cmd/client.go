package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zaffka/wisdom/internal/call"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(client)
}

var client = &cobra.Command{
	Use:   clientMicroServiceName,
	Short: "Client to get a wisdom quote from the server",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.With(zap.String("micro_service", clientMicroServiceName))

		caller := call.NewCaller(
			call.WithServerAddr(viper.GetString(serverAddr)),
			call.WithLogger(log),
			call.WithProtocol(tcp4Network),
		)

		if err := caller.Run(cmd.Context()); err != nil {
			log.Error("failed to get a wisdom quote", zap.Error(err))
		}
	},
}
