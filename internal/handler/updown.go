package handler

import (
	"context"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

const defaultChunkSize = 32 * 1024

func (h *Handler) Download(ctx context.Context, req *connect.Request[v1.DownloadRequest], stream *connect.ServerStream[v1.DownloadResponse]) error {
	target, err := h.safePath(req.Msg.Path)
	if err != nil {
		return err
	}

	f, err := os.Open(target)
	if err != nil {
		if os.IsNotExist(err) {
			return connect.NewError(connect.CodeNotFound, err)
		}
		return connect.NewError(connect.CodeInternal, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	// 断点续传
	offset := req.Msg.Offset
	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	// 探测 MIME 类型
	mimeType := mime.TypeByExtension(filepath.Ext(target))
	if mimeType == "" {
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	// 先发元数据
	if err := stream.Send(&v1.DownloadResponse{
		Payload: &v1.DownloadResponse_Meta{
			Meta: &v1.DownloadMeta{
				Name:     info.Name(),
				Size:     info.Size(),
				MimeType: mimeType,
			},
		},
	}); err != nil {
		return err
	}

	chunkSize := req.Msg.ChunkSize
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	buf := make([]byte, chunkSize)
	for {
		select {
		case <-ctx.Done():
			return connect.NewError(connect.CodeCanceled, ctx.Err())
		default:
		}

		n, err := f.Read(buf)
		if n > 0 {
			if sendErr := stream.Send(&v1.DownloadResponse{
				Payload: &v1.DownloadResponse_Data{Data: buf[:n]},
			}); sendErr != nil {
				return sendErr
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}
}

func (h *Handler) Upload(ctx context.Context, stream *connect.ClientStream[v1.UploadRequest]) (*connect.Response[v1.UploadResponse], error) {
	// 第一条消息必须是 meta
	if !stream.Receive() {
		return nil, connect.NewError(connect.CodeInvalidArgument, stream.Err())
	}
	meta := stream.Msg().GetMeta()
	if meta == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	target, err := h.safePath(meta.Path)
	if err != nil {
		return nil, err
	}

	if !meta.Overwrite {
		if _, err := os.Stat(target); err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, nil)
		}
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	f, err := os.Create(target)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer f.Close()

	var written int64
	for {
		select {
		case <-ctx.Done():
			return nil, connect.NewError(connect.CodeCanceled, ctx.Err())
		default:
		}

		if !stream.Receive() {
			break
		}
		data := stream.Msg().GetData()
		if len(data) == 0 {
			continue
		}
		n, err := f.Write(data)
		written += int64(n)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if err := stream.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.UploadResponse{BytesWritten: written}), nil
}
