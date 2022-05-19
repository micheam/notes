-- +goose Up
-- +goose StatementBegin
CREATE TABLE book (
  id TEXT PRIMARY KEY CHECK( id != '' ),
  title TEXT CHECK( title != '' ),
  created_at TEXT CHECK( created_at != '' ),
  updated_at TEXT CHECK( created_at != '' )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book;
-- +goose StatementEnd
