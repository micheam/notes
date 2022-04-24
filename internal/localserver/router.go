package localserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/micheam/notes"
)

var (
	bookService *notes.BookService
	mu          sync.Mutex
)

func SetBookService(svc *notes.BookService) {
	mu.Lock()
	bookService = svc
	mu.Unlock()
}

func NewRouter() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "Welcome to Notes !\n")
	})

	books := &BookController{bookService: bookService}
	router.POST("/books", books.New)
	router.GET("/books", books.List)
	router.GET("/books/:book_id", books.View)
	router.PUT("/books/:book_id", books.Edit)
	router.DELETE("/books/:book_id", books.Delete)

	return router
}
