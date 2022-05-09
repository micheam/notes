-- +goose Up
create table content (
  id text primary key check( id != '' ),
  book_id text,
  title text check( title != '' ),
  body text,
  created_at text check( created_at != '' ),
  updated_at text check( created_at != '' ),

  foreign key (book_id) references book(id)
);

-- +goose Down
drop table content;
