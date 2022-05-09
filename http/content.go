package http

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micheam/notes"
)

type ContentController struct {
	contService *notes.ContentService
}

func NewContentController(svc *notes.ContentService) *ContentController {
	return &ContentController{svc}
}

func (ctrl ContentController) New(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer r.Body.Close()

	var rowBookID = r.URL.Query().Get("book_id")
	if len(rowBookID) == 0 {
		BadRequest(w, "missing param 'book_id'")
		return
	}
	bookID, err := notes.ParseBookID(rowBookID)
	if err != nil {
		BadRequest(w, "parse book_id: %v", err)
		return
	}

	var title = notes.Title(r.URL.Query().Get("title"))
	if len(title.String()) == 0 {
		BadRequest(w, "missing param 'title'")
		return
	}

	created, err := ctrl.contService.NewContent(r.Context(), bookID, title, r.Body)
	if err != nil {
		Error(w, fmt.Errorf("gen new book: %w", err))
		return
	}

	log.Printf("content created: %v\n", created.ID)
	JSON(w, http.StatusOK, toContentDetail(*created))
}

func (ctrl ContentController) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	ctx := r.Context()
	rowID := p.ByName("content_id")
	id, err := notes.ParseContentID(rowID)
	if err != nil {
		BadRequest(w, "invalid content_id: %v", err)
		return
	}

	cont, err := ctrl.contService.GetContent(ctx, id)
	if err != nil {
		if errors.Is(err, notes.ErrContentNotFound) {
			NotFound(w)
			return
		}
		BadRequest(w, err.Error())
		return
	}

	defer r.Body.Close()
	cont.Body = r.Body
	err = ctrl.contService.UpdateContent(ctx, cont)
	if err != nil {
		Error(w, err)
		return
	}

	JSON(w, http.StatusOK, toContentDetail(*cont))
}

func (ctrl ContentController) List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	ctx := r.Context()
	var rowBookID = r.URL.Query().Get("book_id")
	if len(rowBookID) == 0 {
		BadRequest(w, "missing param 'book_id'")
		return
	}
	bookID, err := notes.ParseBookID(rowBookID)
	if err != nil {
		BadRequest(w, "parse book_id: %v", err)
		return
	}

	got, err := ctrl.contService.ListContent(ctx, bookID)
	if err != nil {
		Error(w, fmt.Errorf("get book: %w", err))
		return
	}
	log.Printf("%v content found\n", got.Len())
	JSON(w, http.StatusOK, map[string]any{"list": toContentSummaryList(got)})
}

func (ctrl ContentController) View(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	ctx := r.Context()
	id, err := notes.ParseContentID(p.ByName("content_id"))
	if err != nil {
		BadRequest(w, "invalid content_id: '%s'", p.ByName("content_id"))
		return
	}
	got, err := ctrl.contService.GetContent(ctx, id)
	if err != nil {
		if errors.Is(err, notes.ErrContentNotFound) {
			NotFound(w)
			return
		}
		Error(w, fmt.Errorf("get content: %w", err))
		return
	}
	log.Printf("content found: %+v\n", got)
	JSON(w, http.StatusOK, toContentDetail(*got))
}

func (ctrl ContentController) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	ctx := r.Context()
	id, err := notes.ParseContentID(p.ByName("content_id"))
	if err != nil {
		BadRequest(w, "invalid content_id: '%s'", p.ByName("content_id"))
		return
	}
	err = ctrl.contService.DeleteContent(ctx, id)
	if err != nil && !errors.Is(err, notes.ErrContentNotFound) {
		Error(w, fmt.Errorf("delete content: %w", err))
		return
	}
	JSON(w, http.StatusOK, map[string]any{})
}

type ContentDetail struct {
	ContentSummary
	Body string `json:"body"`
}

func toContentDetail(cont notes.Content) ContentDetail {
	if cont.Body == nil {
		println("cont body is NULLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL")
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(cont.Body)
	return ContentDetail{
		ContentSummary: toContentSummary(cont),
		Body:           buf.String(),
	}
}

type ContentSummary struct {
	ID        string `json:"id"`
	Book      string `json:"book_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toContentSummary(cont notes.Content) ContentSummary {
	return ContentSummary{
		ID:        cont.ID.String(),
		Book:      cont.Parent.String(),
		Title:     cont.Title.String(),
		CreatedAt: FormatTime(cont.CreatedAt),
		UpdatedAt: FormatTime(cont.UpdatedAt),
	}
}

type ContentSummaryList []ContentSummary

func toContentSummaryList(list notes.ContentList) ContentSummaryList {
	sl := make(ContentSummaryList, len(list))
	for i, e := range list {
		sl[i] = toContentSummary(*e)
	}
	return sl
}
