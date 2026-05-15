package handler

import (
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
)

// safePath 将用户传入的相对路径拼到 base 下，并防止路径穿越
func (h *Handler) safePath(rel string) (string, error) {
	joined := filepath.Join(h.base, filepath.Clean("/"+rel))
	if !strings.HasPrefix(joined, filepath.Clean(h.base)) {
		return "", connect.NewError(connect.CodeInvalidArgument, nil)
	}
	return joined, nil
}
