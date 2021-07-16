package controller

import (
	"talk_api/service"

	"github.com/gin-gonic/gin"
)

type SocketController interface {
	GetSocket(ctx *gin.Context)
}

type socketController struct {
	service service.SocketService
}

func New(service service.SocketService) SocketController {
	return &socketController{
		service: service,
	}
}

// Check MVC
func (c *socketController) GetSocket(ctx *gin.Context) {
	c.service.AddClient(ctx.Writer, ctx.Request)
}
