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

var sea_W, sea_H int

func NotifyGame(c *gin.Context) {
	// log.Println("NOTYFY game.")
	var req request.NotifyGameRequest
	err := c.Bind(&req)
	if err == nil {
		session := sessions.Default(c)
		sessionId := req.SessionId

		current_session := session.Get(sessionId)
		game, ok := current_session.(Game)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}


		if game.FirstPlayer.Id == req.ShotResult.PlayerId {
			// log.Printf("notify player: %v", req.ShotResult.PlayerId)
			game.SecondPlayer.HanldeShotResult(req.ShotResult)
			log.Printf("shoted position: %v", game.SecondPlayer.Shoted)
			log.Println("=================================================================")
			log.Printf("potential target: %v", game.SecondPlayer.PotentialTarget)
			log.Printf("being shot: %v", game.SecondPlayer.BeShoting)
			// log.Printf("done notify %v", req.ShotResult.PlayerId)

		}
		session.Set(sessionId, game)
		session.Save()
		c.JSON(http.StatusOK, 
			gin.H{
			"success": true,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (player *Player) HanldeShotResult(shot_result request.ShotResult) {
	shoted_position := shot_result.Position
	sea_W = player.Sea.W - 1
	sea_H = player.Sea.H - 1

	shoted := &player.Shoted
	*shoted = append(*shoted, Vertex{X: shoted_position.X, Y:shoted_position.Y})


	tracking := updateStatusGirds(&player.Sea, shoted_position, shot_result.Status)
	log.Printf("update status for gird = %v, return %v", player.Sea.Girds[tracking], shot_result.Status)

	potential_target := player.PotentialTarget
	be_shoting := player.BeShoting

	if shot_result.Status == "HIT" {
		if !isEmpty(shot_result.RecognizedWholeShip) {
			whole_ship := shot_result.RecognizedWholeShip
			shoted_type := player.ShotedType
			if shoted_type == nil {
				shoted_type = make(map[string]int)
			}
			if count, ok := shoted_type[whole_ship.Type]; ok {
				count += 1
				shoted_type[whole_ship.Type] = count
			} else {
				shoted_type[whole_ship.Type] = 1
			}
			player.ShotedType = shoted_type
			log.Printf("whole ship type ===================== %v", whole_ship.Type)
			log.Printf("whole ship positions ===================== %v", whole_ship.Positions)

			keys_map := getAllKeyMap(potential_target)
			// keys_map_cp := keys_map
			for _, k := range keys_map {
				vertexs := extractKeyMap(k)
				for _, v := range vertexs {
					for _, p := range whole_ship.Positions {
						if v.X == p.X && v.Y == p.Y {
							delete(potential_target, k)
						}
					}
				}
			}
			keys_map = getAllKeyMap(potential_target)
			if len(keys_map) > 0 {
				next_key := keys_map[0]
				max_len := len(keys_map[0])
				for _, k := range keys_map {
					if len(k) > max_len {
						next_key = k
					}
				}
				player.BeShoting = next_key
				player.PotentialTarget = potential_target
			} else {
				player.BeShoting = ""
			}
		} else {
			if potential_target == nil {
				potential_target = make(map[string][]Vertex)
			}
			potential_target_array := CreatePotentialTargets(shoted_position, potential_target, be_shoting, *shoted)
			var key_map_for_single_hit string
			if be_shoting == "" {
				be_shoting = toKeyMap(shoted_position.X, shoted_position.Y)
			} else {
				key_map_for_single_hit = toKeyMap(shoted_position.X, shoted_position.Y)
				potential_target[key_map_for_single_hit] = CreatePotentialTargetWith1Hit(shoted_position, *shoted)
				be_shoting = be_shoting + "-" + key_map_for_single_hit
			}
			potential_target[be_shoting] = potential_target_array
			player.BeShoting = be_shoting
			player.PotentialTarget = potential_target
		}
	} else if shot_result.Status == "MISS" {
		for k, vertexs := range potential_target {
			for i, vertex := range vertexs {
				if vertex.X == shoted_position.X && vertex.Y == shoted_position.Y {
					vertexs = removeVertex(vertexs, i)
					break
				}
			}
			if len(vertexs) == 0 {
				delete(potential_target, k)
				if player.BeShoting == k {
					player.BeShoting = ""
				}
			} else {
				potential_target[k] = vertexs
				if player.BeShoting == "" {
					player.BeShoting = k
				}
			}
		}
	}
}

func updateStatusGirds(sea *Sea, shoted_position request.Position, status string) (tracking int) {
	for i := 0; i < len(sea.Girds); i++ {
		g := &sea.Girds[i]
		if g.X == shoted_position.X && g.Y == shoted_position.Y {
			g.Status = status
			tracking = i
		}
	}
	return
}

func CreatePotentialTargets(shoted_position request.Position, potential_target map[string][]Vertex, be_shoting string, shoted []Vertex) []Vertex {
	res_targets := make([]Vertex, 0)
	if len(potential_target) <= 0 {
		res_targets = append(res_targets, CreatePotentialTargetWith1Hit(shoted_position, shoted)...)
	} else {
		var potential_vetexes []Vertex = potential_target[be_shoting]
		if len(potential_vetexes) > 0 {
			for i, pos := range potential_vetexes {
				if pos.X == shoted_position.X && pos.Y == shoted_position.Y {
					potential_vetexes = removeVertex(potential_vetexes, i)
					break
				}
			}
			var vetexs []Vertex =  extractKeyMap(be_shoting)
			switch len(vetexs) {
			case 1:
				res_targets = CreatePotentialTargetWith2Hit(shoted_position, potential_vetexes, shoted)			
			case 2:
				res_targets = CreatePotentialTargetWith3Hit(shoted_position, vetexs, shoted)
			case 3:
				res_targets = CreatePotentialTargetWith4Hit(shoted_position, vetexs)
			}

		}	
	}

	return res_targets
}

func CreatePotentialTargetWith4Hit(shoted_position request.Position, shooting_vetexes []Vertex) []Vertex {
	res_targets := make([]Vertex, 0)
	if shooting_vetexes[0].X == shoted_position.X {
		if shoted_position.Y > shooting_vetexes[0].Y {
			res_targets = append(res_targets, Vertex{shoted_position.X - 1, shoted_position.Y - 2})
		} else {
			res_targets = append(res_targets, Vertex{shoted_position.X - 1, shoted_position.Y + 1})
		}
	} else {
		if shoted_position.X > shooting_vetexes[0].X {
			res_targets = append(res_targets, Vertex{shoted_position.X - 2, shoted_position.Y - 1})
		} else {
			res_targets = append(res_targets, Vertex{shoted_position.X + 1, shoted_position.Y - 1})
		}
	}

	return res_targets
}

func CreatePotentialTargetWith3Hit(shoted_position request.Position, shooting_vetexes []Vertex, shoted []Vertex) []Vertex {
	res_targets := make([]Vertex, 0)
	var next_target Vertex

	// On a line
	if (shooting_vetexes[0].X == shoted_position.X && shooting_vetexes[1].X == shoted_position.X) || (shooting_vetexes[0].Y == shoted_position.Y && shooting_vetexes[1].Y == shoted_position.Y) {
		if shooting_vetexes[0].X == shoted_position.X {
			// shoted at behide 
			if shoted_position.Y < shooting_vetexes[0].Y {
				if shoted_position.Y > consts.MIN_HEIGHT {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shoted_position.Y - 1})
				}
				if shooting_vetexes[0].Y < sea_H {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shooting_vetexes[0].Y + 1})
				}
				
			} else if shoted_position.Y > shooting_vetexes[0].Y {	// shoted at up
				if shoted_position.Y < sea_H {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shoted_position.Y + 1})
				}
				if shooting_vetexes[0].Y > consts.MIN_HEIGHT {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shooting_vetexes[0].Y - 1})
				}				
			}
		} else {
			// shoted at left
			if shoted_position.X < shooting_vetexes[0].X {
				if shoted_position.X > consts.MIN_WIDTH {
					res_targets = append(res_targets, Vertex{shoted_position.X - 1, shooting_vetexes[0].Y})
				}
				if shooting_vetexes[0].X < sea_W {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X + 1, shooting_vetexes[0].Y})
				}
			} else if shoted_position.X > shooting_vetexes[0].X {	// shoted at right
				if shoted_position.X < sea_W {
					res_targets = append(res_targets, Vertex{shoted_position.X + 1, shooting_vetexes[0].Y})
				}
				if shooting_vetexes[0].X > consts.MIN_WIDTH {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X - 1, shooting_vetexes[0].Y})
				}				
			}
		}
	} else {
		if shooting_vetexes[0].X == shooting_vetexes[1].X {
			if shooting_vetexes[0].Y < shooting_vetexes[1].Y {
				next_target = Vertex{shoted_position.X, shoted_position.Y - 1}
			} else {
				next_target = Vertex{shoted_position.X, shoted_position.Y + 1}
			}
			// case carrier ship
			if shoted_position.X < shooting_vetexes[0].X {
				if shoted_position.Y <= sea_H - 2 {
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shoted_position.Y + 2})
				} 
				if shoted_position.Y >= 1 {
					
					res_targets = append(res_targets, Vertex{shooting_vetexes[0].X, shoted_position.Y - 1})
				}
			}
		} else if shooting_vetexes[0].Y == shooting_vetexes[1].Y {
			if shooting_vetexes[0].X < shooting_vetexes[1].X {
				next_target = Vertex{shoted_position.X - 1, shoted_position.Y}
			} else {
				next_target = Vertex{shoted_position.X + 1, shoted_position.Y}
			}
			// case carrier ship
			if shoted_position.Y < shooting_vetexes[0].Y {
				if shoted_position.X <= sea_W - 2 {
					res_targets = append(res_targets, Vertex{shoted_position.X + 2, shooting_vetexes[0].Y})
				} 
				if shoted_position.X >= 1 {					
					res_targets = append(res_targets, Vertex{shoted_position.X - 1, shooting_vetexes[0].Y})
				}
			}
		}
		if !isVertexEmpty(next_target) {
			res_targets = append(res_targets, next_target)
		}
	}


	return res_targets
}

