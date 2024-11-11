package models

import (
	"context"
)

type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Interface interface {
	CreateBook(ctx context.Context, book *Book) error
	GetBook(ctx context.Context, id uint) (*Book, error)
	ListBooks(ctx context.Context) ([]Book, error)
	UpdateBook(ctx context.Context, book *Book) error
	DeleteBook(ctx context.Context, id uint) error
}