package service

import (
	"context"
	"net/http"

	"github.com/chi-net/yapfa/gen/yapfa/v1/yapfa_v1connect"
)

func (c *Client) Start(ctx context.Context) {
	mux := http.NewServeMux()
	path, handler := yapfa_v1connect.NewYAPFAHandler(c.handler)
	mux.Handle(path, handler)
	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	s := &http.Server{
		Addr:      ":8080",
		Handler:   mux,
		Protocols: p,
	}
	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
