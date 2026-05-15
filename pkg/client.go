package pkg

import (
	"net/http"

	"github.com/chi-net/yapfa/gen/yapfa/v1/yapfa_v1connect"
)

type Conn struct {
	cli yapfa_v1connect.YAPFAClient
}

// Connect 用于开启服务器的连接
func Connect(host string) *Conn {
	cli := yapfa_v1connect.NewYAPFAClient(http.DefaultClient, host)
	return &Conn{cli: cli}
}
