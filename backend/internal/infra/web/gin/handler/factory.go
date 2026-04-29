package handler

import "sync"

var (
	serverHandler *ServerHandler
	once          sync.Once
)

type ServerHandler struct {
	UserHandler
	TaskHandler
	CsrfHandler
}

func NewHandler() *ServerHandler {
	once.Do(func() {
		serverHandler = &ServerHandler{}
	})
	return serverHandler
}

func (h *ServerHandler) Register(i any) *ServerHandler {
	switch interfaceType := i.(type) {
	case UserHandler:
		serverHandler.UserHandler = interfaceType
	case TaskHandler:
		serverHandler.TaskHandler = interfaceType
	case CsrfHandler:
		serverHandler.CsrfHandler = interfaceType
	}
	return serverHandler
}
