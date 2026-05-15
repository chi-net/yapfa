package handler

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) Upload(context.Context, *connect.ClientStream[v1.UploadRequest]) (*connect.Response[v1.UploadResponse], error) {
	return &connect.Response[v1.UploadResponse]{}, nil
}

func (h *Handler) Download(context.Context, *connect.Request[v1.DownloadRequest], *connect.ServerStream[v1.DownloadResponse]) error {
	return nil
}
