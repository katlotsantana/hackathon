package request

type Player struct {
	Id 		string 		`json:"id" binding:"required"`
}

type StartGameRequest struct {
	SessionId 		string `json:"sessionId" binding:"required"`
	Player1 		Player `json:"player1"`
	Player2 		Player `json:"player2"`
}