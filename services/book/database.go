// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package book

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"sample-app/pkg/log"
	"sample-app/pkg/otelsemconv"
	"sample-app/pkg/tracing"

	"sample-app/pkg/delay"
	"sample-app/services/config"
)

// database simulates Book repository implemented on top of an SQL database
type database struct {
	tracer    trace.Tracer
	logger    log.Factory
	books map[int]*Book
	lock      *tracing.Mutex
}

func newDatabase(tracer trace.Tracer, logger log.Factory) *database {
	return &database{
		tracer: tracer,
		logger: logger,
		lock: &tracing.Mutex{
			SessionBaggageKey: "request",
			LogFactory:        logger,
		},
		books: map[int]*Book{
			1: {
				ID:       1,
				Title:     "Book Title 1",
				Author:   "Author 1",
			},
		},
	}
}

func (d *database) Get(ctx context.Context, bookID int) (*Book, error) {
	d.logger.For(ctx).Info("Loading book", zap.Int("book_id", bookID))

	ctx, span := d.tracer.Start(ctx, "SQL SELECT", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		otelsemconv.PeerServiceKey.String("mysql"),
		attribute.
			Key("sql.query").
			String(fmt.Sprintf("SELECT * FROM book WHERE book_id=%d", bookID)),
	)
	defer span.End()

	if !config.MySQLMutexDisabled {
		// simulate misconfigured connection pool that only gives one connection at a time
		d.lock.Lock(ctx)
		defer d.lock.Unlock()
	}

	// simulate RPC delay
	delay.Sleep(config.MySQLGetDelay, config.MySQLGetDelayStdDev)

	if book, ok := d.books[bookID]; ok {
		return book, nil
	}
	return nil, errors.New("invalid book ID")
}