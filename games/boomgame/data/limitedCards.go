package data

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/scheduler"
	"myBot/helper"
	"myBot/users"
	"strconv"
	"strings"
)

// LimitedCard ...
var LimitedCard = map[string]CardOption{
	"1": CardOption{
		CardName:    "SSR「狙擊之王 騙人布」",
		DisplayName: "狙擊之王 騙人布",
		Class:       "SSR",
		CoreSet:     "1",
		CoolDown:    50,
		ReCoolDown:  50,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("主動:99%%機率消除一個非炸彈數字,失敗引爆,CD%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, "技能CD中.."
				}
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				hit, min, max := g.GetHit()
				if rand.Intn(100) < 99 {
					tmp := rand.Perm(max - min - 1)
					if exist, index := helper.InArray(hit-min-1, tmp); exist {
						tmp = append(tmp[:index], tmp[index+1:]...)
						x := hit
						if len(tmp) > 0 {
							x = tmp[0] + min + 1
						}
						strs = append(strs, fmt.Sprintf("【%s】狙擊王「我是狙擊手，支援是我的天職」", thisPlayer.GetDisplayName()))
						strs = append(strs, fmt.Sprintf("【%s】狙擊王「必殺。火鳥星!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":bullseye:"), x))

						g.GamePhase(strconv.Itoa(x))
					}
				} else {
					strs = append(strs, fmt.Sprintf("【%s】騙人布「我得了不說謊就會死的病」", thisPlayer.GetDisplayName()))
					strs = append(strs, fmt.Sprintf("【%s】騙人布「騙人布橡皮筯!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), hit))
					g.GamePhase(strconv.Itoa(hit))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"5": CardOption{
		CardName:    "SR「燃燒吧 小宇宙 星矢」",
		DisplayName: "燃燒吧 小宇宙 星矢",
		Class:       "SR",
		CoreSet:     "5",
		CoolDown:    9,
		ReCoolDown:  9,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				str := fmt.Sprintf("引爆(防禦):%d%%機率鎖血(99),CD%d", sp, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnShieldFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				sp := 30
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				diceRoll := rand.Intn(100)
				if diceRoll < sp+thisCard.GetLevel()*2 {
					boomCnt := g.GetInfoBoomCnt()
					rankBoomCnt := g.GetRankBoomCnt(thisPlayer.GetUserID())
					if rankBoomCnt+boomCnt > 100 {
						shiled := rankBoomCnt + boomCnt - 100
						g.MakeInfoBoomCnt(-shiled)
						strs = append(strs, fmt.Sprintf("【%s】星矢「燃燒吧 小宇宙」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), g.GetInfoBoomCnt(), -shiled))
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"12": CardOption{
		CardName:    "SR「龍抬頭 一葉之秋」",
		DisplayName: "龍抬頭 一葉之秋",
		Class:       "SR",
		CoreSet:     "12",
		CoolDown:    13,
		ReCoolDown:  13,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("引爆(先手):%d%%機率更改炸彈數字,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnSenteFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				hit, min, max := g.GetHit()
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					if hit == g.GetInfoCurrent() {
						tmp := rand.Perm(max - min - 1)
						if exist, index := helper.InArray(hit-min-1, tmp); exist {
							tmp = append(tmp[:index], tmp[index+1:]...)
							x := hit
							if len(tmp) > 0 {
								x = tmp[0] + min + 1
							}
							strs = append(strs, fmt.Sprintf("【%s】一葉之秋「龍抬頭」%s", thisPlayer.GetDisplayName(), emoji.Emoji(":right_arrow_curving_up:")))
							g.SetHit(x, min, max)
						}
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"13": CardOption{
		CardName:    "SSR「SSR突破卡」",
		DisplayName: "SSR突破卡",
		Class:       "SSR",
		CoreSet:     "13",
		CoolDown:    0,
		ReCoolDown:  0,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
	},
	"14": CardOption{
		CardName:    "SR「SR突破卡」",
		DisplayName: "SR突破卡",
		Class:       "SR",
		CoreSet:     "14",
		CoolDown:    0,
		ReCoolDown:  0,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
	},
	"15": CardOption{
		CardName:    "R「R突破卡」",
		DisplayName: "R突破卡",
		CoreSet:     "15",
		CoolDown:    0,
		ReCoolDown:  0,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
	},
	"19": CardOption{
		CardName:    "SSR「一拳超人  埼玉」",
		DisplayName: "一拳超人  埼玉",
		Class:       "SSR",
		CoreSet:     "19",
		CoolDown:    55,
		ReCoolDown:  55,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("主動:99%%機率消除大量數字,失敗引爆,CD%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, "技能CD中.."
				}
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				hit, min, max := g.GetHit()
				diceRoll := rand.Intn(100)
				if diceRoll < 99 {
					if diceRoll == 0 {
						newMin := hit - 1
						newMax := hit + 1
						g.SetInfoRange(newMin, newMax)
						g.MakeInfoBoomCnt((max - min) - (newMax - newMin))
						strs = append(strs, fmt.Sprintf("【%s】埼玉「認真的一拳!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), g.GetInfoBoomCnt(), (max-min)-(newMax-newMin)))
						g.Show()
					} else {
						newMin := hit - rand.Intn(hit-min) + 1
						newMax := hit + rand.Intn(max-hit) + 1
						g.SetInfoRange(newMin, newMax)
						g.MakeInfoBoomCnt((max - min) - (newMax - newMin))
						strs = append(strs, fmt.Sprintf("【%s】埼玉「普通的一拳」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), g.GetInfoBoomCnt(), (max-min)-(newMax-newMin)))
						g.Show()
					}
				} else {
					g.MakeInfoBoomCnt((max - min - 1))
					strs = append(strs, fmt.Sprintf("【%s】禿頭披風俠「認真掀桌!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), g.GetInfoBoomCnt(), (max-min-1)))
					g.GamePhase(strconv.Itoa(hit))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"saber": CardOption{
		CardName:    "SSR「召喚 誓約勝利之劍 阿爾托莉雅」",
		DisplayName: "召喚 誓約勝利之劍 阿爾托莉雅",
		Class:       "SSR",
		CoreSet:     "saber",
		CoolDown:    0,
		ReCoolDown:  0,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 100
				str := fmt.Sprintf("引爆(攻擊):%d%%機率使炸彈*%d,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown(), (thisCard.GetLevel() + 2))
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 100
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					toCnt := g.GetInfoBoomCnt() * (thisCard.GetLevel() + 2)
					g.MakeInfoBoomCnt(toCnt)
					boomCnt := g.GetInfoBoomCnt()
					strs = append(strs, fmt.Sprintf("【%s】 阿爾托莉雅「誓約勝利之劍!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), boomCnt, toCnt))
				}
				thisCard.ResetCoolDown()
				thisPlayer.GetCardPile().UsedCard("saber")
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"4": CardOption{
		CardName:    "SSR「白金之星 The World! 空条承太郎」",
		DisplayName: "白金之星 The World! 空条承太郎",
		Class:       "SSR",
		CoreSet:     "4",
		CoolDown:    0,
		ReCoolDown:  0,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("主動:炸彈移至下一位玩家,CD%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, "技能CD中..."
				}
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				g.OnPlay()
				toUserID := g.GetQueueNext()
				strs = append(strs, fmt.Sprintf("【%s】白金之星「札．瓦魯斗！」%s%s", thisPlayer.GetDisplayName(), emoji.Emoji(":right_arrow:"), users.UsersList.Data[toUserID].GetDisplayName()))
				thisCard.ResetCoolDown()
				g.Show()
				thisPlayer.GetCardPile().UsedCard("4")
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
