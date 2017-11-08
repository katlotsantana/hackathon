package request

type GameRule struct {	
	BoardWidth 		int `json:"boardWidth" binding:"required"`
	BoardHeight 		int `json:"boardHeight" binding:"required"`
	ShipTypes	[]ShipType 	`json:"ships"`
}


type ShipType struct {
	Type 		string 		`json:"type"`
	Quantity 	int 		`json:"quantity"`
}

type InviteGameRequest struct {
	SessionId     string `json:"sessionId" binding:"required"`
	GameRule 	GameRule `json:"gameRule"`
}