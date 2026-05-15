package handler

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) List(context.Context, *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	return &connect.Response[v1.ListResponse]{}, nil
}
