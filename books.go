package notes

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        BookID
	Title     string
	CreatedAt time.Time
}

func NewBook(title string) (*Book, error) {
	if len(title) == 0 {
		err := errors.New("empty")
		return nil, NewValidationError("title", err)
	}
	return &Book{
		ID:        NewBookID(),
		Title:     title,
		CreatedAt: time.Now(),
	}, nil
}

func (b Book) Valid() error {
	if len(b.Title) == 0 {
		err := errors.New("empty string")
		return NewValidationError("title", err)
	}
	return nil
}

type BookID string

func (b BookID) String() string {
	return string(b)
}

func NewBookID() BookID {
	return BookID(uuid.New().String())
}

func ParseBookID(s string) (BookID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return BookID(u.String()), nil
}

type BookService struct {
	books BookRepository
}

func NewBookService(bookRepo BookRepository) *BookService {
	return &BookService{bookRepo}
}

func (b BookService) SaveBook(ctx context.Context, book *Book) error {
	if err := book.Valid(); err != nil {
		return err
	}
	return b.books.InsertBook(ctx, book)
}

func (b BookService) DeleteBook(ctx context.Context, book *Book) error {
	return b.books.DeleteBook(ctx, book)
}

func (b BookService) GetBook(ctx context.Context, id BookID) (*Book, error) {
	return b.books.GetBook(ctx, id)
}

func (b BookService) ListBooks(ctx context.Context) ([]*Book, error) {
	return b.books.ListBook(ctx)
}
