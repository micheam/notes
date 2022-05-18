package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/micheam/notes"
	"github.com/micheam/notes/internal/fileio"
	"github.com/micheam/notes/internal/localserver"
	"github.com/micheam/notes/internal/sqlite"
)

var db *sqlx.DB
var wdir string

func init() {
	var err error
	wdir, err = fileio.PrepareWDir()
	if err != nil {
		panic(fmt.Errorf("prepare working dir: %w", err))
	}

	datafile, err := datasource()
	if err != nil {
		panic(err)
	}
	db, err = sqlite.Open(datafile)
	if err != nil {
		panic(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		panic(err)
	}

	bookRepo := sqlite.NewBookAccess(db)
	bookSvc := notes.NewBookService(bookRepo)
	localserver.SetBookService(bookSvc)
}

func main() {
	defer func() { _ = db.Close() }()
	if err := exec(); err != nil {
		log.Fatal(err)
	}
}

func exec() error {
	addr := filepath.Join(wdir, "localserver.sock")
	_ = os.Remove(addr)
	router := localserver.NewRouter()
	log.Printf("unix domain socket server start: %s", addr)
	return localserver.ListenAndServe(addr, router)
}

// TODO(micheam): Rethink where settings and data files are stored.
func datasource() (path string, err error) {
	hm := os.Getenv("HOME")
	if len(hm) == 0 {
		return "", fmt.Errorf("env HOME is empty")
	}
	return filepath.Join(hm, ".notes", "books.db"), nil
}
