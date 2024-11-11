package cmd

import (
	"os"

	"sample-app/internal/prometheus"

	"github.com/spf13/cobra"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sample-app/internal/jaegerclientenv2otel"
	"sample-app/pkg/metrics"
	"sample-app/services/config"
)

var (
	logger         *zap.Logger
	metricsFactory metrics.Factory
)

var RootCmd = &cobra.Command{
	Use:   "sample-app",
	Short: "sample-app. - A tracing demo application",
	Long:  `sample-app. - A tracing demo application.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logger.Fatal("We bowled a googly", zap.Error(err))
		os.Exit(-1)
	}
}

func init() {
	addFlags(RootCmd)
	cobra.OnInitialize(onInitialize)
}

func onInitialize() {
	zapOptions := []zap.Option{
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	}
	if !verbose {
		zapOptions = append(zapOptions,
			zap.IncreaseLevel(zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l != zapcore.DebugLevel })),
		)
	}
	logger, _ = zap.NewDevelopment(zapOptions...)

	jaegerclientenv2otel.MapJaegerToOtelEnvVars(logger)

	metricsFactory = prometheus.New().Namespace(metrics.NSOptions{Name: "sample-app", Tags: nil})

	if config.MySQLGetDelay != fixDBConnDelay {
		logger.Info("fix: overriding MySQL query delay", zap.Duration("old", config.MySQLGetDelay), zap.Duration("new", fixDBConnDelay))
		config.MySQLGetDelay = fixDBConnDelay
	}
	if fixDBConnDisableMutex {
		logger.Info("fix: disabling db connection mutex")
		config.MySQLMutexDisabled = true
	}
	if config.RouteWorkerPoolSize != fixRouteWorkerPoolSize {
		logger.Info("fix: overriding route worker pool size", zap.Int("old", config.RouteWorkerPoolSize), zap.Int("new", fixRouteWorkerPoolSize))
		config.RouteWorkerPoolSize = fixRouteWorkerPoolSize
	}

	if bookPort != 8090 {
		logger.Info("changing book service port", zap.Int("old", 8090), zap.Int("new", bookPort))
	}
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running command", zap.Error(err))
	}
	return err
}