// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package book

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"sample-app/pkg/httperr"
	"sample-app/pkg/log"
	"sample-app/pkg/metrics"
	"sample-app/pkg/tracing"
)

// Server implements Book service
type Server struct {
	hostPort string
	tracer   trace.TracerProvider
	logger   log.Factory
	database *database
}

// NewServer creates a new book.Server
func NewServer(hostPort string, otelExporter string, metricsFactory metrics.Factory, logger log.Factory) *Server {
	return &Server{
		hostPort: hostPort,
		tracer:   tracing.InitOTEL("book", otelExporter, metricsFactory, logger),
		logger:   logger,
		database: newDatabase(
			tracing.InitOTEL("mysql", otelExporter, metricsFactory, logger).Tracer("mysql"),
			logger.With(zap.String("component", "mysql")),
		),
	}
}

// Run starts the Book server
func (s *Server) Run() error {
	mux := s.createServeMux()
	s.logger.Bg().Info("Starting", zap.String("address", "http://"+s.hostPort))
	server := &http.Server{
		Addr:              s.hostPort,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}

func (s *Server) createServeMux() http.Handler {
	mux := tracing.NewServeMux(false, s.tracer, s.logger)
	mux.Handle("/book", http.HandlerFunc(s.book))
	return mux
}

func (s *Server) book(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))
	if err := r.ParseForm(); httperr.HandleError(w, err, http.StatusBadRequest) {
		s.logger.For(ctx).Error("bad request", zap.Error(err))
		return
	}

	book := r.Form.Get("book")
	if book == "" {
		http.Error(w, "Missing required 'book' parameter", http.StatusBadRequest)
		return
	}
	bookID, err := strconv.Atoi(book)
	if err != nil {
		http.Error(w, "Parameter 'book' is not an integer", http.StatusBadRequest)
		return
	}

	response, err := s.database.Get(ctx, bookID)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("request failed", zap.Error(err))
		return
	}

	data, err := json.Marshal(response)
	if httperr.HandleError(w, err, http.StatusInternalServerError) {
		s.logger.For(ctx).Error("cannot marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}