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
	"gandalf": CardOption{
		CardName:    fmt.Sprintf("%s「灰袍巫師 甘道夫」", emoji.Emoji(":ghost:")),
		DisplayName: "灰袍巫師 甘道夫",
		Class:       "R",
		CoreSet:     "gandalf",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("隨機(後手):%d%%機率須再喊1次,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
			return func(g scheduler.Game) (r bool, s string) {
				var strs []string
				sp := 20
				queue := g.GetPlayQueue()
				if thisCard.GetCoolDown() > 0 {
					if thisCard.GetReCoolDown()-thisCard.GetCoolDown() < 2 {
						diceRoll := rand.Intn(100)
						if diceRoll < sp+thisCard.GetLevel()*2 {
							strs = append(strs, fmt.Sprintf("%s甘道夫「You shall not PASS!!」%s%s", emoji.Emoji(":ghost:"), emoji.Emoji(":right_arrow:"), g.GetPlayer(queue[len(queue)-1]).GetDisplayName()))
							strs = append(strs, fmt.Sprintf("%d", thisCard.GetCoolDown()))
							r = true
							s = strings.Join(strs, "\n")
							return r, s
						}
					}
					return false, ""
				}

				diceRoll := rand.Intn(100)
				if diceRoll < sp+thisCard.GetLevel()*2 {
					strs = append(strs, fmt.Sprintf("%s甘道夫「You can not PASS!」%s%s", emoji.Emoji(":ghost:"), emoji.Emoji(":right_arrow:"), g.GetPlayer(queue[len(queue)-1]).GetDisplayName()))
					r = true
				}
				thisCard.ResetCoolDown()
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
	"blackflash": CardOption{
		CardName:    fmt.Sprintf("%s「黑閃 虎杖悠仁」", emoji.Emoji(":ghost:")),
		DisplayName: "黑閃 虎杖悠仁",
		Class:       "R",
		CoreSet:     "blackflash",
		CoolDown:    5,
		ReCoolDown:  5,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
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
					strs = append(strs, fmt.Sprintf("%s虎杖悠仁「黑閃!」%s%d(%+d)", emoji.Emoji(":ghost:"), emoji.Emoji(":bomb:"), boomCnt, toCnt))
					i++
				} else {
					i = 0
					thisCard.ResetCoolDown()
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
