package service

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
)

import (
	"hackathon/battleship/schema/request"
)


func StartGame(c *gin.Context) {
	log.Println("start game.")
	var req request.StartGameRequest
	err := c.Bind(&req)
	if err == nil {
		c.JSON(http.StatusOK, 
			gin.H{
			"session Id": req.SessionId,
			"player1": req.Player1.Id,
			"player2": req.Player2.Id,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}