package dicepk

import (
	"MyLine/dice"
	"MyLine/emoji"
	"MyLine/helper"
	"fmt"
	"log"
	"math"
	"myBot/mydb"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
)

type playerType struct {
	UserID      string
	DisplayName string
	HP          int
}

// DicepkType ...
type DicepkType struct {
	bets    int
	turn    int
	tmpDice struct {
		damage int
		boom   int
	}
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

type diceFaces struct {
	emoji string
	value string
}

// ...
var (
	DAMAGE     = diceFaces{emoji.Emoji(":crossed_swords:"), "damage"}
	FOOTPRINTS = diceFaces{emoji.Emoji(":footprints:"), "footprints"}
	BOOM       = diceFaces{emoji.Emoji(":boom:"), "boom"}
)

var diceFace = [6]diceFaces{DAMAGE, BOOM, FOOTPRINTS, FOOTPRINTS, BOOM, DAMAGE}

// Run ...
func (d *DicepkType) Run(input string) []string {
	var texts []string
	// New Game
	if len(d.players) < 2 {
		if strings.HasPrefix(input, "+") {
			d.joinGame()
			c := int(math.Floor(float64(d.turn % 2)))
			d.currentPlayer.id = c
			d.currentPlayer.UserID = d.players[d.turn].UserID
			switch len(d.players) {
			case 0, 1:
				re := regexp.MustCompile(`^\+(\d+)`)
				matches := re.FindStringSubmatch(input)
				d.bets = 1
				if len(matches) > 1 {
					if bet, err := strconv.Atoi(matches[1]); err == nil {
						d.bets = helper.Max(d.bets, bet)
					}
				}
				d.bets = helper.Min(d.bets, 1)
				return []string{fmt.Sprintf("%s 等待挑戰者...%s%d", d.players[0].DisplayName, emoji.Emoji(":money_bag:"), d.bets)}
			case 2:
				return []string{fmt.Sprintf("%s 加入戰局...\n%s%s 的回合", d.players[1].DisplayName, emoji.Emoji(":counterclockwise_arrows_button:"), d.players[c].DisplayName)}
			}
		}
	} else {
		if users.LineUser.UserProfile.UserID == d.currentPlayer.UserID {
			if strings.HasPrefix(input, "+") {
				c := int(math.Floor(float64(d.turn % 2)))
				pkDices := dice.Dice
				pkDices.Roll("3d6")
				for _, v := range pkDices.Rolls {
					//fmt.Printf(" %v", diceFace[v-1])
					switch diceFace[v-1] {
					case BOOM:
						d.tmpDice.boom++
					case DAMAGE:
						d.tmpDice.damage++
					}
				}
				text := fmt.Sprintf("%s(%d)%s(%d)%s", d.players[0].DisplayName, d.players[0].HP, emoji.Emoji(":VS_button:"), d.players[1].HP, d.players[1].DisplayName)
				text += fmt.Sprintf("\n%s 行動: %s %s %s", d.players[c].DisplayName, diceFace[pkDices.Rolls[0]-1].emoji, diceFace[pkDices.Rolls[1]-1].emoji, diceFace[pkDices.Rolls[2]-1].emoji)
				text += fmt.Sprintf("\n%s %d/%s %d", DAMAGE.emoji, d.tmpDice.damage, BOOM.emoji, d.tmpDice.boom)
				if d.tmpDice.boom >= 3 {
					text += d.endTurn()
				}
				texts = append(texts, text)
			} else if strings.HasPrefix(input, "-") {
				text := d.endTurn()
				texts = append(texts, text)
			}
		}
	}

	return texts
}

func (d *DicepkType) endTurn() string {
	var text string
	c := int(math.Floor(float64((d.turn) % 2)))
	t := int(math.Floor(float64((d.turn + 1) % 2)))
	if d.tmpDice.boom < 3 {
		if d.tmpDice.damage > 0 {
			d.players[t].damage(d.tmpDice.damage)
			text = fmt.Sprintf("%s(%d)%s(%d)%s", d.players[0].DisplayName, d.players[0].HP, emoji.Emoji(":VS_button:"), d.players[1].HP, d.players[1].DisplayName)
			text += fmt.Sprintf("\n%s 對 %s 造成傷害 %s %d", d.players[c].DisplayName, d.players[t].DisplayName, emoji.Emoji(":crossed_swords:"), d.tmpDice.damage)
		} else {
			text = fmt.Sprintf("%s(%d)%s(%d)%s", d.players[0].DisplayName, d.players[0].HP, emoji.Emoji(":VS_button:"), d.players[1].HP, d.players[1].DisplayName)
			text += fmt.Sprintf("\n%s 的攻擊未造成傷害", d.players[c].DisplayName)
		}
	} else {
		text = fmt.Sprintf("\n%s 攻擊失敗", d.players[c].DisplayName)
	}
	d.tmpDice.damage = 0
	d.tmpDice.boom = 0
	if d.players[t].HP <= 0 || d.players[c].HP <= 0 {
		if c == 1 {
			switch {
			case d.players[0].HP > d.players[1].HP:
				users.UsersList.Data[d.players[0].UserID].Money += d.bets
				users.UsersList.Data[d.players[1].UserID].Money -= d.bets
				text += fmt.Sprintf("\n勝 %s %d(%+d)%s(%+d)%d %s 負",
					d.players[0].DisplayName,
					users.UsersList.Data[d.players[0].UserID].Money,
					d.bets,
					emoji.Emoji(":money_bag:"),
					-(d.bets),
					users.UsersList.Data[d.players[1].UserID].Money,
					d.players[1].DisplayName)
			case d.players[0].HP < d.players[1].HP:
				users.UsersList.Data[d.players[1].UserID].Money += d.bets
				users.UsersList.Data[d.players[0].UserID].Money -= d.bets
				text += fmt.Sprintf("\n負 %s %d(%+d)%s(%+d)%d %s 勝",
					d.players[0].DisplayName,
					users.UsersList.Data[d.players[0].UserID].Money,
					-(d.bets),
					emoji.Emoji(":money_bag:"),
					d.bets,
					users.UsersList.Data[d.players[1].UserID].Money,
					d.players[1].DisplayName)
			case d.players[0].HP == d.players[1].HP:
				users.UsersList.Data[d.players[0].UserID].Money -= d.bets
				users.UsersList.Data[d.players[1].UserID].Money -= d.bets
				text += fmt.Sprintf("\n同歸餘盡\n%s %d(%+d)%s(%+d)%d %s",
					d.players[0].DisplayName,
					users.UsersList.Data[d.players[0].UserID].Money,
					-(d.bets),
					emoji.Emoji(":money_bag:"),
					-(d.bets),
					users.UsersList.Data[d.players[1].UserID].Money,
					d.players[1].DisplayName)
			}
			d.saveData()
			text += "\nGame Over"
			d.players = nil
			d.turn = 0
			return text
		}
	}
	d.turn++
	text += fmt.Sprintf("\n%s%s 回合", emoji.Emoji(":counterclockwise_arrows_button:"), d.players[int(math.Floor(float64((d.turn)%2)))].DisplayName)
	return text
}

func (d *DicepkType) joinGame() {
	d.players = append(d.players, &playerType{users.LineUser.UserProfile.UserID, users.LineUser.UserProfile.DisplayName, 15})
}

func (p *playerType) damage(dmg int) {
	p.HP -= dmg
}

func (d *DicepkType) saveData() {
	stmt, err := mydb.Db.Prepare("update users set money = $1 where userid = $2")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range d.players {
		_, err = stmt.Exec(users.UsersList.Data[v.UserID].Money, v.UserID)
		if err != nil {
			log.Fatal(err)
		}
	}
	stmt.Close()
}
