// Copyright (c) 2023 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	otelExporter string // otlp, stdout
	verbose      bool

	fixDBConnDelay         time.Duration
	fixDBConnDisableMutex  bool
	fixRouteWorkerPoolSize int

	bookPort int

	basepath string
	jaegerUI string
)

// used by root command
func addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&otelExporter, "otel-exporter", "x", "otlp", "OpenTelemetry exporter (otlp|stdout)")

	cmd.PersistentFlags().DurationVarP(&fixDBConnDelay, "fix-db-query-delay", "D", 300*time.Millisecond, "Average latency of MySQL DB query")
	cmd.PersistentFlags().BoolVarP(&fixDBConnDisableMutex, "fix-disable-db-conn-mutex", "M", false, "Disables the mutex guarding db connection")
	cmd.PersistentFlags().IntVarP(&fixRouteWorkerPoolSize, "fix-route-worker-pool-size", "W", 3, "Default worker pool size")

	// Add flags to choose ports for services
	cmd.PersistentFlags().IntVarP(&bookPort, "book-service-port", "c", 8090, "Port for book service")

	// Flag for serving frontend at custom basepath url
	cmd.PersistentFlags().StringVarP(&basepath, "basepath", "b", "", `Basepath for frontend service(default "/")`)
	cmd.PersistentFlags().StringVarP(&jaegerUI, "jaeger-ui", "j", "http://localhost:16686", "Address of Jaeger UI to create [find trace] links")

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enables debug logging")
}