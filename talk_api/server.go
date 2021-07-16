package main

import (
	"log"
	"net/http"
	"os"
	"talk/talk_api/con"
	"talk/talk_api/controller"
	"talk/talk_api/service"

	"github.com/gin-gonic/gin"
)

var (
	SocketService    service.SocketService       = con.New()
	socketController controller.SocketController = controller.New(SocketService)
)

func main() {
	//Check this
	go SocketService.Run()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := run(port); err != nil {
		log.Printf("error running server %s", err)
	}

}

func run(port string) error {
	health := HealthChecker{}
	router := gin.Default()
	mapRoutes(router, health)
	return router.Run(":" + port)
}

func mapRoutes(r *gin.Engine, health HealthChecker) {
	r.GET("/ping", health.PingHandler)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/msg", socketController.GetSocket)
	}
}

type HealthChecker struct{}

// Ping handler
func (h HealthChecker) PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
