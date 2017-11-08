package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"encoding/gob"
)

import (
	"hackathon/battleship/service"
	// "hackathon/battleship/schema"
)

// func init() {
// 	gob.Register(&schema.Gird{})
// 	gob.Register(&schema.Ship{})
// 	gob.Register(&schema.Coordinates{})
// 	gob.Register(&schema.Sea{})
// 	gob.Register(&schema.Player{})
// 	gob.Register(&schema.Game{})
// }

func init() {
	gob.Register(service.Sea{})
	gob.Register(service.Game{})
}

const REDIS_SERVER = "localhost:6379"

func main() {
	r := gin.Default()
	// store := sessions.NewCookieStore([]byte("something-very-secret"))
	// r.Use(sessions.Sessions("battleship-session", store))

	// sessions.Options = &sessions.Options{
	//     Path:     "/",
	//     MaxAge:   86400,
	//     HttpOnly: true,
	// }

	store, _ := sessions.NewRedisStore(10, "tcp", REDIS_SERVER, "", []byte("something-very-secret"))
	// options := sessions.Options{
	//     Path:     "/",
	//     MaxAge:   86400,
	//     HttpOnly: true,
	// }
	// store.Options(sessions.Options{
	//     Path:     "/",
	//     MaxAge:   86400,
	//     HttpOnly: true,
	// })
	r.Use(sessions.Sessions("battleship-session", store))

	r.POST("/invite", service.InviteGame)
	r.POST("/start", service.StartGame)
	r.POST("/turn", service.TurnGame)
	r.POST("/notify", service.NotifyGame)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
