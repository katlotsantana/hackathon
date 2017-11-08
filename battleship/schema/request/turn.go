package request


type TurnGameRequest struct {
	SessionId 		string `json:"sessionId" binding:"required"`
	TurnNumber 		int `json:"turnNumber" binding:"required"`
}