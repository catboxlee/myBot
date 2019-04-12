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

type scene2InfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene2InfoType) startPhase(g *GameType) {
	var p []string
	for _, v := range g.data.players.Queue {
		p = append(p, v.DisplayName)
	}
	texts = append(texts,
		fmt.Sprintf("[%s%s%s天外奇蹟]\n你們在3萬5千英呎的高空上，房子即將墜落，幸運的是這裡還藏有一件降落傘...\n%s\n%s",
			emoji.Emoji(":balloon:"),
			emoji.Emoji(":house:"),
			emoji.Emoji(":balloon:"),
			strings.Join(p, ", "),
			b.show(g)))
}

func (b *scene2InfoType) runPhase(input string, g *GameType) {
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
					boomDice := &dice.Dice
					boomDice.Roll("1d6")

					if boomDice.Hit == 6 {
						b.intoStage(g)
						g.updateData()
						return
					}
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

func (b *scene2InfoType) stage(g *GameType) {
	if _, exist := b.Info["Stage"]; exist {
		switch b.Info["Stage"] {
		case "A":
			g.data.sceneInfo = &scene2AInfoType{}
			g.data.sceneInfo.(*scene2AInfoType).Info = b.Info
			return
		default:
		}
	}
}

func (b *scene2InfoType) intoStage(g *GameType) {
	g.data.sceneInfo = &scene2AInfoType{}
	b.Info["Stage"] = "A"
	b.Info["Betrayal"] = users.LineUser.UserProfile.UserID
	g.data.sceneInfo.(*scene2AInfoType).Info = b.Info
	g.data.sceneInfo.(*scene2AInfoType).reset()
	g.startPhase()
}

func (b *scene2InfoType) show(g *GameType) string {
	return fmt.Sprintf("%d - %s - %d", helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":closed_umbrella:"), helper.Min(100, int(b.Info["Max"].(float64))))
}

func (b *scene2InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	log.Println(b.Info)
}

func (b *scene2InfoType) gameOver(g *GameType) {
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
