package main

import "github.com/huangrao121/CommunicationApp/BackendService/internal/message/handler"

type HandlerInit struct {
	messageHandler *handler.MessageHandler
}

func NewHandlerInit(messageHandler *handler.MessageHandler) *HandlerInit {
	return &HandlerInit{
		messageHandler: messageHandler,
	}
}
