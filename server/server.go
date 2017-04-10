package server

import (
	"fmt"
	"net/http"
)



type HTTPServer struct {
	port int
}

func NewServer(port int) (*HTTPServer,error) {

	server := &HTTPServer{
		port: port,
	}
	return server,nil
}

func (p *HTTPServer) Run() error {
	addr := fmt.Sprintf(":%d", p.port)

	return http.ListenAndServe(addr, p)
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

