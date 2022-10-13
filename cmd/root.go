package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zaffka/wisdom/pkg/zaplog"
	"go.uber.org/zap"
)

var (
	ver = "dev"

	logger *zap.Logger
)

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: "An RPC server and a client to get a random wisdom quote",
}

func Execute() {
	cobra.OnInitialize(viper.AutomaticEnv, initConfigDefaults)

	hostname, err := os.Hostname()
	if err != nil {
		print(err)
		os.Exit(1)
	}

	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer rootCancel()

	logger = zaplog.New(os.Stderr, zaplog.Opts{
		Host:             hostname,
		Service:          serviceName,
		Version:          ver,
		Debug:            true,
		IsDevEnvironment: true,
	})

	defer func() {
		logger.Info("the app is finished")

		if err := logger.Sync(); err != nil {
			logger.Error("failed to flush logger data", zap.Error(err))
		}
	}()

	if err := rootCmd.ExecuteContext(rootCtx); err != nil {
		logger.Error("failed to execute a root cmd", zap.Error(err))
	}
}
