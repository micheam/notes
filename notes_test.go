package notes_test

import (
	"context"
	"testing"

	. "github.com/golang/mock/gomock"

	"github.com/micheam/notes"
	. "github.com/micheam/notes/internal/testing"
)

type Context struct {
	Context context.Context

	ExistingBook *notes.Book
	Ctrl         *Controller
	ContRepo     *MockContentRepository
	BookRepo     *MockBookRepository
}

func genContext(t *testing.T) Context {
	t.Helper()
	ctrl := NewController(t)
	contRepo := NewMockContentRepository(ctrl)
	bookRepo := NewMockBookRepository(ctrl)

	// inject existing book
	aBook := RandBook(t)
	bookRepo.EXPECT().
		GetBook(Any(), Eq(aBook.ID)).
		Return(aBook, nil).
		AnyTimes()

	return Context{
		Context:      context.Background(),
		Ctrl:         ctrl,
		ExistingBook: aBook,
		ContRepo:     contRepo,
		BookRepo:     bookRepo,
	}
}
