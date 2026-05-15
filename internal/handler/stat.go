package handler

import (
	"context"
	"os"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) Stat(ctx context.Context, req *connect.Request[v1.StatRequest]) (*connect.Response[v1.StatResponse], error) {
	target, err := h.safePath(req.Msg.Path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.StatResponse{
		File: fileInfoFromOS(info.Name(), req.Msg.Path, info),
	}), nil
}
