package boomgame

import (
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/users"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type scene3InfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene3InfoType) startPhase(g *GameType) {
	var p []string
	for _, v := range g.data.players.Queue {
		p = append(p, v.DisplayName)
	}
	texts = append(texts,
		fmt.Sprintf("[%s核爆危機]\n這是一顆核彈，請在5回合內拆除...\n%s\n%s",
			emoji.Emoji(":radioactive:"),
			strings.Join(p, ", "),
			b.show(g)))
}

func (b *scene3InfoType) runPhase(input string, g *GameType) {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
			// 數字 - 檢查炸彈
			if x > b.Info["Min"].(float64) && x < b.Info["Max"].(float64) {
				b.Info["Current"] = x
				g.recordPlayers()
				switch {
				case b.Info["Current"] == b.Info["Hit"]:
					b.gameOver(g)
					g.showRank()
					g.checkRank()
					g.reset()
					g.startPhase()
				case b.Info["Current"].(float64) < b.Info["Hit"].(float64):
					b.Info["Min"] = b.Info["Current"].(float64)
					g.show()
				case b.Info["Current"].(float64) > b.Info["Hit"].(float64):
					b.Info["Max"] = b.Info["Current"].(float64)
					g.show()
				}
				g.updateData()
			}
		}
	}
}

func (b *scene3InfoType) stage(g *GameType) {

}

func (b *scene3InfoType) show(g *GameType) string {
	return fmt.Sprintf("<%s%d> %d - %s - %d", 
			emoji.Emoji(":hourglass_not_done:"),
			int(b.Info["turn"].(float64)), 
			helper.Max(1, int(b.Info["Min"].(float64))), 
			emoji.Emoji(":bomb:"), 
			helper.Min(100, int(b.Info["Max"].(float64))))
}

func (b *scene3InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	b.Info["Turn"] = float64(5)
}

func (b *scene3InfoType) gameOver(g *GameType) {
	var str []string
	for _, u := range g.data.players.List {
		if u.UserID == users.LineUser.UserProfile.UserID {
			str = append(str, fmt.Sprintf("%s %s %d", u.DisplayName, emoji.Emoji(":umbrella:"), int(b.Info["Hit"].(float64))))
		} else {
			str = append(str, fmt.Sprintf("%s %s", u.DisplayName, emoji.Emoji(":collision:")))
			if _, exist := g.rank[u.UserID]; exist {
				g.rank[u.UserID].Boom++
			} else {
				g.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: 1}
			}
		}
	}
	texts = append(texts, strings.Join(str, "\n"))
}
