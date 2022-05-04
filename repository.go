package notes

import "context"

type BookRepository interface {
	GetBook(context.Context, BookID) (*Book, error)
	ListBooks(context.Context) ([]*Book, error)
	InsertBook(context.Context, *Book) error
	UpdateBook(context.Context, *Book) error
	DeleteBook(context.Context, *Book) error
}
