package service


import (
	"log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

import (
	"hackathon/battleship/consts"
	"hackathon/battleship/schema/request"
)


const (
	CARRIER = "CV"
	BATTLE_SHIP = "BB"
	OIL_RIG = "OR"
	DESTROYER = "DD"
	CRUISER = "CA"
)

// type ShipCreator interface {
// 	Create() []Vertex
// }

type ShipPattern struct {
	Id 			string
	Name 		string
	Demension	int
}

func Create(ship ShipPattern) (positions []Vertex) {
	positions = make([]Vertex, ship.Demension)

	for i := 0; i < len(positions); i++ {
		if i == 0 {
			x := random(0, 19)
			y := random(0, 7)
			positions[i] = Vertex{X:x, Y:y}
			// log.Printf("x = %v, y = %v", x, y)
		} else {
			switch ship.Id {
			case CARRIER:
				if i == len(positions) - 1 {
					positions[i] = Vertex{X: positions[1].X, Y: positions[2].Y - 1}
				} else {
					positions[i] = Vertex{X: positions[i-1].X + 1, Y: positions[i-1].Y}
				}
			case BATTLE_SHIP:
				positions[i] = Vertex{X: positions[i-1].X + 1, Y: positions[i-1].Y}
			case OIL_RIG:
				if i == 1 {
					positions[i] = Vertex{X: positions[i-1].X + 1, Y: positions[i-1].Y}
				} else if i == 2 {
					positions[i] = Vertex{X: positions[i-1].X, Y: positions[i-1].Y + 1}
				} else if i == 3 {
					positions[i] = Vertex{X: positions[i-1].X - 1, Y: positions[i-1].Y}
				}
			case DESTROYER:				
				positions[i] = Vertex{X: positions[i-1].X, Y: positions[i-1].Y + 1}
			case CRUISER:
				positions[i] = Vertex{X: positions[i-1].X, Y: positions[i-1].Y + 1}
			}
		}
	}

	// log.Printf("Positions before return : %v.", positions)

	return positions
}

var ship_pattern_map = map[string]ShipPattern{
	"CV": {Id: CARRIER, Name: "Carrier", Demension: 5},
	"BB": {Id: BATTLE_SHIP, Name: "Battle Ship", Demension: 4},
	"OR": {Id: OIL_RIG, Name: "Oil Rig", Demension: 4},
	"DD": {Id: DESTROYER, Name: "Destroyer", Demension: 2},
	"CA": {Id: CRUISER, Name: "Cruiser", Demension: 3},
}

var sea_width, sea_height int

func StartGame(c *gin.Context) {
	// log.Println("start game.")
	var req request.StartGameRequest
	err := c.Bind(&req)
	if err == nil {
		// var sea schema.Sea
		session := sessions.Default(c)
		sessionId := req.SessionId

		current_session := session.Get(sessionId)
		sea, ok := current_session.(Sea)
		if !ok {
			game, ok := current_session.(Game)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			sea = game.FirstPlayer.Sea
		}
		sea_width = sea.W - 1
		sea_height = sea.H - 1
		// log.Println("load from session ok.")

		sea.createShips()
		game := Game{
			FirstPlayer: Player{
				Id: req.Player1.Id, 
				Sea: sea,
			},
			SecondPlayer: Player{
				Id: req.Player2.Id,
				Sea: Sea{
					W: sea.W,
					H: sea.H,
					Girds: sea.Girds,
					Ships: sea.Ships,
					ShipType: sea.ShipType,
				},
			},
		}

		session.Set(sessionId, game)
		session.Save()

		c.JSON(http.StatusOK, 
			gin.H{
			"session Id": req.SessionId,
			"ships": sea.Ships,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}


func (sea *Sea) createShips() /*map[string]Vertex*/ {
	var ships_map = make(map[string]Vertex)
	for i := 0; i < len(sea.Ships); i++ {
		ship := &sea.Ships[i]
		ship.createOnMap(ships_map, ship_pattern_map[ship.Type])
		log.Println("ships map: ", ships_map)
	}

	// return ships_map
}



func (ship *Ship) createOnMap(ships_map map[string]Vertex, ship_pattern ShipPattern) {
	log.Printf("Starting create position for ship: %v.", ship.Type)
	// log.Println("currently ships map: %v.", ships_map)

	// var ship_creator ShipCreator = ship_pattern
	var positions []Vertex
	for {
		// ship_creator.Create()
		positions = Create(ship_pattern)
		// log.Println("created position for ship: %v.", positions)

		if checkOverlap(ships_map, positions) {
			break
		}
	}

	for _, pos := range positions {
		ships_map[toKeyMap(pos.X, pos.Y)] = pos
	}
	ship.Positions = positions
	return
}

/*
	Check list positions have overlap or not
*/

func checkOverlap(ships_map map[string]Vertex, positions_check []Vertex) bool {
	status := true
	for _, p := range positions_check {
		/*
			Check a vertext outsite sea map.
		*/
		if p.X > sea_width || p.Y > sea_height {
			status = false
			break
		}
		if p.X < 0 || p.Y < 0 {
			status = false
			break
		}
		if _, ok := ships_map[toKeyMap(p.X, p.Y)]; ok {
			status = false
			break
		}
		if _, ok := ships_map[toKeyMap(p.X + consts.DISTANCE, p.Y)]; ok {
			status = false
			break
		}
		if _, ok := ships_map[toKeyMap(p.X, p.Y + consts.DISTANCE)]; ok {
			status = false
			break
		}
		if _, ok := ships_map[toKeyMap(p.X - consts.DISTANCE, p.Y)]; ok {
			status = false
			break
		}
		if _, ok := ships_map[toKeyMap(p.X, p.Y - consts.DISTANCE)]; ok {
			status = false
			break
		}
	}
	// log.Println("status check pos: %v", status)
	return status
}

/* ================================================================== */

func CreateShips(sea Sea, ships_type_shooted map[string]int) (res_positions []Vertex) {
	res_positions = make([]Vertex, 0)
	var ships []Ship = make([]Ship, 0)
	ship_types_defined := copyMap(sea.ShipType)
	log.Printf("================= ships defined ======================== : %v", ship_types_defined)
	log.Printf("================= ships shoted ======================== : %v", ships_type_shooted)
	for _, item :=  range sea.Ships {
		count, ok := ships_type_shooted[item.Type]
		if  !ok {
			ships = append(ships, item)
		}
		if ok && count < ship_types_defined[item.Type] {
			ships = append(ships, item)
			ship_types_defined[item.Type] = ship_types_defined[item.Type] - 1
		}
	}
	log.Printf("================= ships defined ======================== : %v", ship_types_defined)
	log.Printf(" ================ total ships count ================ : %v", len(sea.Ships))
	log.Printf(" ================ ships remains count ================ : %v", len(ships))

	ships_map := make(map[string]Vertex)
	for _, ship := range ships {
		res_positions = append(res_positions, CreateShipOnMap(ship, ships_map, ship_pattern_map[ship.Type], sea)...)
	}

	// log.Printf("res positions: %v", res_positions)

	return
}

func CreateShipOnMap(ship Ship, ship_map map[string]Vertex, ship_pattern ShipPattern, sea Sea) (res_positions []Vertex) {
	
	count := 1
	for {
		if count > 100 {
			break
		}
		log.Printf("---------try create ship ---------------%v------------", count)
		res_positions = Create(ship_pattern)
		if Validate(sea, res_positions, ship_map, true) {
			break
		}
		res_positions = nil
		count += 1
	}
	
	return res_positions
}

func Validate(sea Sea, positions []Vertex, ships_map map[string]Vertex, is_check_on_map bool) bool {
	status := true
	temp_positions := GetPositionAroundShip(positions)
	for _, p := range temp_positions {
		if p.IsOutSideMap() {
			status = false
			break
		}
		if p.IsOverlap(ships_map) {
			status = false
			break
		}
		if is_check_on_map {
			if !p.IsValidOnMap(sea.Girds) {
				status = false
				break
			}			
		}
	}

	return status
}

func GetPositionAroundShip(positions []Vertex) (res_positions []Vertex) {
	res_positions = make([]Vertex, 0)
	var temp_vertex Vertex
	for _, p := range positions {
		if !CheckVertexInArray(res_positions, p) {
			if p.X < 19 {
				temp_vertex = Vertex{p.X + 1, p.Y}
				res_positions = append(res_positions, temp_vertex)
			} 
			if p.X > consts.MIN_WIDTH { 
				temp_vertex = Vertex{p.X - 1, p.Y}
				res_positions = append(res_positions, temp_vertex)
			}
			if p.Y < 7 {
				temp_vertex = Vertex{p.X, p.Y + 1}
				res_positions = append(res_positions, temp_vertex)
			}
			if p.Y > consts.MIN_HEIGHT {
				temp_vertex = Vertex{p.X, p.Y - 1}
				res_positions = append(res_positions, temp_vertex)
			}
			/*if p.X < sea_width && p.Y < sea_height {
				temp_vertex = Vertex{p.X + 1, p.Y + 1}
				res_positions = append(res_positions, temp_vertex)
			}
			if p.X < sea_width && p.Y > 0 {
				temp_vertex = Vertex{p.X + 1, p.Y - 1}
				res_positions = append(res_positions, temp_vertex)
			}
			if p.X > 0 && p.Y < sea_height {
				temp_vertex = Vertex{p.X - 1, p.Y + 1}
				res_positions = append(res_positions, temp_vertex)
			}
			if p.X > 0 && p.Y > 0 {
				temp_vertex = Vertex{p.X - 1, p.Y - 1}
				res_positions = append(res_positions, temp_vertex)
			}*/
		}
	}
	return
}

func CheckVertexInArray(arr []Vertex, v Vertex) bool {
	for _, p := range arr {
		if p.X == v.X && p.Y == v.Y {
			return true
		}
	}
	return false
}

func (position Vertex) IsOutSideMap() bool {
	res := false
	if position.X > 19 || position.Y > 7 || position.X < 0 || position.Y < 0 {
		res = true
	}
	return res
}

func (p Vertex) IsOverlap(ships_map map[string]Vertex) bool {
	if _, ok := ships_map[toKeyMap(p.X, p.Y)]; ok {
		return true
	}
	if _, ok := ships_map[toKeyMap(p.X + consts.DISTANCE, p.Y)]; ok {
		return true
	}
	if _, ok := ships_map[toKeyMap(p.X, p.Y + consts.DISTANCE)]; ok {
		return true
	}
	if _, ok := ships_map[toKeyMap(p.X - consts.DISTANCE, p.Y)]; ok {
		return true
	}
	if _, ok := ships_map[toKeyMap(p.X, p.Y - consts.DISTANCE)]; ok {
		return true
	}

	return false
}

func (position Vertex) IsValidOnMap(girds []Gird) bool {
	is_ok := true
	for _, g := range girds {
		if position.X == g.X && position.Y == g.Y && g.Status != "" {
			is_ok = false
			break
		}
	}
	return is_ok
}
