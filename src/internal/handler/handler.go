package handler

import "go.uber.org/dig"

type Handler struct {
	UserHandler UserHandler
}

type HandlerParams struct {
	dig.In

	UserHandler UserHandler
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		UserHandler: params.UserHandler,
	}
}
