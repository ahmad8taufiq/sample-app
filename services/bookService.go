package services

import (
	"context"
	"errors"
	"fmt"

	"sample-app/models"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type BookService struct {
	tracer trace.Tracer
}

func NewBookService() *BookService {
	return &BookService{
		tracer: otel.Tracer("book-service"),
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book) error {
	ctx, span := s.tracer.Start(ctx, "CreateBook")
	defer span.End()

	span.SetAttributes(
		attribute.String("book.title", book.Title),
		attribute.String("book.author", book.Author),
	)

	result := models.DB.WithContext(ctx).Create(book)
	if result.Error != nil {
		span.RecordError(result.Error)
		return fmt.Errorf("failed to create book: %v", result.Error)
	}

	return nil
}

func (s *BookService) GetBook(ctx context.Context, id uint) (*models.Book, error) {
	ctx, span := s.tracer.Start(ctx, "GetBook")
	defer span.End()

	span.SetAttributes(attribute.Int64("book.id", int64(id)))

	var book models.Book
	result := models.DB.WithContext(ctx).First(&book, id)
	if result.Error != nil {
		span.RecordError(result.Error)
		if errors.Is(result.Error, models.DB.Error) {
			return nil, fmt.Errorf("book not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get book: %v", result.Error)
	}

	return &book, nil
}

func (s *BookService) ListBooks(ctx context.Context) ([]models.Book, error) {
	ctx, span := s.tracer.Start(ctx, "ListBooks")
	defer span.End()

	var books []models.Book
	result := models.DB.WithContext(ctx).Find(&books)
	if result.Error != nil {
		span.RecordError(result.Error)
		return nil, fmt.Errorf("failed to list books: %v", result.Error)
	}

	span.SetAttributes(attribute.Int("books.count", len(books)))
	return books, nil
}

func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) error {
	ctx, span := s.tracer.Start(ctx, "UpdateBook")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("book.id", int64(book.ID)),
		attribute.String("book.title", book.Title),
		attribute.String("book.author", book.Author),
	)

	result := models.DB.WithContext(ctx).Save(book)
	if result.Error != nil {
		span.RecordError(result.Error)
		return fmt.Errorf("failed to update book: %v", result.Error)
	}

	return nil
}

func (s *BookService) DeleteBook(ctx context.Context, id uint) error {
	ctx, span := s.tracer.Start(ctx, "DeleteBook")
	defer span.End()

	span.SetAttributes(attribute.Int64("book.id", int64(id)))

	result := models.DB.WithContext(ctx).Delete(&models.Book{}, id)
	if result.Error != nil {
		span.RecordError(result.Error)
		return fmt.Errorf("failed to delete book: %v", result.Error)
	}

	return nil
}