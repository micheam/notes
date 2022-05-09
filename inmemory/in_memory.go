package inmemory

import (
	"context"
	"sync"

	"github.com/micheam/notes"
)

var books = sync.Map{}

type BookAccess struct{}

func (BookAccess) GetBook(_ context.Context, id notes.BookID) (*notes.Book, error) {
	got, ok := books.Load(id)
	if !ok {
		return nil, notes.ErrNotFound
	}
	return got.(*notes.Book), nil
}

func (BookAccess) ListBooks(_ context.Context) ([]*notes.Book, error) {
	var list []*notes.Book
	books.Range(func(_, book any) bool {
		list = append(list, book.(*notes.Book))
		return true
	})
	return list, nil
}

func (BookAccess) InsertBook(_ context.Context, book *notes.Book) error {
	books.Store(book.ID, book)
	return nil
}

func (BookAccess) UpdateBook(_ context.Context, book *notes.Book) error {
	_, ok := books.Load(book.ID)
	if !ok {
		return notes.ErrNotFound
	}
	books.Store(book.ID, book)
	return nil
}

func (BookAccess) DeleteBook(_ context.Context, book *notes.Book) error {
	_, ok := books.Load(book.ID)
	if !ok {
		return notes.ErrNotFound
	}
	books.Delete(book.ID)
	return nil
}
