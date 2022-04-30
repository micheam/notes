package localserver

import (
	"net"
	"net/http"
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

func ListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}
	udss := &UDSServer{Server: server}
	return udss.ListenAndServe()
}
