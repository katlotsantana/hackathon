package request

type GameRule struct {	
	BoardWidth 		int `json:"boardWidth" binding:"required"`
	BoardHeight 		int `json:"boardHeight" binding:"required"`
	Ships	[] Ship 	`json:"ships"`
}


type Ship struct {
	Type 		string 		`json:"type"`
	Quantity 	int 		`json:"quantity"`
}

type InviteGameRequest struct {
	SessionId     string `json:"sessionId" binding:"required"`
	GameRule 	GameRule `json:"gameRule"`
}