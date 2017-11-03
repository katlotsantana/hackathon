package main

import (
	"github.com/gin-gonic/gin"
)

import (
	"./service"
)

//var DB = make(map[string]string)

func main() {
	r := gin.Default()

	// Ping test
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	// // Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := DB[user]
	// 	if ok {
	// 		c.JSON(200, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(200, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	r.POST("/invite", service.InviteGame)
	r.POST("/start", service.StartGame)
	// r.POST("/turn", gameShot)
	// r.POST("/notify", gameNotify)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}


// func gameStart(c *gin.Context) {
// 	p0 := Position{0,0}
// 	var pos [] Position

// 	pos = append(pos, p0)
// 	ship1 := Ship{"AA", pos}

// 	c.JSON(http.StatusOK, gin.H{
// 		"ship tpe": ship1.ship_type,
// 		"ship's size": len(ship1.position)})
// 	// c.String(http.StatusOK, "hello battleship game.")
// }

// func inviteGame(c *gin.Context) {
// 	// session_id := c.PostForm("sessionId")
// 	// game_info := c.Params
// 	log.Println("invite game.")
// 	var req InviteGameRequest
// 	err := c.Bind(&req)
// 	if err == nil {
// 		c.JSON(http.StatusOK, 
// 			gin.H{
// 			"session Id": req.SessionId,
// 			"board width": req.GameRule.BoardWidth,
// 			"board height": req.GameRule.BoardHeight,
// 		})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	}
// }

func gameShot(c *gin.Context) {

}

func gameNotify(c *gin.Context) {

}

// type GameRule struct {	
// 	BoardWidth 		int `json:"boardWidth" binding:"required"`
// 	BoardHeight 		int `json:"boardHeight" binding:"required"`	
// }

// type InviteGameRequest struct {
// 	SessionId     string `json:"sessionId" binding:"required"`
// 	GameRule 	GameRule `json:"gameRule"`
// }

// type Ship struct {
// 	ship_type string
// 	position [] Position
// }

// type Position struct {
// 	x int
// 	y int
// }
