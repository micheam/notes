package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/micheam/notes"
)

type ContentRow struct {
	ID        string `db:"id"`
	BookID    string `db:"book_id"`
	Title     string `db:"title"`
	Body      string `db:"body"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (row ContentRow) ToContent() *notes.Content {
	return &notes.Content{
		ID:        notes.ContentID(row.ID),
		Parent:    notes.BookID(row.BookID),
		Title:     notes.Title(row.Title),
		Body:      strings.NewReader(row.Body),
		CreatedAt: parseDatetime(row.CreatedAt),
		UpdatedAt: parseDatetime(row.UpdatedAt),
	}
}

func (row *ContentRow) FromContent(cont notes.Content) error {
	*row = ContentRow{
		ID:        cont.ID.String(),
		BookID:    cont.Parent.String(),
		Title:     cont.Title.String(),
		CreatedAt: formatDatetime(cont.CreatedAt),
		UpdatedAt: formatDatetime(cont.UpdatedAt),
	}
	if cont.Body == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(cont.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	row.Body = buf.String()
	return nil
}

type ContentAccess struct {
	db *sqlx.DB
}

func NewContentAccess(db *sqlx.DB) *ContentAccess { return &ContentAccess{db} }

func (g ContentAccess) Get(ctx context.Context, id notes.ContentID) (*notes.Content, error) {
	query := `
        SELECT
          id, book_id, title, body, created_at, updated_at
        FROM content
        WHERE id=?;
    `
	var row = new(ContentRow)
	if err := g.db.GetContext(ctx, row, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, notes.ErrContentNotFound
		}
		return nil, err
	}
	return row.ToContent(), nil
}

func (g ContentAccess) Delete(ctx context.Context, cont *notes.Content) error {
	query := `
      DELETE FROM content
      WHERE id=:id;
    `
	_, err := g.db.NamedExecContext(ctx, query, cont)
	if err != nil {
		return err
	}
	return nil
}

func (g ContentAccess) Insert(ctx context.Context, cont *notes.Content) error {
	row := new(ContentRow)
	err := row.FromContent(*cont)
	if err != nil {
		return fmt.Errorf("row from content: %w", err)
	}
	query := `
        INSERT INTO content
          (id, book_id, title, body, created_at, updated_at)
        VALUES
          (:id, :book_id, :title, :body, :created_at, :updated_at);
    `
	_, err = g.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return err
	}
	return nil
}

func (g ContentAccess) List(ctx context.Context, book *notes.Book) ([]*notes.Content, error) {
	query := `
        SELECT
          id, book_id, title, body, created_at, updated_at
        FROM content
        WHERE book_id=?
        ORDER BY created_at ASC
    `
	var rows = []ContentRow{}
	if err := g.db.SelectContext(ctx, &rows, query, book.ID); err != nil {
		return nil, err
	}
	list := make([]*notes.Content, len(rows))
	for i := range rows {
		row := rows[i]
		list[i] = row.ToContent()
	}
	return list, nil
}

func (g ContentAccess) Update(ctx context.Context, cont *notes.Content) error {
	cont.UpdatedAt = time.Now()

	var row = new(ContentRow)
	err := row.FromContent(*cont)
	if err != nil {
		return fmt.Errorf("row from content: %w", err)
	}

	query := `
        UPDATE content
        SET title=:title, body=:body, updated_at=:updated_at
        WHERE id=:id
    `
	_, err = g.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return err
	}
	return nil
}
