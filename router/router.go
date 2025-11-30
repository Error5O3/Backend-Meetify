package router

import (
	"server/internal/user"
	"server/internal/event"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, eventHandler *event.Handler) {
	r = gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.LoginUser)
	r.POST("/event", eventHandler.CreateEvent)
	r.GET("/events/:event_id/grid" ,eventHandler.GetEventGrid)
	r.POST("/availability", eventHandler.MarkAvailable)
	r.DELETE("/availability",eventHandler.UnmarkAvailable)

}

func Start(address string) error {
	return r.Run(address)
}
