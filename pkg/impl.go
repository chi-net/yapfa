package pkg

import (
	"context"
	"io"

	"connectrpc.com/connect"
	v1 "github.com/chi-net/yapfa/gen/yapfa/v1"
)

// List 列出指定路径下的文件和目录
func (c *Conn) List(ctx context.Context, path string) ([]*v1.FileInfo, error) {
	resp, err := c.cli.List(ctx, connect.NewRequest(&v1.ListRequest{Path: path}))
	if err != nil {
		return nil, err
	}
	return resp.Msg.Files, nil
}

// Stat 获取指定路径的文件或目录元信息
func (c *Conn) Stat(ctx context.Context, path string) (*v1.FileInfo, error) {
	resp, err := c.cli.Stat(ctx, connect.NewRequest(&v1.StatRequest{Path: path}))
	if err != nil {
		return nil, err
	}
	return resp.Msg.File, nil
}

// Download 下载指定路径的文件，offset 为断点续传起始字节（0 表示从头），返回 io.ReadCloser
func (c *Conn) Download(ctx context.Context, path string, offset int64) (io.ReadCloser, error) {
	stream, err := c.cli.Download(ctx, connect.NewRequest(&v1.DownloadRequest{
		Path:   path,
		Offset: offset,
	}))
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	go func() {
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {

			}
		}(pw)
		for stream.Receive() {
			if data := stream.Msg().GetData(); len(data) > 0 {
				if _, err := pw.Write(data); err != nil {
					err := pw.CloseWithError(err)
					if err != nil {
						return
					}
					return
				}
			}
		}
		if err := stream.Err(); err != nil {
			err := pw.CloseWithError(err)
			if err != nil {
				return
			}
		}
	}()
	return pr, nil
}

// Upload 上传文件，从 r 读取数据写入 path，overwrite 控制是否覆盖已有文件
func (c *Conn) Upload(ctx context.Context, path string, size int64, overwrite bool, r io.Reader) (int64, error) {
	stream := c.cli.Upload(ctx)

	// 先发元数据
	if err := stream.Send(&v1.UploadRequest{
		Payload: &v1.UploadRequest_Meta{
			Meta: &v1.UploadMeta{
				Path:      path,
				Size:      size,
				Overwrite: overwrite,
			},
		},
	}); err != nil {
		return 0, err
	}

	// 分块发送数据
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			if sendErr := stream.Send(&v1.UploadRequest{
				Payload: &v1.UploadRequest_Data{Data: buf[:n]},
			}); sendErr != nil {
				return 0, sendErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	resp, err := stream.CloseAndReceive()
	if err != nil {
		return 0, err
	}
	return resp.Msg.BytesWritten, nil
}

// Delete 删除指定路径的文件或目录，recursive 为 true 时相当于 rm -rf
func (c *Conn) Delete(ctx context.Context, path string, recursive bool) error {
	_, err := c.cli.Delete(ctx, connect.NewRequest(&v1.DeleteRequest{
		Path:      path,
		Recursive: recursive,
	}))
	return err
}

// Move 移动或重命名文件/目录，overwrite 控制目标已存在时是否覆盖
func (c *Conn) Move(ctx context.Context, from, to string, overwrite bool) error {
	_, err := c.cli.Move(ctx, connect.NewRequest(&v1.MoveRequest{
		From:      from,
		To:        to,
		Overwrite: overwrite,
	}))
	return err
}

// Copy 复制文件或目录，recursive 为 true 时支持复制目录
func (c *Conn) Copy(ctx context.Context, from, to string, recursive, overwrite bool) error {
	_, err := c.cli.Copy(ctx, connect.NewRequest(&v1.CopyRequest{
		From:      from,
		To:        to,
		Recursive: recursive,
		Overwrite: overwrite,
	}))
	return err
}

// Mkdir 创建目录，parents 为 true 时自动创建父目录（相当于 mkdir -p）
func (c *Conn) Mkdir(ctx context.Context, path string, parents bool) error {
	_, err := c.cli.Mkdir(ctx, connect.NewRequest(&v1.MkdirRequest{
		Path:    path,
		Parents: parents,
	}))
	return err
}
