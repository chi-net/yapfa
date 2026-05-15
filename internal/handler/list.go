package handler

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) List(ctx context.Context, req *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	target, err := h.safePath(req.Msg.Path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	files := make([]*v1.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, fileInfoFromOS(entry.Name(), filepath.Join(req.Msg.Path, entry.Name()), info))
	}

	return connect.NewResponse(&v1.ListResponse{Files: files}), nil
}

func fileInfoFromOS(name, path string, info fs.FileInfo) *v1.FileInfo {
	return &v1.FileInfo{
		Name:        name,
		Path:        path,
		Size:        info.Size(),
		IsDir:       info.IsDir(),
		Permissions: info.Mode().String(),
		ModifiedAt:  info.ModTime().Unix(),
	}
}
