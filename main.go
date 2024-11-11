// main.go
package main

import (
	"context"
	"log"
	"net/http"

	"sample-app/handlers"

	"sample-app/models"
	"sample-app/services"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("book-service"),
		)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

func main() {
	// Initialize tracer
	cleanup := initTracer()
	defer cleanup()

	// Connect to database
	models.ConnectDatabase()

	// Initialize services and handlers
	bookService := services.NewBookService()
	bookHandler := handlers.NewBookHandler(bookService)

	// Set up router with OpenTelemetry instrumentation
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("book-service"))

	// Register routes
	r.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	r.HandleFunc("/books", bookHandler.ListBooks).Methods("GET")
	r.HandleFunc("/books/{id}", bookHandler.GetBook).Methods("GET")
	r.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// Start server
	log.Println("Server is running on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}