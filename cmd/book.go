package cmd

import (
	"net"
	"sample-app/pkg/log"
	"strconv"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"sample-app/services/book"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Starts Book service",
	Long:  `Starts Book service.`,
	RunE: func(_ *cobra.Command, _ /* args */ []string) error {
		zapLogger := logger.With(zap.String("service", "book"))
		logger := log.NewFactory(zapLogger)
		server := book.NewServer(
			net.JoinHostPort("0.0.0.0", strconv.Itoa(bookPort)),
			otelExporter,
			metricsFactory,
			logger,
		)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(bookCmd)
}