package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/micheam/notes"
	"github.com/micheam/notes/internal/fileio"
	"github.com/micheam/notes/internal/localserver"
	"github.com/micheam/notes/internal/persistence/inmemory"
)

func init() {
	bookRepo := new(inmemory.BookAccess)
	bookSvc := notes.NewBookService(bookRepo)
	localserver.SetBookService(bookSvc)
}

func main() {
	if err := exec(); err != nil {
		log.Fatal(err)
	}
}

func exec() error {
	wd, err := fileio.PrepareWDir()
	if err != nil {
		return fmt.Errorf("prepare working dir: %w", err)
	}
	addr := filepath.Join(wd, "notes-localserver.sock")
	_ = os.Remove(addr)
	router := localserver.NewRouter()
	log.Printf("unix domain socket server start: %s", addr)
	return localserver.ListenAndServe(addr, router)
}
