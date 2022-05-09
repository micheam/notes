package notes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

type Content struct {
	Parent    BookID
	ID        ContentID
	Title     Title
	Body      io.Reader // TODO: Introduce Markdown format string types.
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContentList []*Content

func (cl ContentList) Len() int {
	return len(cl)
}

func (c Content) Validate() error {
	if err := c.Title.Validate(); err != nil {
		return err
	}
	return nil
}

func (c Content) Read(p []byte) (n int, err error) {
	return c.Body.Read(p)
}

type ContentService struct {
	contRepository ContentRepository
	bookRepository BookRepository
}

func NewContentService(c ContentRepository, b BookRepository) *ContentService {
	return &ContentService{c, b}
}

// NewContent generates new content and returns the result.
//
// Return an error if any of the following incorrect conditions are detected.
//
//   * Specified book does not exist.
//   * Invalid content title
//
// Note that title validation is equivalent to Title.Validate().
func (cs ContentService) NewContent(
	ctx context.Context, bookID BookID, title Title, body io.Reader) (*Content, error) {

	if bookID.Empty() {
		return nil, fmt.Errorf("book_id: %w", ErrInvalidArgument)
	}
	if err := title.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	cont := &Content{
		Parent:    bookID,
		ID:        NewContentID(),
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := cs.bookRepository.GetBook(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("get book: %w", err)
	}
	err = cs.contRepository.Insert(ctx, cont)
	if err != nil {
		return nil, fmt.Errorf("insert content: %w", err)
	}
	// TODO: send log to sync server
	return cont, nil
}

// UpdateContent edit an content and register it.
//
// Return an error if any of the following incorrect conditions are detected.
//
//   * Specified content does not exist.
//   * Invalid content property
//
// Note that title validation is equivalent to Title.Validate().
func (cs ContentService) UpdateContent(ctx context.Context, cont *Content) error {

	if err := cont.Validate(); err != nil {
		return err
	}

	_, err := cs.contRepository.Get(ctx, cont.ID)
	if err != nil {
		return fmt.Errorf("get existing content: %w", err)
	}

	cont.UpdatedAt = time.Now()
	err = cs.contRepository.Update(ctx, cont)
	if err != nil {
		return fmt.Errorf("update content: %w", err)
	}
	// TODO: send log to sync server
	return nil
}

func (cs ContentService) DeleteContent(ctx context.Context, id ContentID) error {

	cont, err := cs.contRepository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			err = ErrContentNotFound
		}
		return fmt.Errorf("get content(%v): %w", id, err)
	}

	err = cs.contRepository.Delete(ctx, cont)
	if err != nil {
		return fmt.Errorf("delete content(%v): %w", id, err)
	}
	// TODO: send log to sync server
	return nil
}

func (cs ContentService) ListContent(ctx context.Context, bookID BookID) (ContentList, error) {
	foundBook, err := cs.bookRepository.GetBook(ctx, bookID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			err = ErrBookNotFound
		}
		return nil, fmt.Errorf("get book(%v): %w", bookID, err)
	}
	return cs.contRepository.List(ctx, foundBook)
}

func (cs ContentService) GetContent(ctx context.Context, id ContentID) (*Content, error) {
	return cs.contRepository.Get(ctx, id)
}
