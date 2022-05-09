package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func BadRequest(w http.ResponseWriter, f string, args ...any) {
	JSON(w,
		http.StatusBadRequest,
		map[string]any{"message": fmt.Sprintf(f, args...)})
}

func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, map[string]any{})
}

func Error(w http.ResponseWriter, err error) {
	JSON(w,
		http.StatusInternalServerError,
		map[string]any{"message": err.Error()})
}

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ERROR] failed to marshal response data: %v", err)
		w.Write([]byte(`{}`))
		return
	}
	w.Write(b)
}
