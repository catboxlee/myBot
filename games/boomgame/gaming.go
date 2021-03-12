package boomgame

import (
	"fmt"
	"log"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/cards"
	"myBot/games/boomgame/title"
	"myBot/helper"
	"myBot/users"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func (g *GameType) resetInfo() {
	inning := 0
	if g.Info != nil {
		inning = g.Info.Inning
	}
	boomHit := rand.Intn(100) + 1
	g.Info = &InfoType{
		Hit:       boomHit,
		Min:       0,
		Max:       101,
		Turn:      0,
		Inning:    inning,
		BoomCnt:   1,
		Queue:     []string{},
		PlayQueue: []string{},
	}
	log.Println("Info", g.Info)
}

func (g *GameType) showGameInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%d - %s - %d", helper.Max(1, g.Info.Min), emoji.Emoji(":bomb:"), helper.Min(100, g.Info.Max)))
	if g.Info.BoomCnt > 99 {
		strs = append(strs, fmt.Sprintf("%s(%d)", emoji.Emoji(":bomb:"), (g.Info.BoomCnt)))
	} else {
		strs = append(strs, fmt.Sprintf("%s(%d)", strings.Repeat(emoji.Emoji(":bomb:"), g.Info.BoomCnt), g.Info.BoomCnt))
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) runPhase(input string) {
	var strs []string
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		if x, err := strconv.Atoi(matches[1]); err == nil {
			// 數字 - 檢查炸彈
			if x > g.Info.Min && x < g.Info.Max {
				g.OnPlay()
				g.Info.Turn++
				g.Info.Current = x
				g.doCoolDown(g.Info.CurrentUserID)
				if s := g.mythosCardTrigger(); len(s) > 0 {
					strs = append(strs, s)
				}
				log.Println(g.GetHit())
				if s := g.senteCardTrigger(); len(s) > 0 {
					strs = append(strs, s)
				}
				switch {
				case g.Info.Current == g.Info.Hit:
					if s := g.onHit(); len(s) > 0 {
						strs = append(strs, s)
					}
					if len(strs) > 0 {
						texts = append(texts, strings.Join(strs, "\n"))
					}
					g.showRank()
					// 結算
					g.checkRank()
					g.updateData()
					return
				case g.Info.Current < g.Info.Hit:
					g.Info.Min = helper.Max(g.Info.Current, g.Info.Min)
					if s := g.passCardTrigger(); len(s) > 0 {
						strs = append(strs, s)
					}
					if s := g.mythosPassCardTrigger(); len(s) > 0 {
						strs = append(strs, s)
					}
				case g.Info.Current > g.Info.Hit:
					g.Info.Max = helper.Min(g.Info.Current, g.Info.Max)
					if s := g.passCardTrigger(); len(s) > 0 {
						strs = append(strs, s)
					}
					if s := g.mythosPassCardTrigger(); len(s) > 0 {
						strs = append(strs, s)
					}
				}
				if len(strs) > 0 {
					texts = append(texts, strings.Join(strs, "\n"))
				}
				g.Show()
				g.updateData()
			}
		}
	}
}

