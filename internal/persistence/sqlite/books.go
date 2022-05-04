package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/micheam/notes"
)

type BookRow struct {
	ID        string `db:"id"`
	Title     string `db:"title"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (row BookRow) ToBook() *notes.Book {
	return &notes.Book{
		ID:        notes.BookID(row.ID),
		Title:     row.Title,
		CreatedAt: parseDatetime(row.CreatedAt),
		UpdatedAt: parseDatetime(row.UpdatedAt),
	}
}

type BookAccess struct {
	db *sqlx.DB
}

var _ notes.BookRepository = (*BookAccess)(nil)

func NewBookAccess(db *sqlx.DB) *BookAccess { return &BookAccess{db} }

func (b BookAccess) GetBook(ctx context.Context, id notes.BookID) (*notes.Book, error) {
	query := `
        SELECT 
            id, title,
            created_at,
            updated_at
        FROM book
        WHERE id=?
    `
	var row = new(BookRow)
	if err := b.db.GetContext(ctx, row, query, id); err != nil {
		return nil, err
	}
	return row.ToBook(), nil
}

func (b BookAccess) ListBooks(ctx context.Context) ([]*notes.Book, error) {
	query := `
        SELECT 
            id, title,
            created_at,
            updated_at
        FROM book
        ORDER BY created_at ASC
    `
	var rows = []BookRow{}
	if err := b.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}
	books := make([]*notes.Book, len(rows))
	for i := range rows {
		row := rows[i]
		books[i] = row.ToBook()
	}
	return books, nil
}

func (b BookAccess) InsertBook(ctx context.Context, book *notes.Book) error {
	query := `
        INSERT INTO book
            (id, title, created_at, updated_at)
            VALUES 
            (:id, :title, :created_at, :updated_at)`
	_, err := b.db.NamedExecContext(ctx, query, book)
	if err != nil {
		return err
	}
	return nil
}

func (b BookAccess) UpdateBook(ctx context.Context, book *notes.Book) error {
	book.UpdatedAt = time.Now()
	query := `
        UPDATE book 
        SET title=:title, updated_at=:updated_at
        WHERE id=:id
    `
	_, err := b.db.NamedExecContext(ctx, query, BookRow{
		ID:        book.ID.String(),
		Title:     book.Title,
		CreatedAt: formatDatetime(book.CreatedAt),
		UpdatedAt: formatDatetime(book.UpdatedAt),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b BookAccess) DeleteBook(ctx context.Context, book *notes.Book) error {
	query := `DELETE FROM book WHERE id=:id`
	_, err := b.db.NamedExecContext(ctx, query, book)
	if err != nil {
		return err
	}
	return nil
}
