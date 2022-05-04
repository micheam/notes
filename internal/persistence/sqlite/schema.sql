CREATE TABLE book (
  id TEXT PRIMARY KEY CHECK( id != '' ),
  title TEXT CHECK( title != '' ),
  created_at INTEGER CHECK( created_at != 0 ),
  updated_at INTEGER CHECK( created_at != 0 )
)
