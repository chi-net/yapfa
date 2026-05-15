package handler

import "github.com/chi-net/yapfa/gen/yapfa/v1/yapfa_v1connect"

type Handler struct {
	base string
}

func New(base string) *Handler {
	return &Handler{
		base: base,
	}
}

func (h *Handler) Instance() yapfa_v1connect.YAPFAHandler {
	return h
}
