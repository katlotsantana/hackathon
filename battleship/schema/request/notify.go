package request


type ShotResult struct {
	PlayerId		string 				`json:"playerId" binding:"required"`
	Position		Position			`json:"position" binding:"required"`
	Status			string				`json:"status" binding:"required"`
	RecognizedWholeShip		RecognizedWholeShip 	`json:"recognizedWholeShip"`
}

type Position struct {
	X 			int						`json:"x"`
	Y 			int						`json:"y"`
}

type RecognizedWholeShip struct {
	Type 			string 				`json:"type"`
	Positions		[]Position 			`json:"positions"`
}

type NotifyGameRequest struct {
	SessionId     			string 					`json:"sessionId" binding:"required"`
	ShotResult 				ShotResult 				`json:"shotResult" binding:"required"`
}