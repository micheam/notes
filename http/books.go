package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micheam/notes"
)

type BookController struct {
	bookService *notes.BookService
}

func NewBookController(svc *notes.BookService) *BookController {
	return &BookController{svc}
}

// New creates New Book, then return it.
//
// Request:
//   {
//      "title": "string",
//   }
func (ctrl BookController) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, "json body required")
		return
	}
	defer func() { _ = r.Body.Close() }()

	req := &BookNewRequest{}
	err = json.Unmarshal(b, req)
	if err != nil {
		BadRequest(w, "illegal format of request: %v", err)
		return
	}

	created, err := notes.NewBook(req.Title)
	if err != nil {
		Error(w, fmt.Errorf("gen new book: %w", err))
		return
	}
	if err := ctrl.bookService.SaveBook(r.Context(), created); err != nil {
		Error(w, fmt.Errorf("save book: %w", err))
		return
	}

	log.Printf("book created: %v\n", created)
	JSON(w, BookDetail{
		ID:        created.ID.String(),
		Title:     created.Title,
		CreatedAt: FormatTime(created.CreatedAt),
	})
}

func (ctrl BookController) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	foundBooks, err := ctrl.bookService.ListBooks(ctx)
	if err != nil {
		Error(w, fmt.Errorf("list books: %w", err))
		return
	}

	var books = make([]BookDetail, len(foundBooks))
	for i, b := range foundBooks {
		books[i] = BookDetail{
			ID:        b.ID.String(),
			Title:     b.Title,
			CreatedAt: FormatTime(b.CreatedAt),
		}
	}
	// if len(books) > 100 {
	// 	books = books[0:100]
	// }
	JSON(w, map[string]interface{}{"books": books})
}

func (ctrl BookController) View(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	bookID, err := notes.ParseBookID(p.ByName("book_id"))
	if err != nil {
		BadRequest(w, "invalid book_id")
		return
	}
	got, err := ctrl.bookService.GetBook(ctx, bookID)
	if err != nil {
		Error(w, fmt.Errorf("get book: %w", err))
		return
	}
	JSON(w, BookDetail{
		ID:        got.ID.String(),
		Title:     got.Title,
		CreatedAt: FormatTime(got.CreatedAt),
	})
}

func (ctrl BookController) Edit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, "json body required")
		return
	}
	defer func() { _ = r.Body.Close() }()

	req := &BookEditRequest{}
	err = json.Unmarshal(b, req)
	if err != nil {
		BadRequest(w, "illegal format of request: %v", err)
		return
	}

	// get existing book
	bookID, err := notes.ParseBookID(p.ByName("book_id"))
	if err != nil {
		BadRequest(w, "invalid book_id")
		return
	}
	book, err := ctrl.bookService.GetBook(ctx, bookID)
	if err != nil {
		Error(w, fmt.Errorf("get book: %w", err))
		return
	}

	book.Title = req.Title
	err = ctrl.bookService.SaveBook(ctx, book)
	if err != nil {
		Error(w, fmt.Errorf("save book: %w", err))
		return
	}

	JSON(w, BookDetail{
		ID:        book.ID.String(),
		Title:     book.Title,
		CreatedAt: FormatTime(book.CreatedAt),
	})
}

func (ctrl BookController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	bookID, err := notes.ParseBookID(p.ByName("book_id"))
	if err != nil {
		NotFound(w, "")
		return
	}
	existing, err := ctrl.bookService.GetBook(ctx, bookID)
	if err != nil {
		if errors.Is(err, notes.ErrNotFound) {
			JSON(w, map[string]any{}) // OK: Already gone
			return
		}
		Error(w, fmt.Errorf("get existing book: %w", err))
		return
	}
	err = ctrl.bookService.DeleteBook(ctx, existing)
	if err != nil {
		Error(w, fmt.Errorf("delete book: %w", err))
		return
	}
	JSON(w, map[string]any{})
}

type BookNewRequest struct {
	Title string `json:"title"`
}

type BookEditRequest struct {
	Title string `json:"title"`
}

type BookDetail struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}
