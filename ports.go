package notes

import "context"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -package notes_test -source=$GOFILE -destination ports_mock_test.go

type BookRepository interface {
	GetBook(context.Context, BookID) (*Book, error)
	ListBooks(context.Context) ([]*Book, error)
	InsertBook(context.Context, *Book) error
	UpdateBook(context.Context, *Book) error
	DeleteBook(context.Context, *Book) error
}

type ContentRepository interface {
	Get(context.Context, ContentID) (*Content, error)
	List(context.Context, *Book) ([]*Content, error)
	Insert(context.Context, *Content) error
	Update(context.Context, *Content) error
	Delete(context.Context, *Content) error
}
