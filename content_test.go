package notes_test

import (
	"errors"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/micheam/notes"
	. "github.com/micheam/notes/internal/testing"
)

func TestContentService_NewContent(t *testing.T) {
	t.Run("return error on invalid input", func(t *testing.T) {
		var (
			ctx   = genContext(t)
			book  = RandBook(t)
			title = notes.Title("")
		)
		defer ctx.Ctrl.Finish()
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		_, err := sut.NewContent(ctx.Context, book.ID, title, nil)
		// Verify ...
		assert.Error(t, err)
		assert.ErrorIs(t, err, notes.ErrInvalidTitle)
	})
	t.Run("return error if no book found", func(t *testing.T) {
		var (
			ctx   = genContext(t)
			book  = RandBook(t)
			title = notes.Title(RandStr(10))
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.BookRepo.EXPECT().
			GetBook(ctx.Context, book.ID).
			Return(nil, notes.ErrNotFound)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		_, err := sut.NewContent(ctx.Context, book.ID, title, nil)
		// Verify ...
		assert.Error(t, err)
		assert.ErrorIs(t, err, notes.ErrNotFound)
	})
	t.Run("success to register content", func(t *testing.T) {
		var (
			ctx   = genContext(t)
			book  = ctx.ExistingBook
			title = notes.Title(RandStr(10))
		)
		defer ctx.Ctrl.Finish()
		// Expectations
		ctx.ContRepo.EXPECT().
			Insert(ctx.Context, Any()).
			Return(nil)
		// Exercise
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		_, err := sut.NewContent(ctx.Context, book.ID, title, nil)
		// Verify
		assert.NoError(t, err)
	})
	t.Run("return error if content registration failed", func(t *testing.T) {
		var (
			ctx   = genContext(t)
			book  = ctx.ExistingBook
			title = notes.Title(RandStr(10))
			anerr = errors.New("!!FAILED TO INSERT CONTENT!!")
		)
		defer ctx.Ctrl.Finish()
		// Expectations
		ctx.ContRepo.EXPECT().
			Insert(ctx.Context, Any()).
			Return(anerr)
		// Exercise
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		_, err := sut.NewContent(ctx.Context, book.ID, title, nil)
		// Verify
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, anerr)
		}
	})
}

func TestContentService_UpdateContent(t *testing.T) {
	t.Run("Return error on the invalid content title", func(t *testing.T) {
		var (
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		cont.Title = notes.Title("") // Empty title
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.UpdateContent(ctx.Context, cont)
		// Verify ...
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, notes.ErrInvalidTitle)
		}
	})
	t.Run("Return error if target content does not exist", func(t *testing.T) {
		var (
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().
			Get(ctx.Context, cont.ID).
			Return(nil, notes.ErrNotFound)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.UpdateContent(ctx.Context, cont)
		// Verify ...
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, notes.ErrNotFound)
		}
	})
	t.Run("Update existing content", func(t *testing.T) {
		var (
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().
			Get(ctx.Context, cont.ID).
			Return(cont, nil)
		ctx.ContRepo.EXPECT().
			Update(ctx.Context, EqEntity(cont)). // TODO: Match with ID
			Return(nil)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.UpdateContent(ctx.Context, cont)
		// Verify ...
		assert.NoError(t, err)
	})
}

func TestContentService_DeleteContent(t *testing.T) {
	t.Run("Return error if content not found", func(t *testing.T) {
		var (
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().
			Get(ctx.Context, cont.ID).
			Return(nil, notes.ErrNotFound)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.DeleteContent(ctx.Context, cont.ID)
		// Verify ...
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, notes.ErrContentNotFound)
		}
	})
	t.Run("return error if faild to delete from storage", func(t *testing.T) {
		var (
			ctx            = genContext(t)
			cont           = RandContent(t)
			errFailedToDel = errors.New("failed to delete")
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().Get(ctx.Context, cont.ID).Return(cont, nil)
		ctx.ContRepo.EXPECT().Delete(ctx.Context, cont).Return(errFailedToDel)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.DeleteContent(ctx.Context, cont.ID)
		// Verify ...
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, errFailedToDel)
		}
	})
	t.Run("delete content", func(t *testing.T) {
		var (
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().Get(ctx.Context, cont.ID).Return(cont, nil)
		ctx.ContRepo.EXPECT().Delete(ctx.Context, cont).Return(nil)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		err := sut.DeleteContent(ctx.Context, cont.ID)
		// Verify ...
		assert.NoError(t, err)
	})
}

func TestContentService_ListContent(t *testing.T) {
	t.Run("return error if the specified Book does not exist", func(t *testing.T) {
		ctx := genContext(t)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.BookRepo.EXPECT().
			GetBook(ctx.Context, Any()).
			Return(nil, notes.ErrNotFound)
		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		_, err := sut.ListContent(ctx.Context, notes.NewBookID())
		// Verify ...
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, notes.ErrBookNotFound)
		}
	})
	t.Run("return found content", func(t *testing.T) {
		var ( // Setup ...
			ctx  = genContext(t)
			book = ctx.ExistingBook
			list = []*notes.Content{
				RandContent(t),
			}
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().List(ctx.Context, book).Return(list, nil)

		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		got, err := sut.ListContent(ctx.Context, book.ID)
		// Verify ...
		if assert.NoError(t, err) {
			assert.EqualValues(t, list, got)
		}
	})
}

func TestContentService_GetContent(t *testing.T) {
	t.Run("return found content", func(t *testing.T) {
		var ( // Setup ...
			ctx  = genContext(t)
			cont = RandContent(t)
		)
		defer ctx.Ctrl.Finish()
		// Expectations ...
		ctx.ContRepo.EXPECT().Get(ctx.Context, cont.ID).Return(cont, nil)

		// Exercise ...
		sut := notes.NewContentService(ctx.ContRepo, ctx.BookRepo)
		got, err := sut.GetContent(ctx.Context, cont.ID)
		// Verify ...
		if assert.NoError(t, err) {
			assert.EqualValues(t, cont, got)
		}
	})
}
