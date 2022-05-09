package http

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/micheam/notes"
)

var (
	bookService *notes.BookService
	contService *notes.ContentService
	mu          sync.Mutex
)

func SetBookService(svc *notes.BookService) {
	mu.Lock()
	bookService = svc
	mu.Unlock()
}

func SetContentService(svc *notes.ContentService) {
	mu.Lock()
	contService = svc
	mu.Unlock()
}

func NewRouter() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "Welcome to Notes !\n")
	})

	//
	// NOTE:
	//   Since we intend to move to a gRPC-based system,
	//   we do not insist on following REST principles
	//   in our endpoint design.
	//

	// books
	books := &BookController{bookService: bookService}
	router.POST("/books", books.New)
	router.GET("/books", books.List)
	router.GET("/books/:book_id", books.View)
	router.PUT("/books/:book_id", books.Edit)
	router.DELETE("/books/:book_id", books.Delete)

	// content
	content := &ContentController{contService: contService}
	router.GET("/content", content.List)
	router.POST("/content", content.New)
	router.GET("/content/:content_id", content.View)
	router.PUT("/content/:content_id", content.Update)
	router.DELETE("/content/:content_id", content.Delete)

	return router
}
