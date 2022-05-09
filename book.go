package notes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Book struct {
	ID        BookID    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewBook(title string) (*Book, error) {
	if len(title) == 0 {
		err := errors.New("empty")
		return nil, NewValidationError("title", err)
	}
	now := time.Now()
	return &Book{
		ID:        NewBookID(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (b Book) Validate() error {
	if len(b.Title) == 0 {
		err := errors.New("empty string")
		return NewValidationError("title", err)
	}
	return nil
}

type BookService struct {
	books BookRepository
}

func NewBookService(bookRepo BookRepository) *BookService {
	return &BookService{bookRepo}
}

func (b BookService) SaveBook(ctx context.Context, book *Book) error {
	if err := book.Validate(); err != nil {
		return err
	}

	got, err := b.books.GetBook(ctx, book.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get book: %w", err)
		}
		return b.books.InsertBook(ctx, book)
	}

	got.Title = book.Title
	got.UpdatedAt = time.Now()
	return b.books.UpdateBook(ctx, got)
}

func (b BookService) DeleteBook(ctx context.Context, book *Book) error {
	return b.books.DeleteBook(ctx, book)
}

func (b BookService) GetBook(ctx context.Context, id BookID) (*Book, error) {
	return b.books.GetBook(ctx, id)
}

func (b BookService) ListBooks(ctx context.Context) ([]*Book, error) {
	return b.books.ListBooks(ctx)
}
