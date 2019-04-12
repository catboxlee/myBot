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

type scene3InfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene3InfoType) startPhase(g *GameType) {
	var p []string
	for _, v := range g.data.players.Queue {
		p = append(p, v.DisplayName)
	}
	texts = append(texts,
		fmt.Sprintf("[%s核爆危機]\n這是一顆核彈，請在%d回合內拆除...\n%s\n%s",
			emoji.Emoji(":radioactive:"),
			int(b.Info["Turn"].(float64)),
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
					b.Info["Turn"] = b.Info["Turn"].(float64) - float64(1)
					g.show()
					b.gameOver2(g)
				case b.Info["Current"].(float64) > b.Info["Hit"].(float64):
					b.Info["Max"] = b.Info["Current"].(float64)
					b.Info["Turn"] = b.Info["Turn"].(float64) - float64(1)
					g.show()
					b.gameOver2(g)
				}
				g.updateData()
			}
		}
	}
}

func (b *scene3InfoType) stage(g *GameType) {
	if _, exist := b.Info["Stage"]; exist {
		switch b.Info["Stage"] {
		default:
		}
	}
}

func (b *scene3InfoType) intoStage(g *GameType) {
	/*
		g.data.sceneInfo = &scene3AInfoType{}
		b.Info["Stage"] = "A"
		b.Info["Betrayal"] = users.LineUser.UserProfile.UserID
		g.data.sceneInfo.(*scene3AInfoType).Info = b.Info
		g.data.sceneInfo.(*scene3AInfoType).reset()
		g.startPhase()
	*/
}

func (b *scene3InfoType) show(g *GameType) string {
	return fmt.Sprintf("<%s%d>%d - %s - %d",
		emoji.Emoji(":hourglass_not_done:"),
		int(b.Info["Turn"].(float64)),
		helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":radioactive:"), helper.Min(100, int(b.Info["Max"].(float64))))
}

func (b *scene3InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	b.Info["Turn"] = float64(6)
	log.Println(b.Info)
}

func (b *scene3InfoType) gameOver(g *GameType) {
	texts = append(texts, fmt.Sprintf("%s 解除核彈 %s %d", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":red_heart:"), int(b.Info["Hit"].(float64))))
	if _, exist := g.rank[users.LineUser.UserProfile.UserID]; exist {
		g.rank[users.LineUser.UserProfile.UserID].Boom = helper.Max(0, g.rank[users.LineUser.UserProfile.UserID].Boom-1)
	}
}

func (b *scene3InfoType) gameOver2(g *GameType) {
	if b.Info["Turn"].(float64) > float64(0) {
		return
	}
	var str []string
	str = append(str, fmt.Sprintf("%s 核彈引爆", emoji.Emoji(":collision:")))
	for _, u := range g.data.players.List {
		str = append(str, fmt.Sprintf("%s %s", u.DisplayName, emoji.Emoji(":collision:")))
		if _, exist := g.rank[u.UserID]; exist {
			g.rank[u.UserID].Boom++
		} else {
			g.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: 1}
		}
	}
	texts = append(texts, strings.Join(str, "\n"))

	g.showRank()
	g.checkRank()
	g.reset()
	g.startPhase()
}
