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
)

type scene0InfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (g *GameType) setSceneInfo() {
	switch g.scene {
	case 3:
		g.data.sceneInfo = &scene3InfoType{}
	case 2:
		g.data.sceneInfo = &scene2InfoType{}
	default:
		g.scene = 1
		g.data.sceneInfo = &scene0InfoType{}
		g.data.players.List = make(map[string]playerType)
		g.data.players.Queue = nil
	}
}
func (b *scene0InfoType) startPhase(g *GameType) {
	texts = append(texts, fmt.Sprintf("[%s終極密碼]\n%s", emoji.Emoji(":bomb:"), b.show(g)))
}

func (b *scene0InfoType) runPhase(input string, g *GameType) {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
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
						b.intoStage(g)
					case b.Info["Current"].(float64) > b.Info["Hit"].(float64):
						b.Info["Max"] = b.Info["Current"].(float64)
						g.show()
						b.intoStage(g)
					}
					g.updateData()
				}
			}
		}
	}
}

func (b *scene0InfoType) stage(g *GameType) {
	if _, exist := b.Info["Stage"]; exist {
		switch b.Info["Stage"] {
		case "A":
			g.data.sceneInfo = &scene0AInfoType{}
			g.data.sceneInfo.(*scene0AInfoType).Info = b.Info
			return
		default:
		}
	}
}

func (b *scene0InfoType) intoStage(g *GameType) {
	b.Info["Turn"] = b.Info["Turn"].(float64) + float64(1)

	if b.Info["Turn"].(float64) > float64(3) {
		return
	}
	boomDice := &dice.Dice
	boomDice.Roll("1d6")

	if boomDice.Hit == 6 {
		boomDice := &dice.Dice
		boomDice.Roll("1d1")

		switch boomDice.Hit {
		case 1:
			g.data.sceneInfo = &scene0AInfoType{}
			b.Info["Stage"] = "A"
			b.Info["Betrayal"] = users.LineUser.UserProfile.UserID
			g.data.sceneInfo.(*scene0AInfoType).Info = b.Info
			g.startPhase()
		}
	}
}

func (b *scene0InfoType) show(g *GameType) string {
	return fmt.Sprintf("%d - %s - %d", helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":bomb:"), helper.Min(100, int(b.Info["Max"].(float64))))
}

func (b *scene0InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	b.Info["Turn"] = float64(0)
	log.Println(b.Info)
}

func (b *scene0InfoType) gameOver(g *GameType) {
	texts = append(texts, fmt.Sprintf("%s %s %d", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":collision:"), int(b.Info["Hit"].(float64))))
	if _, exist := g.rank[users.LineUser.UserProfile.UserID]; exist {
		g.rank[users.LineUser.UserProfile.UserID].Boom++
	} else {
		g.rank[users.LineUser.UserProfile.UserID] = &rankType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName, Boom: 1}
	}
}
