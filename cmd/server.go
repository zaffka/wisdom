package cmd

import (
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zaffka/wisdom/internal/serve"
	"github.com/zaffka/wisdom/pkg/pow"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(server)
}

var server = &cobra.Command{
	Use:   serverMicroServiceName,
	Short: "Server to serve wisdom quotes via TCP4",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		log := logger.With(
			zap.String("micro_service", serverMicroServiceName),
			zap.String("address", viper.GetString(serverAddr)))

		log.Info("starting the server")

		listener, err := net.Listen("tcp", viper.GetString(serverAddr))
		if err != nil {
			log.Error("failed to init net listener", zap.Error(err))

			return
		}
		defer func() {
			if err := listener.Close(); err != nil {
				log.Error("failed to close a listener", zap.Error(err))
			}
		}()

		server := serve.NewServer(
			serve.WithListener(listener),
			serve.WithLogger(log),
			serve.WithInitialPoWComplexity(viper.GetInt64(complexityPoW)),
			serve.WithPoWBlockFunc(pow.NewBlock),
		)

		go server.Run(ctx)
		<-ctx.Done()
	},
}
