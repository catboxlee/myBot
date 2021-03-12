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

// ResidentCard ...
var ResidentCard = map[string]CardOption{
	"theworld": CardOption{
		CardName:    "SR「無駄無駄 世界」",
		DisplayName: "無駄無駄 世界",
		CoreSet:     "theworld",
		CoolDown:    20,
		ReCoolDown:  20,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				n := 1
				if thisCard.GetLevel() >= 4 {
					n = 2
				}
				str := fmt.Sprintf("隨機:%d%%機率凍結1位玩家%d張卡片(%d回合),CD%d", sp+thisCard.GetLevel()*2, n, n, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPassFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				sp := 30
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				queue := g.GetQueue()
				if len(queue) < 2 {
					return false, ""
				}
				tmp := rand.Perm(len(queue))
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					n := 1
					if thisCard.GetLevel() >= 4 {
						n = 2
					}
					strs = append(strs, fmt.Sprintf("【%s】世界「無駄無駄無駄無駄」", thisPlayer.GetDisplayName()))
					for i := 0; i < len(tmp); i++ {
						uid := queue[tmp[i]]
						if uid != thisPlayer.GetUserID() {
							if cos := g.GetPlayer(uid).GetRandCards(n); len(cos) > 0 {
								for _, co := range cos {
									strs = append(strs, fmt.Sprintf("%s%s - %s%s(%d)", emoji.Emoji(":Japanese_prohibited_button:"), users.UsersList.Data[uid].GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":Japanese_prohibited_button:"), co.GetFreeze()))
									co.MakeFreeze(n)
								}
							}
						}
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"dio_theworld": CardOption{
		CardName:    "SSR「我不想當人了 Dio」",
		DisplayName: "我不想當人了 Dio",
		CoreSet:     "dio_theworld",
		CoolDown:    80,
		ReCoolDown:  80,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				n := 2
				if thisCard.GetLevel() >= 4 {
					n = 3
				}
				str := fmt.Sprintf("隨機:%d%%機率凍結所有玩家(除自己外)%d張卡片(%d回合),CD%d", sp+thisCard.GetLevel()*2, n, n, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPassFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				sp := 30
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				queue := g.GetQueue()
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					n := 2
					if thisCard.GetLevel() >= 4 {
						n = 3
					}
					strs = append(strs, fmt.Sprintf("【%s】Dio「我不想當人了!」", thisPlayer.GetDisplayName()))
					strs = append(strs, fmt.Sprintf("【%s】Dio「札．瓦魯斗！」", thisPlayer.GetDisplayName()))
					for _, uid := range queue {
						if uid != thisPlayer.GetUserID() {
							if cos := g.GetPlayer(uid).GetRandCards(n); len(cos) > 0 {
								for _, co := range cos {
									strs = append(strs, fmt.Sprintf("%s%s - %s%s(%d)", emoji.Emoji(":Japanese_prohibited_button:"), users.UsersList.Data[uid].GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":Japanese_prohibited_button:"), co.GetFreeze()))
									co.MakeFreeze(n)
								}
							}
						}
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"usopp": CardOption{
		CardName:    "SR「我是不會騙人的 騙人布」",
		DisplayName: "我是不會騙人的 騙人布",
		CoreSet:     "usopp",
		CoolDown:    50,
		ReCoolDown:  50,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 65
				str := fmt.Sprintf("主動:%d%%機率消除一個非炸彈數字,失敗引爆,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func() (bool, string) {
			return func() (r bool, s string) {
				var strs []string
				sp := 65
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				hit, min, max := g.GetHit()
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					tmp := rand.Perm(max - min - 1)
					if exist, index := helper.InArray(hit-min-1, tmp); exist {
						tmp = append(tmp[:index], tmp[index+1:]...)
						x := hit
						if len(tmp) > 0 {
							x = tmp[0] + min + 1
						}
						strs = append(strs, fmt.Sprintf("【%s】騙人布「我是不會騙人的」", thisPlayer.GetDisplayName()))
						strs = append(strs, fmt.Sprintf("【%s】騙人布「鉛星!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":radio_button:"), x))
						g.GamePhase(strconv.Itoa(x))
					}
				} else {
					strs = append(strs, fmt.Sprintf("【%s】騙人布「我說的話連我自己都不信」", thisPlayer.GetDisplayName()))
					strs = append(strs, fmt.Sprintf("【%s】騙人布「臭蛋星!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":radio_button:"), hit))
					g.GamePhase(strconv.Itoa(hit))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
