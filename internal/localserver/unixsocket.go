package localserver

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/micheam/notes/internal/fileio"
)

type UDSServer struct {
	*http.Server
}

func (srv *UDSServer) ListenAndServe() error {
	addr := srv.Addr
	ln, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func Start() error {
	wd, err := fileio.PrepareWDir()
	if err != nil {
		return fmt.Errorf("prepare working dir: %w", err)
	}
	addr := socketAddr(wd)
	_ = os.Remove(addr)
	router := NewRouter()
	return ListenAndServe(addr, router)
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}
	udss := &UDSServer{Server: server}
	return udss.ListenAndServe()
}

func socketAddr(prefix string) string {
	return filepath.Join(prefix, "notes-localserver.sock")
}
