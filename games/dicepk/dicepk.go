package dicepk

import (
	"MyLine/dice"
	"fmt"
	"math"
	"strings"
)

type playerType struct {
	UserID      string
	DisplayName string
	HP          int
}

type DicepkType struct {
	trun          int
	currentPlayer struct {
		id     int
		UserID string
	}
	players []*playerType
}

// DicePK ...
var DicePK = make(map[string]*DicepkType)
var pler = struct {
	UserID      string
	DisplayName string
}{
	"userid",
	"catbox",
}

// Run ...
func (d *DicepkType) Run(input string) []string {
	if strings.HasPrefix(input, "+") {
		// New Game
		if len(d.players) < 2 {
			d.joinGame()
			c := int(math.Floor(float64(d.trun / 2)))
			d.currentPlayer.id = c
			d.currentPlayer.UserID = d.players[c].UserID
			fmt.Println(d.currentPlayer)
			return []string{"..."}
		}

		if "userid" == d.currentPlayer.UserID {
			pkDices := dice.Dice
			pkDices.Roll("3d6")
			fmt.Println("Roll")
			fmt.Println(pkDices)
		}
	}
	if strings.HasPrefix(input, "-") {
		d.players[0].damage()
	}

	return []string{"..."}
}

func (d *DicepkType) joinGame() {
	d.players = append(d.players, &playerType{"userid", "catbox", 15})
	d.players = append(d.players, &playerType{"userid2", "catbox2", 15})
}

func (p *playerType) damage() {
	p.HP--
}
