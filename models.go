package notes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// =====================================
//
// BookID
//
// =====================================

type BookID string

func (b BookID) String() string {
	return string(b)
}

func NewBookID() BookID {
	return BookID("book:" + uuid.New().String())
}

func ParseBookID(s string) (BookID, error) {
	pref, rowuuid, found := strings.Cut(s, ":")
	if !found {
		return "", fmt.Errorf("%w: missing prefix", ErrInvalidBookID)
	}
	if pref != "book" {
		return "", fmt.Errorf("%w: wrong prefix", ErrInvalidBookID)
	}
	_, err := uuid.Parse(rowuuid)
	if err != nil {
		return "", fmt.Errorf("parse uuid: %w", err)
	}
	return BookID("book:" + rowuuid), nil
}

func MustParseBookID(s string) BookID {
	bid, err := ParseBookID(s)
	if err != nil {
		panic(err)
	}
	return bid
}

// =====================================
//
// ContentID
//
// =====================================

type ContentID string

func NewContentID() ContentID {
	v := "content:" + uuid.New().String()
	return ContentID(v)
}

func ParseContentID(s string) (ContentID, error) {
	pref, rowuuid, found := strings.Cut(s, ":")
	if !found {
		return "", fmt.Errorf("%w: missing prefix", ErrInvalidContentID)
	}
	if pref != "content" {
		return "", fmt.Errorf("%w: wrong prefix", ErrInvalidContentID)
	}
	_, err := uuid.Parse(rowuuid)
	if err != nil {
		return "", fmt.Errorf("parse uuid: %w", err)
	}
	return ContentID("content:" + rowuuid), nil
}

func MustParseContentID(s string) ContentID {
	cid, err := ParseContentID(s)
	if err != nil {
		panic(err)
	}
	return cid
}

func (c ContentID) String() string {
	return string(c)
}

// =====================================
//
// Title
//
// =====================================

type Title string

func (t Title) Validate() error {
	if len(t) == 0 {
		return fmt.Errorf("%w: empty title", ErrInvalidTitle)
	}
	return nil
}

func (t Title) String() string {
	return string(t)
}

// =====================================
//
// TimeStamp
//
// =====================================

type TimeStamp struct {
	value time.Time
}

var (
	jstLocation *time.Location
	jstOnce     sync.Once
)

func JST() *time.Location {
	if jstLocation == nil {
		jstOnce.Do(func() {
			l, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				l = time.FixedZone("JST2", +9*60*60)
			}
			jstLocation = l
		})
	}
	return jstLocation
}

func (t TimeStamp) Format() string {
	return t.value.Format(time.RFC3339)
}

// =====================================
//
// Content
//
// =====================================

type Content struct {
	Parent    BookID
	ID        ContentID
	Title     Title
	Body      string // TOOD: Introduce Markdown format string types.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGeneralContent(p BookID, t Title) (*Content, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	id := NewContentID()
	now := time.Now()
	return &Content{
		Parent:    p,
		ID:        id,
		Title:     t,
		Body:      "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// =====================================
//
// Book
//
// =====================================

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
