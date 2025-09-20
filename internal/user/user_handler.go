package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (s *Handler) CreateUser(c *gin.Context) {
	var user UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := s.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		// Return HTTP 500 status code with error message if service call fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return HTTP 200 status code with the response from the service
	c.JSON(http.StatusOK, res)

}
