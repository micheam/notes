package main

import (
	"log"

	"github.com/micheam/notes"
	"github.com/micheam/notes/internal/localserver"
	"github.com/micheam/notes/internal/persistence/inmemory"
)

func init() {
	bookRepo := new(inmemory.BookAccess)
	bookSvc := notes.NewBookService(bookRepo)
	localserver.SetBookService(bookSvc)
}

func main() {
	err := localserver.Start()
	if err != nil {
		log.Fatal(err)
	}
}
