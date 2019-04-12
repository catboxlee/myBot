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

type scene0AInfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene0AInfoType) startPhase(g *GameType) {
	var p []string
	for _, v := range g.data.players.Queue {
		p = append(p, v.DisplayName)
	}
	texts = append(texts,
		fmt.Sprintf("[%s炸彈狂魔 %s]\n阻止他獲得炸彈!\n%s",
			g.data.players.List[b.Info["Betrayal"].(string)].DisplayName,
			emoji.Emoji(":smiling_face_with_horns:"),
			b.show(g)))
}

func (b *scene0AInfoType) runPhase(input string, g *GameType) {
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
				b.Info["Betrayal"] = users.LineUser.UserProfile.UserID
				g.updateData()
			}
		}
	}
}

func (b *scene0AInfoType) show(g *GameType) string {
	return fmt.Sprintf("%d - %s - %d", helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":smiling_face_with_horns:"), helper.Min(100, int(b.Info["Max"].(float64))))
}

func (b *scene0AInfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	//b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	log.Println(b.Info)
	//b.info = nil
}

func (b *scene0AInfoType) gameOver(g *GameType) {
	var str []string
	betrayal := g.data.players.List[b.Info["Betrayal"].(string)]
	if users.LineUser.UserProfile.UserID == b.Info["Betrayal"].(string) {
		for _, u := range g.data.players.List {
			if u.UserID == users.LineUser.UserProfile.UserID {
				str = append(str, fmt.Sprintf("%s %s %d", u.DisplayName, emoji.Emoji(":smiling_face_with_sunglasses:"), int(b.Info["Hit"].(float64))))
			} else {
				str = append(str, fmt.Sprintf("%s %s", u.DisplayName, emoji.Emoji(":collision:")))
				if _, exist := g.rank[u.UserID]; exist {
					g.rank[u.UserID].Boom++
				} else {
					g.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: 1}
				}
			}
		}
	} else {
		str = append(str, fmt.Sprintf("%s %s %d", betrayal.DisplayName, emoji.Emoji(":collision:"), int(b.Info["Hit"].(float64))))
		if _, exist := g.rank[betrayal.UserID]; exist {
			g.rank[betrayal.UserID].Boom++
		} else {
			g.rank[betrayal.UserID] = &rankType{UserID: betrayal.UserID, DisplayName: betrayal.DisplayName, Boom: 1}
		}
	}
	texts = append(texts, strings.Join(str, "\n"))
}

func (b *scene0AInfoType) stage(g *GameType) {

}
