package boomgame

import (
	"fmt"
	"log"
	"math"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
)

type scene4InfoType struct {
	Info map[string]interface{} `json:"info"`
}

func (b *scene4InfoType) startPhase(g *GameType) {
	texts = append(texts, fmt.Sprintf("[%s終極密碼2.0]\n%s", emoji.Emoji(":bomb:"), b.show(g)))
}

func (b *scene4InfoType) runPhase(input string, g *GameType) {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
			// 數字 - 檢查炸彈
			if x > b.Info["Min"].(float64) && x < b.Info["Max"].(float64) {
				b.Info["Current"] = x
				b.Info["LastPlayerID"] = b.Info["CurrentPlayerID"]
				b.Info["CurrentPlayerID"] = users.LineUser.UserProfile.UserID
				g.recordPlayers()
				switch {
				case b.Info["Current"] == b.Info["Hit"]:
					b.gameOver(g)
					// 結算
					users.LineUser.SaveUserData()
					g.showRank()
					g.checkRank()
					g.reset()
					g.startPhase()
				case b.Info["Current"].(float64) < b.Info["Hit"].(float64):
					b.Info["Min"] = b.Info["Current"].(float64)
					b.chkFate(g)
					g.show()
				case b.Info["Current"].(float64) > b.Info["Hit"].(float64):
					b.Info["Max"] = b.Info["Current"].(float64)
					b.chkFate(g)
					g.show()
				}
				g.updateData()
			}
		}
	}
}

func (b *scene4InfoType) stage(g *GameType) {
	if _, exist := b.Info["Stage"]; exist {
		switch b.Info["Stage"] {
		default:
		}
	}
}

func (b *scene4InfoType) intoStage(g *GameType) {
}

func (b *scene4InfoType) show(g *GameType) string {
	str := ""
	if b.Info["BoomCnt"].(float64) > float64(1) {
		if b.Info["BoomCnt"].(float64) > float64(99) {
			str = fmt.Sprintf("\n%s(%d)", emoji.Emoji(":bomb:"), int(b.Info["BoomCnt"].(float64)))
		} else {
			str = fmt.Sprintf("\n%s(%d)", strings.Repeat(emoji.Emoji(":bomb:"), int(b.Info["BoomCnt"].(float64))), int(b.Info["BoomCnt"].(float64)))
		}
	}
	return fmt.Sprintf("%d - %s - %d%s", helper.Max(1, int(b.Info["Min"].(float64))), emoji.Emoji(":bomb:"), helper.Min(100, int(b.Info["Max"].(float64))), str)
}

func (b *scene4InfoType) reset() {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	b.Info = make(map[string]interface{})
	b.Info["Hit"] = float64(boomDice.Hit)
	b.Info["Current"] = float64(0)
	b.Info["Min"] = float64(0)
	b.Info["Max"] = float64(101)
	b.Info["Turn"] = float64(0)
	b.Info["BoomCnt"] = float64(1)
	b.Info["CurrentPlayerID"] = ""
	b.Info["LastPlayerID"] = ""
	log.Println(b.Info)
}

func (b *scene4InfoType) gameOver(g *GameType) {
	str := b.chkChance(g)
	if len(str) > 0 {
		texts = append(texts, str)
	}

	if b.Info["BoomCnt"].(float64) > float64(99) {
		texts = append(texts, fmt.Sprintf("%s %s(%d)", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, emoji.Emoji(":collision:"), int(b.Info["BoomCnt"].(float64))))
	} else {
		texts = append(texts, fmt.Sprintf("%s %s(%d)", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, strings.Repeat(emoji.Emoji(":collision:"), int(b.Info["BoomCnt"].(float64))), int(b.Info["BoomCnt"].(float64))))
	}
	if _, exist := g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID]; exist {
		g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID].Boom += int(b.Info["BoomCnt"].(float64))
	} else {
		g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID] = &rankType{UserID: g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID, DisplayName: g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, Boom: int(b.Info["BoomCnt"].(float64))}
	}

}

