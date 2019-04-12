package boomgame

import (
	"myBot/dice"
	"myBot/emoji"
	"myBot/users"
	"fmt"
	"log"
	"strings"
)

type scene4InfoType struct {
	Turn        int                          `json:"turn"`
	Last        int                          `json:"last"`
	TargetPoint int                          `json:"targetpoint"`
	Players     map[string]*scene4PlayerType `json:"Players"`
	Info        map[string]interface{}       `json:"info"`
}

type scene4PlayerType struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"displayname"`
	Footprints  int    `json:"footprints"`
	Turn        int    `json:"turn"`
	TmpDice     struct {
		Footprints int `json:"footprints"`
		Booms      int `json:"booms"`
		Damages    int `json:"damages"`
	} `json:"tmpdice"`
}

type diceFaces struct {
	emoji string
	value string
}

// ...
var (
	DAMAGE     = diceFaces{emoji.Emoji(":sleeping_face:"), "damage"}
	FOOTPRINTS = diceFaces{emoji.Emoji(":footprints:"), "footprints"}
	BOOM       = diceFaces{emoji.Emoji(":zombie:"), "boom"}
)

var diceFace = [6]diceFaces{DAMAGE, BOOM, FOOTPRINTS, FOOTPRINTS, BOOM, DAMAGE}

func (b *scene4InfoType) startPhase(g *GameType) {
	
	texts = append(texts,
		fmt.Sprintf("[%s屍速列車%s]\n逃往第%d節車廂...\n[+]擲骰 [-]結束回合",
			emoji.Emoji(":bullet_train:"),
			emoji.Emoji(":zombie:"),
			b.TargetPoint))
}

func (b *scene4InfoType) runPhase(input string, g *GameType) {
	var text string
	if strings.HasPrefix(input, "+") {
		if b.checkPlayerBoom(g) {
			return
		}
		//c := int(math.Floor(float64(int(b.Info["Turn"].(float64)) % 2)))
		g.recordPlayers()
		pkDices := dice.Dice
		pkDices.Roll("3d6")
		for _, v := range pkDices.Rolls {
			//fmt.Printf(" %v", diceFace[v-1])
			switch diceFace[v-1] {
			case BOOM:
				b.Players[users.LineUser.UserProfile.UserID].TmpDice.Booms++
			case FOOTPRINTS:
				b.Players[users.LineUser.UserProfile.UserID].TmpDice.Footprints++
			case DAMAGE:
				b.Players[users.LineUser.UserProfile.UserID].TmpDice.Damages++
			}
		}

		text += fmt.Sprintf("%s: %s %s %s", b.Players[users.LineUser.UserProfile.UserID].DisplayName, diceFace[pkDices.Rolls[0]-1].emoji, diceFace[pkDices.Rolls[1]-1].emoji, diceFace[pkDices.Rolls[2]-1].emoji)
		text += fmt.Sprintf("\n%s%d(+%d)/%s%d", FOOTPRINTS.emoji, b.Players[users.LineUser.UserProfile.UserID].Footprints+b.Players[users.LineUser.UserProfile.UserID].TmpDice.Footprints, b.Players[users.LineUser.UserProfile.UserID].TmpDice.Footprints, BOOM.emoji, b.Players[users.LineUser.UserProfile.UserID].TmpDice.Booms)

		if b.Players[users.LineUser.UserProfile.UserID].TmpDice.Booms < 3 {
			texts = append(texts, text)
		} else {
			text += b.endTurn(g)
			texts = append(texts, text)
			g.show()
			b.gameOver(g)
		}
		g.updateData()
	} else if strings.HasPrefix(input, "-") {
		if b.checkPlayerBoom(g) {
			return
		}
		g.recordPlayers()
		text := b.endTurn(g)
		if len(text) > 0 {
			texts = append(texts, text)
		}
		g.show()
		b.gameOver(g)
		g.updateData()
	}
	return
}

func (b *scene4InfoType) checkPlayerBoom(g *GameType) bool {
	if _, exist := b.Players[users.LineUser.UserProfile.UserID]; exist {
		if b.Last > 0 && b.Players[users.LineUser.UserProfile.UserID].Turn < b.Turn {
			return true
		}
	} else {
		b.Players[users.LineUser.UserProfile.UserID] = &scene4PlayerType{}
		b.Players[users.LineUser.UserProfile.UserID].UserID = users.LineUser.UserProfile.UserID
		b.Players[users.LineUser.UserProfile.UserID].DisplayName = users.LineUser.UserProfile.DisplayName
	}
	return false
}

func (b *scene4InfoType) running(g *GameType) {

}
func (b *scene4InfoType) resting(g *GameType) {

}

func (b *scene4InfoType) stage(g *GameType) {
	if _, exist := b.Info["Stage"]; exist {
		switch b.Info["Stage"] {
		default:
		}
	}
}

func (b *scene4InfoType) intoStage(g *GameType) {
	/*
		g.data.sceneInfo = &scene3AInfoType{}
		b.Info["Stage"] = "A"
		b.Info["Betrayal"] = users.LineUser.UserProfile.UserID
		g.data.sceneInfo.(*scene3AInfoType).Info = b.Info
		g.data.sceneInfo.(*scene3AInfoType).reset()
		g.startPhase()
	*/
}

func (b *scene4InfoType) show(g *GameType) string {
	var str []string
	for _, v := range g.data.players.Queue {
		if b.Players[v.UserID].Footprints >= b.TargetPoint {
			str = append(str, fmt.Sprintf("%s %s %d", b.Players[v.UserID].DisplayName, emoji.Emoji(":smiling_face_with_smiling_eyes:"), b.Players[v.UserID].Footprints))
		} else {
			if b.Last > 0 {
				if b.Players[v.UserID].Turn < b.Last {
					str = append(str, fmt.Sprintf("%s %s %d", b.Players[v.UserID].DisplayName, emoji.Emoji(":zombie:"), b.Players[v.UserID].Footprints))
				} else {
					str = append(str, fmt.Sprintf("%s %s %d", b.Players[v.UserID].DisplayName, emoji.Emoji(":face_screaming_in_fear:"), b.Players[v.UserID].Footprints))
				}
			} else {
				str = append(str, fmt.Sprintf("%s %s %d", b.Players[v.UserID].DisplayName, FOOTPRINTS.emoji, b.Players[v.UserID].Footprints))
			}
		}
	}
	return strings.Join(str, "\n")
}

func (b *scene4InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d2")
	b.Turn = 1
	b.Last = 0
	b.TargetPoint = 13 + boomDice.Hit
	b.Info = make(map[string]interface{})
	b.Players = make(map[string]*scene4PlayerType)
	log.Println(b.Info)
}

func (b *scene4InfoType) gameOver(g *GameType) {
	if b.Last == 0 {
		return
	}
	if users.LineUser.UserProfile.UserID != g.data.players.Queue[len(g.data.players.Queue)-1].UserID {
		return
	}
	for _, v := range g.data.players.List {
		if b.Players[v.UserID].Turn >= b.Last && b.Players[v.UserID].Footprints < b.TargetPoint {
			if _, exist := g.rank[v.UserID]; exist {
				g.rank[v.UserID].Boom++
			} else {
				g.rank[v.UserID] = &rankType{UserID: v.UserID, DisplayName: v.DisplayName, Boom: 1}
			}
		}
	}

	g.checkRank()
	g.reset()
	g.startPhase()
}

func (b *scene4InfoType) endTurn(g *GameType) string {
	var text string
	p := b.Players[users.LineUser.UserProfile.UserID]
	if p.TmpDice.Booms < 3 {
		if p.TmpDice.Footprints > 0 {
			p.Footprints += p.TmpDice.Footprints
		} else {
			text += fmt.Sprintf("\n%s 什麼事都沒做", p.DisplayName)
		}
	} else {
		text += fmt.Sprintf("\n%s 逃跑失敗", p.DisplayName)
	}
	b.Turn++
	p.TmpDice.Damages = 0
	p.TmpDice.Booms = 0
	p.TmpDice.Footprints = 0
	if b.Last == 0 && p.Footprints >= b.TargetPoint {
		b.Last = p.Turn
	} else {
		p.Turn++
	}

	return text
}
