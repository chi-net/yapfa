package handler

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) Delete(context.Context, *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	return &connect.Response[v1.DeleteResponse]{}, nil
}

func (h *Handler) Move(context.Context, *connect.Request[v1.MoveRequest]) (*connect.Response[v1.MoveResponse], error) {
	return &connect.Response[v1.MoveResponse]{}, nil
}

func (h *Handler) Copy(context.Context, *connect.Request[v1.CopyRequest]) (*connect.Response[v1.CopyResponse], error) {
	return &connect.Response[v1.CopyResponse]{}, nil
}

func (h *Handler) Mkdir(context.Context, *connect.Request[v1.MkdirRequest]) (*connect.Response[v1.MkdirResponse], error) {
	return &connect.Response[v1.MkdirResponse]{}, nil
}
