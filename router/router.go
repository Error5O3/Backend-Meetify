package router

import (
	"server/internal/event"
	"server/internal/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, eventHandler *event.Handler) {
	r = gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://192.168.4.76:8080"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {

			return origin == "http://192.168.4.76:8080"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/signup", userHandler.CreateUser)
	r.GET("/events/:event_id", eventHandler.GetEventGrid)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/event", eventHandler.CreateEvent)
	r.GET("/events/:event_id/grid", eventHandler.GetEventGrid)
	r.POST("/availability", eventHandler.MarkAvailable)
	r.DELETE("/availability", eventHandler.UnmarkAvailable)

}

func Start(address string) error {
	return r.Run(address)
}
