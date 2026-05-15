package handler

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) Stat(context.Context, *connect.Request[v1.StatRequest]) (*connect.Response[v1.StatResponse], error) {
	return &connect.Response[v1.StatResponse]{}, nil
}
