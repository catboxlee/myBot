package data

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/scheduler"
	"strings"
)

// MythosCard ...
var MythosCard = map[string]CardOption{
	/*
			"killer_queen": CardOption{
				CardName:    fmt.Sprintf("%s「皇后殺手 吉良吉影」", emoji.Emoji(":ghost:")),
				DisplayName: "皇后殺手 吉良吉影",
				Class:       "SR",
				CoreSet:     "killer_queen",
				CoolDown:    9,
				ReCoolDown:  9,
				Unique:      true,
				DescFunc: func(thisCard scheduler.Card) func() string {
					return func() string {
						sp := 30
						str := fmt.Sprintf("隨機(先手):%d%%機率引爆炸彈,CD%d\n", 1, thisCard.GetReCoolDown())
						str += fmt.Sprintf("隨機(後手):%d%%機率消除一個非炸彈數字,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
						thisCard.SetDesc(str)
						return str
					}
				},
				OnMythosFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
					return func(g scheduler.Game) (r bool, s string) {
						var strs []string
						sp := 1
						if thisCard.GetCoolDown() > 0 {
							return false, ""
						}
						hit, _, _ := g.GetHit()
						current := g.GetInfoCurrent()
						if hit == current {
							return false, ""
						}

						diceRoll := rand.Intn(100)
						if diceRoll < sp {
							g.SetInfoCurrent(hit)
							strs = append(strs, fmt.Sprintf("%s吉良吉影「到極限了，就是現在，按下去！」", emoji.Emoji(":ghost:")))
							strs = append(strs, fmt.Sprintf("%s皇后殺手「穿心一擊！」%s%d", emoji.Emoji(":ghost:"), emoji.Emoji(":collision:"), hit))
							r = true
							thisCard.ResetCoolDown()
						}
						s = strings.Join(strs, "\n")
						return r, s
					}
				},
				OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
					return func(g scheduler.Game) (r bool, s string) {
						var strs []string
						sp := 30
						if thisCard.GetCoolDown() > 0 {
							return false, ""
						}
						hit, min, max := g.GetHit()
						current := g.GetInfoCurrent()
						if hit == current {
							return false, ""
						}

						diceRoll := rand.Intn(100)
						if diceRoll < sp+thisCard.GetLevel()*2 {

							rg := max - min - 2
							if rg > 1 {
								triggerNumber := rand.Intn(rg)
								triggerNumber += min + 1
								if triggerNumber >= hit {
									triggerNumber++
								}
								switch {
								case triggerNumber < hit:
									min = triggerNumber
								case triggerNumber > hit:
									max = triggerNumber
								}
								g.SetInfoRange(min, max)
								strs = append(strs, fmt.Sprintf("%s吉良吉影「皇后殺手！」%s%d", emoji.Emoji(":ghost:"), emoji.Emoji(":bullseye:"), triggerNumber))
							}
							r = true
						}
						thisCard.ResetCoolDown()
						s = strings.Join(strs, "\n")
						return r, s
					}
				},
			},

		"sakuragi": CardOption{
			CardName:    fmt.Sprintf("%s「你還差得遠呢 越前龍馬」", emoji.Emoji(":ghost:")),
			DisplayName: "你還差得遠呢 越前龍馬",
			Class:       "R",
			CoreSet:     "sakuragi",
			CoolDown:    11,
			ReCoolDown:  11,
			Unique:      true,
			DescFunc: func(thisCard scheduler.Card) func() string {
				return func() string {
					sp := 20
					str := fmt.Sprintf("後手:%d%%機率炸彈隨機轉移玩家,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
					thisCard.SetDesc(str)
					return str
				}
			},
			OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
				return func(g scheduler.Game) (r bool, s string) {
					var strs []string
					if thisCard.GetCoolDown() > 0 || len(g.GetPlayQueue()) <= len(g.GetQueue()) {
						return
					}
					sp := 20
					if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
						queue := g.GetQueue()
						roll := rand.Perm(len(queue))[0]
						toPlayer := queue[roll]
						strs = append(strs, fmt.Sprintf("%s越前龍馬「腳邊截擊!」%s%s", emoji.Emoji(":ghost:"), emoji.Emoji(":right_arrow:"), g.GetPlayer(toPlayer).GetDisplayName()))
						thisCard.ResetCoolDown()
					}
					s = strings.Join(strs, "\n")
					return r, s
				}
			},
		},
	*/
	"blackflash": CardOption{
		CardName:    fmt.Sprintf("%s「黑閃 虎杖悠仁」", emoji.Emoji(":ghost:")),
		DisplayName: "黑閃 虎杖悠仁",
		Class:       "SR",
		CoreSet:     "blackflash",
		CoolDown:    5,
		ReCoolDown:  5,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				str := fmt.Sprintf("後手:%d%%機率使炸彈增加,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
			var i int = 0
			return func(g scheduler.Game) (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					boomCnt := g.GetInfoBoomCnt()
					toCnt := 5 + i*2
					g.MakeInfoBoomCnt(toCnt)
					boomCnt += toCnt
					strs = append(strs, fmt.Sprintf("%s%s%d(%+d)", emoji.Emoji(":ghost:"), emoji.Emoji(":bomb:"), boomCnt, toCnt))
					i++
					thisCard.ResetCoolDown()
				} else {
					i = 0
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"altria": CardOption{
		CardName:    fmt.Sprintf("%s「召喚 Saber 阿爾托莉雅」", emoji.Emoji(":ghost:")),
		DisplayName: "召喚 Saber 阿爾托莉雅",
		Class:       "R",
		CoreSet:     "altria",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("後手:%d%%機率,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
			return func(g scheduler.Game) (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return
				}
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					toPlayerID := g.GetPlayQueueLast()
					strs = append(strs, fmt.Sprintf("%s阿爾托莉雅「%s, 你就是我的Master嗎」", emoji.Emoji(":ghost:"), g.GetPlayer(toPlayerID).GetDisplayName()))
					strs = append(strs, fmt.Sprintf("【%s】獲得%s", g.GetPlayer(toPlayerID).GetDisplayName(), g.GetPlayer(toPlayerID).GetCardPile().TakeCard("saber")))
					thisCard.ResetCoolDown()
					r = true
				}
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
	"starplatinum": CardOption{
		CardName:    fmt.Sprintf("%s「召喚 白金之星」", emoji.Emoji(":ghost:")),
		DisplayName: "召喚 白金之星",
		Class:       "SSR",
		CoreSet:     "starplatinum",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("後手:%d%%機率,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
			return func(g scheduler.Game) (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return
				}
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					toPlayerID := g.GetPlayQueueLast()
					strs = append(strs, fmt.Sprintf("%s「這是替身攻擊!!」", emoji.Emoji(":ghost:")))
					strs = append(strs, fmt.Sprintf("【%s】獲得%s", g.GetPlayer(toPlayerID).GetDisplayName(), g.GetPlayer(toPlayerID).GetCardPile().TakeCard("4")))
					thisCard.ResetCoolDown()
					r = true
				}
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
}
