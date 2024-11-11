// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package book

import (
	"context"
)

// Book contains data about a book.
type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Interface exposed by the Book service.
type Interface interface {
	Create(ctx context.Context, book *Book) error
	Get(ctx context.Context, id uint) (*Book, error)
	List(ctx context.Context) ([]Book, error)
	UpdateBook(ctx context.Context, book *Book) error
	DeleteBook(ctx context.Context, id uint) error
}
