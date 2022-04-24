package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

func main() {
	// err := cli.New().Run(os.Args)
	// if err != nil {
	// 	fmt.Println("ERROR: ", err)
	// 	os.Exit(1)
	// }

	conn, err := net.Dial("unix", socketAddr())
	if err != nil {
		panic(err)
	}

	name, _ := os.Hostname()
	req, err := http.NewRequest("GET", "http://localhost:8888/hello/"+name, nil)
	if err != nil {
		panic(err)
	}

	_ = req.Write(conn)
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}

func socketAddr() string {
	return filepath.Join(os.TempDir(), "notes.localserver.sock")
}
