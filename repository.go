package notes

import "context"

type BookRepository interface {
	GetBook(context.Context, BookID) (*Book, error)
	ListBook(context.Context) ([]*Book, error)
	InsertBook(context.Context, *Book) error
	UpdataBook(context.Context, *Book) error
	DeleteBook(context.Context, *Book) error
}
