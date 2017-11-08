package service

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"net/http"
)

import (
	"hackathon/battleship/schema/request"
)




func InviteGame(c *gin.Context) {
	log.Println("Hello, this is hackathon. I wanna invite game with you.")
	var req request.InviteGameRequest
	err := c.Bind(&req)
	if err == nil {
		var sea Sea
		sessionId := req.SessionId
		session := sessions.Default(c)
		old_sea := session.Get(sessionId)

		if old_sea == nil {
			// log.Println("create game.")
			sea = createSea(req.GameRule)
		} else {
			// log.Println("use session.")
			s, ok := old_sea.(Sea)
			if ok {
				sea = s
			}
			// log.Println("status: ", ok)
		}
		session.Set(sessionId, sea)
		session.Save()

		c.JSON(http.StatusOK, 
			gin.H{
			"success": "true",
			"number of girds": len(sea.Girds),
			"total ships": len(sea.Ships),
		})
	} else {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"success": "false"})
	}
}

func createSea(game_rule request.GameRule) (s Sea) {
	var ships []Ship;
	w := game_rule.BoardWidth
	h := game_rule.BoardHeight

	ship_types := game_rule.ShipTypes
	ship_type_map := make(map[string]int, len(ship_types))
	/*for t := 0; t < len(ship_types); t++ {
		ship_type := ship_types[t]
		ss := createShip(ship_type)
		ships = append(ships, ss...)
	}*/
	for _, st := range ship_types {
		ships = append(ships, createShip(st)...)
		ship_type_map[st.Type] = st.Quantity
	}

	var girds = make([]Gird, w*h)

	for i := 0; i <= w-1; i++ {
		for j := 0; j <= h-1; j++ {
			girds[i*h + j] = Gird{X: i, Y: j}
		}
	}

	s = Sea{
		W: w, 
		H: h, 
		Girds: girds, 
		Ships: ships,
		ShipType: ship_type_map,
	}
	return
}

func createShip(ship_type request.ShipType) (ss []Ship) {
	for i := 0; i < ship_type.Quantity; i++ {
		ship := Ship{Type: ship_type.Type}
		ss = append(ss, ship)
	}
	return 
}
