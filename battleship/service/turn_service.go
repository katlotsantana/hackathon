package service

import (
	"log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

import (
	"hackathon/battleship/schema/request"
)


func TurnGame(c *gin.Context) {
	var req request.TurnGameRequest
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

		var shot_pos Vertex
		if len(game.SecondPlayer.PotentialTarget) > 0 {
			shot_pos = TargetShot(&game.SecondPlayer)
		} else {
			shot_pos = HunterShot(&game.SecondPlayer)
		}
		log.Printf("shot at vertext: %v", shot_pos)

		c.JSON(http.StatusOK, 
			gin.H{
			"firePosition": shot_pos,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func TargetShot(player *Player) Vertex {
	log.Println("........................Starting target shot ........................")
	potential_target := player.PotentialTarget

	var target Vertex
	var be_shoting string
	var potentials []Vertex = potential_target[player.BeShoting]
	if len(potentials) > 0 {
		switch vetexs := extractKeyMap(be_shoting); len(vetexs) {
		case 1:
			target = potentials[0]
		case 2:
			if vetexs[0].X == vetexs[1].X {
				for _, n := range potentials {
					if n.X == vetexs[0].X {
						target = n
						break						
					}					
				}
			} else if vetexs[0].Y == vetexs[1].Y {
				for _, n := range potentials {
					if n.Y == vetexs[0].Y {
						target = n
						break						
					}					
				}
			}
			if isVertexEmpty(target) {
				log.Println("Target is empty.--------------------")
				target = potentials[0]
			}
		case 3:
			if vetexs[0].X == vetexs[1].X && vetexs[1].Y == vetexs[2].Y {
				for _, n := range potentials {
					if n.Y == vetexs[0].Y {
						target = n
						break
					}
				}
			} else if vetexs[0].Y == vetexs[1].Y && vetexs[1].X == vetexs[2].X {
				for _, n := range potentials {
					if n.X == vetexs[0].X {
						target = n
						break
					}
				}
			} else {
				if vetexs[0].X == vetexs[1].X {
					for _, n := range potentials {
						if n.X == vetexs[0].X {
							target = n
							break
						}
					}
				} else if vetexs[0].Y == vetexs[1].Y {
					for _, n := range potentials {
						if n.Y == vetexs[0].Y {
							target = n
							break
						}
					}
				}
			}
			if isVertexEmpty(target) {
				target = potentials[0]
			}
		case 4:
			target = potentials[0]
		default:
			target = potentials[0]
		}		
	} else {
		for k := range potential_target {
			var pts []Vertex = potential_target[k]
	        if len(pts) > 0 {
	        	target = pts[0]
	        	be_shoting = k
	        	break
	        }
	    }
	    player.BeShoting = be_shoting
	}
	
	if CheckPositionShooted(player.Shoted, target) {
		for i, pos := range potentials {
			if pos.X == target.X && pos.Y == target.Y {
				potentials = removeVertex(potentials, i)
				break
			}
		}
		potential_target[player.BeShoting] = potentials
		player.PotentialTarget = potential_target
		return TargetShot(player)
	}

	return target
}

func HunterShot(player *Player) (Vertex) {
	log.Println("........................Starting hunter shot ........................")
	try_number := 1
	try_map := make(map[int][]Vertex, try_number)

	for t := 0; t < try_number; t++ {
		try_map[t] = CreateShips(player.Sea, player.ShotedType)
	}

	arr1 := try_map[0]
	/*arr2 := try_map[1]
	res_arr := make([]Vertex, 0)
	for _, v1 := range arr1 {
		for _, v2 := range arr2 {
			if v1.X == v2.X && v1.Y == v2.Y {
				res_arr = append(res_arr, v1)
			}
		}
	}
	log.Printf("arr1: %v", arr1)
	log.Printf("arr2: %v", arr2)

	if len(res_arr) == 0 {
		return HunterShot(player)
		// return arr1[0]
	}

	res_index := random(0, len(res_arr)-1)
	target := res_arr[res_index]
	*/
	if len(arr1) == 0 {
		return HunterShotRandom(player)
	}
	target := arr1[random(0, len(arr1)-1)]
	if CheckPositionShooted(player.Shoted, target) {
		log.Println("........................this position shoted. Retry it ........................")
		return HunterShot(player)
	}
	return target
}

func HunterShotRandom(player *Player) (Vertex) {
	log.Println("........................Starting hunter shot ........................")
	girds := player.Sea.Girds
	indexs := make([]int, 0, len(girds)/2)
	if len(player.Shoted) == 0 {
		girds := player.Sea.Girds
		for i, _ := range girds {
			if i%2 == 0 {
				indexs = append(indexs, i)
			}
		}
	} else {
		var odd_index []int
		var even_index []int
		for i, g := range player.Sea.Girds {
			if g.Status == "" {
				if i%2 == 0 {
					even_index = append(even_index, i)
				} else {
					odd_index = append(odd_index, i)
				}
			}
		}
		if len(even_index) <= len(odd_index) {
			indexs = append(indexs, even_index...)
		} else {
			indexs = append(indexs, odd_index...)
		}
	}

	random_shot := random(indexs[0], indexs[len(indexs)-1])
	r_gird := girds[random_shot]
	var target = Vertex{X: r_gird.X, Y: r_gird.Y}
	if CheckPositionShooted(player.Shoted, target) {
		
		return HunterShotRandom(player)
	}

	return target
}
