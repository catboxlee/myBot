package boomgame

import (
	"fmt"
	"log"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
)

type scene2AInfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene2AInfoType) startPhase(g *GameType) {
	var p []string
	for _, v := range g.data.players.Queue {
		p = append(p, v.DisplayName)
	}
	texts = append(texts,
		fmt.Sprintf("[%s%s%s天外奇蹟 II]\n%s獲得降落傘，請尋找逃生出口...\n%s\n%s",
			emoji.Emoji(":balloon:"),
			emoji.Emoji(":house:"),
			emoji.Emoji(":balloon:"),
			g.data.players.List[b.Info["Betrayal"].(string)].DisplayName,
			strings.Join(p, ", "),
			b.show(g)))
}

func (b *scene2AInfoType) runPhase(input string, g *GameType) {
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

func (b *scene2AInfoType) stage(g *GameType) {

}

func (b *scene2AInfoType) show(g *GameType) string {
	return fmt.Sprintf("%d - %s - %d", helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":door:"), helper.Min(int(b.Info["MaxLimit"].(float64))-1, int(b.Info["Max"].(float64))))
}

func (b *scene2AInfoType) reset() {
	if _, exist := b.Info["MaxLimit"]; exist {
		if b.Info["MaxLimit"].(float64)-float64(10) < float64(11) {
			b.Info["MaxLimit"] = float64(helper.Max(2, int(b.Info["MaxLimit"].(float64))-1))
		} else {
			b.Info["MaxLimit"] = float64(helper.Max(11, int(b.Info["MaxLimit"].(float64))-10))
		}
	} else {
		b.Info["MaxLimit"] = float64(101)
	}
	boomDice := &dice.Dice
	boomDice.Roll("1d" + strconv.Itoa(int(b.Info["MaxLimit"].(float64))-1))
	//b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = b.Info["MaxLimit"]
	log.Println(b.Info)
}

func (b *scene2AInfoType) gameOver(g *GameType) {
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
