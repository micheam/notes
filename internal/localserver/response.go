package localserver

import (
	"bytes"
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
	status := http.StatusBadRequest
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	msg := new(bytes.Buffer)
	_, err := fmt.Fprintf(msg, f, args...)
	if err != nil {
		log.Printf("[ERROR] failed to write response: %v", err)
		w.Write([]byte(http.StatusText(status)))
		return
	}
	w.Write(msg.Bytes())
}

func NotFound(w http.ResponseWriter, f string, args ...any) {
	status := http.StatusNotFound
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	msg := new(bytes.Buffer)
	_, err := fmt.Fprintf(msg, f, args...)
	if err != nil {
		log.Printf("[ERROR] failed to write response: %v", err)
		w.Write([]byte(http.StatusText(status)))
		return
	}
	w.Write(msg.Bytes())
}

func Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	log.Printf("[ERROR] unhandle error: %v", err)
	b := []byte(err.Error())
	w.Write(b)
}

func JSON(w http.ResponseWriter, data any) {
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ERROR] failed to marshal response data: %v", err)
		w.Write([]byte(`{}`))
		return
	}
	w.Write(b)
}
