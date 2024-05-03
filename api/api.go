package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nbitslabs/chaintips/storage"
)

type Server struct {
	db storage.Storage
}

func NewServer(db storage.Storage) *Server {
	return &Server{db: db}
}

func (s *Server) Serve() {
	router := gin.Default()

	router.GET("/chaintips/:id", s.GetChainTip)

	router.Run(":8080")
}

func (s *Server) GetChainTip(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid chain ID"})
		return
	}

	tip, err := s.db.GetNotableBlocks(intId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tip)
}