func CreatePotentialTargetWith2Hit(shoted_position request.Position, potential_vetexes []Vertex, shoted []Vertex) []Vertex {
	res_targets := make([]Vertex, 0)
	res_targets = append(res_targets, CreatePotentialTargetWith1Hit(shoted_position, shoted)...)
	res_targets = append(res_targets, potential_vetexes...)

	return res_targets
}

func CreatePotentialTargetWith1Hit(shoted_position request.Position, shoted []Vertex) []Vertex {
	res_targets := make([]Vertex, 0)
	var next_target Vertex
	if shoted_position.X < sea_W {
		next_target = Vertex{X: shoted_position.X + 1, Y:shoted_position.Y}
		if !CheckPositionShooted(shoted, next_target) {
			 res_targets = append(res_targets, next_target)
		}
	}
	if shoted_position.Y < sea_H {
		next_target = Vertex{X: shoted_position.X, Y:shoted_position.Y + 1}
		if !CheckPositionShooted(shoted, next_target) {
			 res_targets = append(res_targets, next_target)
		}
	}
	if shoted_position.X > consts.MIN_WIDTH {
		next_target = Vertex{X: shoted_position.X - 1, Y:shoted_position.Y}
		if !CheckPositionShooted(shoted, next_target) {
			 res_targets = append(res_targets, next_target)
		}
	}
	if shoted_position.Y > consts.MIN_HEIGHT {
		next_target = Vertex{X: shoted_position.X, Y:shoted_position.Y - 1}
		if !CheckPositionShooted(shoted, next_target) {
			 res_targets = append(res_targets, next_target)
		}
	}

	return res_targets
}

func CheckPositionShooted(shoted []Vertex, intend_shot Vertex) bool {
	for _, sh := range shoted {
		if sh.X == intend_shot.X && sh.Y == intend_shot.Y {
			return true
		}
	}

	return false
}
