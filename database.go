package main

import (
	"sample-app/models"
	"sample-app/pkg/log"
	"sample-app/pkg/tracing"

	"go.opentelemetry.io/otel/trace"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	tracer    trace.Tracer
	logger    log.Factory
	books map[int]*models.Book
	lock      *tracing.Mutex
}

func newDatabase(tracer trace.Tracer, logger log.Factory) *database {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.Book{})
    if err != nil {
        return nil
    }

	return &database{
		tracer: tracer,
		logger: logger,
		lock: &tracing.Mutex{
			SessionBaggageKey: "request",
			LogFactory:        logger,
		},
	}
}