func (g *GameType) onHit() string {
	var strs []string
	boomerID := g.Info.PlayQueue[len(g.Info.PlayQueue)-1]
	var boomerIDs []string
	// 引爆前: 反擊, 回擊, 燕返, 時停, 回復...
	if s := g.hitCardTrigger(&boomerID); len(s) > 0 {
		strs = append(strs, s)
	}

	boomerID = g.Info.PlayQueue[len(g.Info.PlayQueue)-1]
	boomerIDs = append(boomerIDs, boomerID)

	// 引爆: 盾, 回復, 連鎖, 無傷, 鎖血...
	if s := g.shieldCardTrigger(&boomerID); len(s) > 0 {
		strs = append(strs, s)
	}

	// 結算
	if g.Info.BoomCnt > int(99) {
		strs = append(strs, fmt.Sprintf("%s %s(%d)", users.UsersList.Data[boomerID].GetDisplayName(), emoji.Emoji(":collision:"), g.Info.BoomCnt))
	} else {
		strs = append(strs, fmt.Sprintf("%s %s(%d)", users.UsersList.Data[boomerID].GetDisplayName(), strings.Repeat(emoji.Emoji(":collision:"), helper.Max(1, g.Info.BoomCnt)), g.Info.BoomCnt))
	}

	// 結算獎金
	strs = append(strs, fmt.Sprintf("[[單局獎勵]]"))
	for _, val := range g.Info.Queue {
		if exist, _ := helper.InArray(val, boomerIDs); exist {
			g.rank[val].Boom += g.Info.BoomCnt
			g.rank[val].PlayTimes++
			if s := g.checkAchieve(&boomerID); len(s) > 0 {
				strs = append(strs, s)
			}
			//strs = append(strs, fmt.Sprintf("%s %s%d(%+d)", UserID, emoji.Emoji(":collision:"), g.playerList[UserID].Boom, g.Info.BoomCnt))
		} else {
			users.UsersList.Data[val].MakeGemStone(singleGameBonusGem)
			strs = append(strs, fmt.Sprintf("%s %s%d(%+d)", users.UsersList.Data[val].GetDisplayName(), emoji.Emoji(":gem_stone:"), users.UsersList.Data[val].GetGemStone(), singleGameBonusGem))

			g.rank[val].PlayTimes++
			g.rank[val].WinTimes++
		}
		users.UsersList.SaveData(val)
		g.Players.Data[val].SaveData()
	}
	g.Info.Inning++
	return strings.Join(strs, "\n")

}
func (g *GameType) checkAchieve(boomerID *string) string {
	var strs []string
	thisPlayer := g.Player(*boomerID)
	// 一拳
	if len(g.Info.PlayQueue) > 1 && g.Info.Turn <= len(g.Info.Queue) && g.Info.Queue[g.Info.Turn-1] != g.Info.PlayQueue[len(g.Info.PlayQueue)-1] {
		if exist, _ := helper.InArray("one shot", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "one shot")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("one shot")))
		}
		thisPlayer.Property.OneShot++
		thisPlayer.TakeCard("19")
		strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["19"].GetDisplayName()))
	}
	// 狙擊王
	if len(g.Info.Queue) == 1 && g.Info.Turn == 1 {
		if exist, _ := helper.InArray("sogeking", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "sogeking")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("sogeking")))
		}
		thisPlayer.TakeCard("1")
		strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["1"].GetDisplayName()))
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) checkSeasonAchieve(boomerID string) string {
	var strs []string
	thisPlayer := g.Player(boomerID)
	// 三連冠
	if thisPlayer.Property.LosingStreak >= 3 {
		if exist, _ := helper.InArray("three consecutive", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "three consecutive")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("three consecutive")))
			thisPlayer.TakeCard("12")
			strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["12"].GetDisplayName()))
		}
	}
	// 五連冠
	if thisPlayer.Property.LosingStreak >= 5 {
		if exist, _ := helper.InArray("five consecutive", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "five consecutive")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("five consecutive")))
			thisPlayer.TakeCard("12")
			strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["12"].GetDisplayName()))
		}
	}
	// 七連冠
	if thisPlayer.Property.LosingStreak >= 7 {
		if exist, _ := helper.InArray("seven consecutive", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "seven consecutive")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("seven consecutive")))
			thisPlayer.TakeCard("12")
			strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["12"].GetDisplayName()))
		}
	}
	// 十連冠
	if thisPlayer.Property.LosingStreak >= 10 {
		if exist, _ := helper.InArray("ten consecutive", thisPlayer.Titles); !exist {
			thisPlayer.Titles = append(thisPlayer.Titles, "ten consecutive")
			strs = append(strs, fmt.Sprintf("%s 獲得稱號<%s>", thisPlayer.GetDisplayName(), title.Title("ten consecutive")))
			thisPlayer.TakeCard("12")
			strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["12"].GetDisplayName()))
			thisPlayer.TakeCard("12")
			strs = append(strs, fmt.Sprintf("%s 獲得卡片<%s>", thisPlayer.GetDisplayName(), thisPlayer.CardPile.Cards["12"].GetDisplayName()))
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) hitCardTrigger(boomerID *string) string {
	var strs []string

	for _, co := range g.Player(*boomerID).CardPile.Cards {
		if co.OnHitFunc != nil {
			if r, s := co.OnHitFunc(); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
				break
			}
		}
	}
	if *boomerID != g.Info.PlayQueue[len(g.Info.PlayQueue)-1] {
		*boomerID = g.Info.PlayQueue[len(g.Info.PlayQueue)-1]
		if s := g.hitCardTrigger(boomerID); len(s) > 0 {
			strs = append(strs, s)
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) senteCardTrigger() string {
	var strs []string

	for _, co := range g.Player(g.Info.CurrentUserID).CardPile.Cards {
		if co.OnSenteFunc != nil {
			if r, s := co.OnSenteFunc(); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
				break
			}
		}
	}
	return strings.Join(strs, "\n")
}
func (g *GameType) passCardTrigger() string {
	var strs []string

	for _, co := range g.Player(g.Info.CurrentUserID).CardPile.Cards {
		if co.OnPassFunc != nil {
			if r, s := co.OnPassFunc(); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
				break
			}
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) mythosCardTrigger() string {
	var strs []string

	for _, co := range g.mythosCards {
		if co.OnMythosFunc != nil {
			if r, s := co.OnMythosFunc(g); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
				break
			}
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) mythosPassCardTrigger() string {
	var strs []string

	for _, co := range g.mythosCards {
		if co.OnMythosPassFunc != nil {
			if r, s := co.OnMythosPassFunc(g); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
				break
			}
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) shieldCardTrigger(boomerID *string) string {
	var strs []string
	var cos []*cards.CardOption
	for _, uid := range g.GetQueue() {
		for _, co := range g.Player(uid).CardPile.Cards {
			if uid == *boomerID {
				if co.OnShieldFunc != nil {
					cos = append(cos, co)
				}
			} else {
				if co.OnAttackFunc != nil {
					cos = append(cos, co)
				}
			}
		}
	}
	if len(cos) > 0 {
		tmp := rand.Perm(len(cos))
		for _, id := range tmp {
			co := cos[tmp[id]]
			if co.OnAttackFunc != nil {
				_, s := co.OnAttackFunc()
				if len(s) > 0 {
					strs = append(strs, s)
				}
			}
			if co.OnShieldFunc != nil {
				_, s := co.OnShieldFunc()
				if len(s) > 0 {
					strs = append(strs, s)
				}
			}
		}
	}
	return strings.Join(strs, "\n")
}
func (g *GameType) newSeason() {
	g.season++
	g.rank = make(map[string]*rankType)
	g.Info.Inning = 0
	g.resetInfo()
	g.Show()
}

func (g *GameType) newRound() {
	g.Info.Inning++
	g.resetInfo()
	g.Show()
}

func (g *GameType) showRank() {
	var strs []string
	values := make([]*rankType, 0, len(g.rank))

	for _, v := range g.rank {
		values = append(values, v)
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Boom > values[j].Boom
	})

	strs = append(strs, fmt.Sprintf("S%d Rank：", g.season))

	for _, val := range values {
		strs = append(strs, fmt.Sprintf("%s %s%d", users.UsersList.Data[val.UserID].GetDisplayName(), emoji.Emoji(":collision:"), val.Boom))
	}

	if len(strs) > 0 {
		texts = append(texts, strings.Join(strs, "\n"))
	}
}

func (g *GameType) checkRank() {
	var strs []string
	var isBoomKing bool = false
	var boomKings []*rankType
	var otherPlayer []*rankType

	values := make([]*rankType, 0, len(g.rank))

	for _, v := range g.rank {
		values = append(values, v)
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Boom > values[j].Boom
	})

	for _, val := range values {
		if val.Boom >= seasonBoomCount {
			isBoomKing = true
			boomKings = append(boomKings, val)
		} else {
			otherPlayer = append(otherPlayer, val)
		}
	}

	if isBoomKing {
		strs = append(strs, fmt.Sprintf("%s S%d 爆爆王 %s", emoji.Emoji(":confetti_ball:"), g.season, emoji.Emoji(":confetti_ball:")))
		// King
		for _, val := range boomKings {
			strs = append(strs, fmt.Sprintf("%s %s%d", users.UsersList.Data[val.UserID].GetDisplayName(), emoji.Emoji(":collision:"), val.Boom))
			g.Player(val.UserID).Property.LosingStreak++
			g.Player(val.UserID).Property.WinningStreak = 0
			if s := g.checkSeasonAchieve(val.UserID); len(s) > 0 {
				strs = append(strs, s)
			}
			g.Player(val.UserID).SaveData()
			users.UsersList.SaveData(val.UserID)
		}
		// bonus
		strs = append(strs, fmt.Sprintf("[[賽季結算獎勵]]"))
		for _, val := range otherPlayer {
			g.Player(val.UserID).Property.WinningStreak++
			g.Player(val.UserID).Property.LosingStreak = 0
			if s := g.checkSeasonAchieve(val.UserID); len(s) > 0 {
				strs = append(strs, s)
			}
			users.UsersList.Data[val.UserID].MakeGemStone(seasonGameBonusGem)
			strs = append(strs, fmt.Sprintf("%s %s%d(%+d)", users.UsersList.Data[val.UserID].GetDisplayName(), emoji.Emoji(":gem_stone:"), users.UsersList.Data[val.UserID].GetGemStone(), seasonGameBonusGem))
			users.UsersList.SaveData(val.UserID)
			g.Player(val.UserID).SaveData()
		}
		if len(strs) > 0 {
			texts = append(texts, strings.Join(strs, "\n"))
		}
		g.newSeason()
	} else {
		g.newRound()
	}
}
