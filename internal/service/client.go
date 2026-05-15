package service

import "github.com/chi-net/yapfa/gen/yapfa/v1/yapfa_v1connect"

type Client struct {
	handler yapfa_v1connect.YAPFAHandler
}

func NewClient(handler yapfa_v1connect.YAPFAHandler) *Client {
	return &Client{
		handler: handler,
	}
}
