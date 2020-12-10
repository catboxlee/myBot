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
		if len(matches) > 1 {
			if x, err := strconv.ParseFloat(matches[1], 64); err == nil {
				// 數字 - 檢查炸彈
				if x > b.Info["Min"].(float64) && x < b.Info["Max"].(float64) {
					b.Info["Current"] = x
					g.recordPlayers()
					switch {
					case b.Info["Current"] == b.Info["Hit"]:
						b.Info["CurrentPlayerID"] = users.LineUser.UserProfile.UserID
						b.gameOver(g)
						g.showRank()
						g.checkRank()
						g.reset()
						g.startPhase()
					case b.Info["Current"].(float64) < b.Info["Hit"].(float64):
						b.Info["Min"] = b.Info["Current"].(float64)
						b.Info["CurrentPlayerID"] = users.LineUser.UserProfile.UserID
						b.Info["LastPlayerID"] = users.LineUser.UserProfile.UserID
						b.chkFate(g)
						g.show()
						b.intoStage(g)
					case b.Info["Current"].(float64) > b.Info["Hit"].(float64):
						b.Info["Max"] = b.Info["Current"].(float64)
						b.Info["CurrentPlayerID"] = users.LineUser.UserProfile.UserID
						b.Info["LastPlayerID"] = users.LineUser.UserProfile.UserID
						b.chkFate(g)
						g.show()
						b.intoStage(g)
					}
					g.updateData()
				}
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
	b.Info["LastPlayerID"] = string("")
	b.Info["CurrentPlayerID"] = string("")
	log.Println(b.Info)
}

func (b *scene4InfoType) gameOver(g *GameType) {

	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	lucky := boomDice.Hit
	if 10 >= lucky {
		var str []string
		reBooms := 0
		boomDice := &dice.Dice
		str = append(str, fmt.Sprintf("%s 翻開陷阱卡「神聖慧星。反射力量」", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName))
		for _, u := range g.data.players.List {
			if u.UserID != b.Info["CurrentPlayerID"].(string) {
				boomDice.Roll("1d3")
				switch int(boomDice.Hit) {
				case 2:
					str = append(str, fmt.Sprintf("%s 「雪花之壁！」 %s(%d)", u.DisplayName, emoji.Emoji(":collision:"), int(math.Ceil(b.Info["BoomCnt"].(float64)/3))))
					if _, exist := g.rank[u.UserID]; exist {
						g.rank[u.UserID].Boom += int(math.Ceil(b.Info["BoomCnt"].(float64) / 3))
					} else {
						g.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: int(math.Ceil(b.Info["BoomCnt"].(float64) / 3))}
					}
				case 1:
					str = append(str, fmt.Sprintf("%s 「燕返！」", u.DisplayName))
					reBooms += int(b.Info["BoomCnt"].(float64))
				default:
					str = append(str, fmt.Sprintf("%s %s(%d)", u.DisplayName, emoji.Emoji(":collision:"), int(b.Info["BoomCnt"].(float64))))
					if _, exist := g.rank[u.UserID]; exist {
						g.rank[u.UserID].Boom += int(b.Info["BoomCnt"].(float64))
					} else {
						g.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: int(b.Info["BoomCnt"].(float64))}
					}
				}
			}
		}
		if reBooms > 0 {
			boomDice.Roll("1d3")
			if int(boomDice.Hit) == 1 {
				str = append(str, fmt.Sprintf("%s 「雪花之壁！」 %s(%d)", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, emoji.Emoji(":collision:"), int(math.Ceil(float64(reBooms)/3))))
				reBooms = int(math.Ceil(float64(reBooms) / 3))
			} else {
				str = append(str, fmt.Sprintf("%s %s(%d)", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, emoji.Emoji(":collision:"), int(reBooms)))
			}
			if _, exist := g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID]; exist {
				g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID].Boom += reBooms
			} else {
				g.rank[g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID] = &rankType{UserID: g.data.players.List[b.Info["CurrentPlayerID"].(string)].UserID, DisplayName: g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName, Boom: reBooms}
			}
		}
		texts = append(texts, strings.Join(str, "\n"))
	} else if 30 >= lucky {

		texts = append(texts, b.chkChance(g))

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
	} else {
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

}

func (b *scene4InfoType) chkChance(g *GameType) string {
	boomDice := &dice.Dice
	var strs string
	boomDice.Roll("1d2")
	switch int(boomDice.Hit) {
	case 2:
		if len(b.Info["LastPlayerID"].(string)) > 0 {
			strs = fmt.Sprintf("【%s】 不二周助「燕返！」", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName)
			b.Info["LastPlayerID"], b.Info["CurrentPlayerID"] = b.Info["CurrentPlayerID"], b.Info["LastPlayerID"]
			boomDice.Roll("1d100")
			if 30 >= int(boomDice.Hit) {
				strs = fmt.Sprintf("%s\n%s", strs, b.chkChance(g))
			}
		}
	default:
		strs = fmt.Sprintf("【%s】 Shielder瑪修「頌為堅城的雪花之壁！」", g.data.players.List[b.Info["CurrentPlayerID"].(string)].DisplayName)
		b.Info["BoomCnt"] = math.Ceil(b.Info["BoomCnt"].(float64) / 3)
	}
	return strs
}

func (b *scene4InfoType) chkFate(g *GameType) {
	boomDice := &dice.Dice
	boomDice.Roll("1d100")
	lucky := boomDice.Hit
	if 100 >= lucky {
		if 3 == int(b.Info["Max"].(float64)-b.Info["Min"].(float64)) {
			texts = append(texts, fmt.Sprintf("「信仰之躍！！！」"))
			b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) * 3
		} else {
			boomDice.Roll("1d3")
			switch int(boomDice.Hit) {
			case 3:
				texts = append(texts, fmt.Sprintf("御坂美琴「超電磁砲！」"))
				b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) * 5
			case 2:
				texts = append(texts, fmt.Sprintf("惠惠「Explosion！」"))
				b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) * 3
			default:
				texts = append(texts, fmt.Sprintf("漩渦鳴人「影分身之術！！」"))
				b.Info["BoomCnt"] = (b.Info["BoomCnt"].(float64)) * 2
			}
		}
	}
}
