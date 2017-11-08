package service

/*import (
	"hackathon/battleship/schema/request"
)*/

type Game struct {
	FirstPlayer  Player `json:"FirstPlayer"`
	SecondPlayer Player `json:"SecondPlayer"`
}

type Player struct {
	Id     				string        		`json:"id"`
	Sea      			Sea           		`json:"Sea"`
	Shoted				[]Vertex
	ShotedType 			map[string]int
	PotentialTarget  	map[string][]Vertex
	BeShoting			string
	// GunShot  int           `json:"GunShot"`
	// Suffered []Coordinates `json:"Suffered"`
}

type Sea struct {
	W 		int 	`json:"width"`
	H		int		`json:"height"`
	Girds	[]Gird 	
	Ships     []Ship `json:"ships"`
	ShipType 	map[string]int
	// ShipTypes	[]request.ShipType
}

type Gird struct {
	X int
	Y int
	Status string
}

type Ship struct {
	Type string           	`json:"type"`
	Positions []Vertex 		`json:"positions"`
}

type Vertex struct {
	X int `json:"x"`
	Y int `json:"y"`
}
