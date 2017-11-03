package service

import (
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
)

import (
	"hackathon/battleship/schema/request"
)


func InviteGame(c *gin.Context) {
	log.Println("invite game.")
	var req request.InviteGameRequest
	err := c.Bind(&req)
	if err == nil {
		c.JSON(http.StatusOK, 
			gin.H{
			"session Id": req.SessionId,
			"board width": req.GameRule.BoardWidth,
			"board height": req.GameRule.BoardHeight,
			"ship count": len(req.GameRule.Ships),
			"ship": req.GameRule.Ships[0].Quantity,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}