func (b *scene4InfoType) chkChance(g *GameType) string {
	boomDice := &dice.Dice
	var strs string
	var swallowReturn = users.UsersList.Data[b.Info["CurrentPlayerID"].(string)].SwallowReturn + g.data.players.List[b.Info["CurrentPlayerID"].(string)].SwallowReturn
	if swallowReturn > 0 {
		if len(b.Info["LastPlayerID"].(string)) > 0 {
			boomDice.Roll("1d100")
			if swallowReturn >= int(boomDice.Hit) {
				strs += fmt.Sprintf("【%s】「燕返%d%%！」發動!\n", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, swallowReturn)
				tmp := g.data.players.List[b.Info["CurrentPlayerID"].(string)]
				if users.UsersList.Data[b.Info["CurrentPlayerID"].(string)].SwallowReturn <= 30 {
					users.UsersList.Data[b.Info["CurrentPlayerID"].(string)].SwallowReturn += 2
					strs += fmt.Sprintf("【%s】獲得 燕返(常駐)%d%%(+2%%)！\n", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, users.UsersList.Data[b.Info["CurrentPlayerID"].(string)].SwallowReturn)
				}
				tmp.SwallowReturn = 0
				g.data.players.List[b.Info["CurrentPlayerID"].(string)] = tmp
				b.Info["LastPlayerID"], b.Info["CurrentPlayerID"] = b.Info["CurrentPlayerID"], b.Info["LastPlayerID"]
				strs += fmt.Sprintf("%s%s", strs, b.chkChance(g))
			} else {
				strs += fmt.Sprintf("【%s】「燕返%d%%！」失敗.\n", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, swallowReturn)
				tmp := g.data.players.List[b.Info["CurrentPlayerID"].(string)]
				tmp.SwallowReturn = 0
				g.data.players.List[b.Info["CurrentPlayerID"].(string)] = tmp
			}
		}
	}

	boomDice.Roll("1d100")
	if 30 >= int(boomDice.Hit) {
		boomDice.Roll("1d1")
		switch int(boomDice.Hit) {
		case 1:
			strs += fmt.Sprintf("【%s】 瑪修「頌為堅城的雪花之壁！」\n", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName)
			b.Info["BoomCnt"] = math.Ceil(b.Info["BoomCnt"].(float64) / 3)
		}
	}
	return strs
}

func (b *scene4InfoType) chkFate(g *GameType) {
	boomDice := &dice.Dice
	boomDice.Roll("1d3")
	switch int(boomDice.Hit) {
	case 3:
		tmp := g.data.players.List[b.Info["CurrentPlayerID"].(string)]
		boomDice.Roll("1d10")
		tmp.SwallowReturn += boomDice.Hit
		g.data.players.List[b.Info["CurrentPlayerID"].(string)] = tmp
		texts = append(texts, fmt.Sprintf("【%s】獲得 燕返%d%%(+%d%%)", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, users.UsersList.Data[b.Info["CurrentPlayerID"].(string)].SwallowReturn+g.data.players.List[b.Info["CurrentPlayerID"].(string)].SwallowReturn, boomDice.Hit))
	}

	if 3 == int(b.Info["Max"].(float64)-b.Info["Min"].(float64)) {
		texts = append(texts, fmt.Sprintf("「信仰之躍！！！」"))
		b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) + 10
	} else {
		boomDice.Roll("1d3")
		switch int(boomDice.Hit) {
		case 3:
			texts = append(texts, fmt.Sprintf("御坂美琴「超電磁砲！」"))
			b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) + 5
		case 2:
			texts = append(texts, fmt.Sprintf("惠惠「Explosion！」"))
			b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) + 3
		default:
			texts = append(texts, fmt.Sprintf("漩渦鳴人「影分身之術！！」"))
			b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) + 2
		}
	}
}
