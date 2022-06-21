package sqlite

import (
	"context"
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
		Body:      row.Body,
		CreatedAt: parseDatetime(row.CreatedAt),
		UpdatedAt: parseDatetime(row.UpdatedAt),
	}
}

type ContentAccess struct {
	db *sqlx.DB
}

func NewContentAccess(db *sqlx.DB) *ContentAccess { return &ContentAccess{db} }

func (g ContentAccess) Get(ctx context.Context, id notes.ContentID) (*notes.Content, error) {
	query := `
        SELECT
          id, book_id, title, body, created_at, updated_at
        FROM general
        WHERE id=?;
    `
	var row = new(ContentRow)
	if err := g.db.GetContext(ctx, row, query, id); err != nil {
		return nil, err
	}
	return row.ToContent(), nil
}

func (g ContentAccess) Delete(ctx context.Context, cont *notes.Content) error {
	query := `
      DELETE FROM general
      WHERE id=:id;
    `
	_, err := g.db.NamedExecContext(ctx, query, cont)
	if err != nil {
		return err
	}
	return nil
}

func (g ContentAccess) Insert(ctx context.Context, cont *notes.Content) error {
	row := &ContentRow{
		ID:        cont.ID.String(),
		BookID:    cont.Parent.String(),
		Title:     cont.Title.String(),
		Body:      cont.Body,
		CreatedAt: formatDatetime(cont.CreatedAt),
		UpdatedAt: formatDatetime(cont.UpdatedAt),
	}
	query := `
        INSERT INTO general
          (id, book_id, title, body, created_at, updated_at)
        VALUES
          (:id, :book_id, :title, :body, :created_at, :updated_at);
    `
	_, err := g.db.NamedExecContext(ctx, query, row)
	return err
}

func (g ContentAccess) List(ctx context.Context) ([]*notes.Content, error) {
	query := `
        SELECT
          id, book_id, title, body, created_at, updated_at
        FROM general
        ORDER BY created_at ASC
    `
	var rows = []ContentRow{}
	if err := g.db.SelectContext(ctx, &rows, query); err != nil {
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
	query := `
        UPDATE general
        SET title=:title, body=:body, updated_at=:updated_at
        WHERE id=:id
    `
	_, err := g.db.NamedExecContext(ctx, query, ContentRow{
		ID:        cont.ID.String(),
		BookID:    cont.Parent.String(),
		Title:     cont.Title.String(),
		Body:      cont.Body,
		CreatedAt: formatDatetime(cont.CreatedAt),
		UpdatedAt: formatDatetime(cont.UpdatedAt),
	})
	if err != nil {
		return err
	}
	return nil
}