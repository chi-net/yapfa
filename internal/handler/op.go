package handler

import (
	"context"
	"io"
	"os"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

func (h *Handler) Delete(ctx context.Context, req *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	target, err := h.safePath(req.Msg.Path)
	if err != nil {
		return nil, err
	}

	if req.Msg.Recursive {
		err = os.RemoveAll(target)
	} else {
		err = os.Remove(target)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.DeleteResponse{Success: true}), nil
}

func (h *Handler) Move(ctx context.Context, req *connect.Request[v1.MoveRequest]) (*connect.Response[v1.MoveResponse], error) {
	from, err := h.safePath(req.Msg.From)
	if err != nil {
		return nil, err
	}
	to, err := h.safePath(req.Msg.To)
	if err != nil {
		return nil, err
	}

	if !req.Msg.Overwrite {
		if _, err := os.Stat(to); err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, nil)
		}
	}

	if err := os.Rename(from, to); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.MoveResponse{Success: true}), nil
}

func (h *Handler) Copy(ctx context.Context, req *connect.Request[v1.CopyRequest]) (*connect.Response[v1.CopyResponse], error) {
	from, err := h.safePath(req.Msg.From)
	if err != nil {
		return nil, err
	}
	to, err := h.safePath(req.Msg.To)
	if err != nil {
		return nil, err
	}

	if err := copyPath(from, to, req.Msg.Recursive, req.Msg.Overwrite); err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.CopyResponse{Success: true}), nil
}

func (h *Handler) Mkdir(ctx context.Context, req *connect.Request[v1.MkdirRequest]) (*connect.Response[v1.MkdirResponse], error) {
	target, err := h.safePath(req.Msg.Path)
	if err != nil {
		return nil, err
	}

	if req.Msg.Parents {
		err = os.MkdirAll(target, 0755)
	} else {
		err = os.Mkdir(target, 0755)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.MkdirResponse{Success: true}), nil
}

// copyPath 递归复制文件或目录
func copyPath(from, to string, recursive, overwrite bool) error {
	info, err := os.Stat(from)
	if err != nil {
		return connect.NewError(connect.CodeNotFound, err)
	}

	if info.IsDir() {
		if !recursive {
			return connect.NewError(connect.CodeInvalidArgument, nil)
		}
		return copyDir(from, to, overwrite)
	}
	return copyFile(from, to, overwrite)
}

func copyFile(from, to string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(to); err == nil {
			return connect.NewError(connect.CodeAlreadyExists, nil)
		}
	}

	src, err := os.Open(from)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	return nil
}

func copyDir(from, to string, overwrite bool) error {
	if err := os.MkdirAll(to, 0755); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	entries, err := os.ReadDir(from)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	for _, entry := range entries {
		srcPath := from + "/" + entry.Name()
		dstPath := to + "/" + entry.Name()
		if err := copyPath(srcPath, dstPath, true, overwrite); err != nil {
			return err
		}
	}
	return nil
}
