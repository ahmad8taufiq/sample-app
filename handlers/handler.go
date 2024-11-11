// handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sample-app/models"
	"sample-app/services"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// CreateBook handles the creation of a new book
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("http-handler").Start(r.Context(), "CreateBookHandler")
	defer span.End()

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.bookService.CreateBook(ctx, &book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// GetBook handles retrieving a book by ID
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("http-handler").Start(r.Context(), "GetBookHandler")
	defer span.End()

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBook(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// ListBooks handles retrieving all books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("http-handler").Start(r.Context(), "ListBooksHandler")
	defer span.End()

	books, err := h.bookService.ListBooks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// UpdateBook handles updating a book
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("http-handler").Start(r.Context(), "UpdateBookHandler")
	defer span.End()

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	book.ID = uint(id)

	if err := h.bookService.UpdateBook(ctx, &book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DeleteBook handles deleting a book
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("http-handler").Start(r.Context(), "DeleteBookHandler")
	defer span.End()

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.bookService.DeleteBook(ctx, uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}