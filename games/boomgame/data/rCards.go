package data

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/scheduler"
	"myBot/users"
	"strings"
)

// RCard ...
var RCard = map[string]CardOption{
	/*
		"8": CardOption{
			CardName:    fmt.Sprintf("R「%s」", emoji.Emoji(":gem_stone:")),
			DisplayName: fmt.Sprintf("%s", emoji.Emoji(":gem_stone:")),
			Class:       "R",
			CoreSet:     "8",
			CoolDown:    0,
			ReCoolDown:  0,
			Unique:      false,
			DescFunc: func(thisCard scheduler.Card) func() string {
				return func() string {
					str := fmt.Sprintf("使用:獲得鑽石%s(0~20)", emoji.Emoji(":gem_stone:"))
					thisCard.SetDesc(str)
					return str
				}
			},
			OnPlayFunc: func(thisCard scheduler.Card) func() (bool, string) {
				return func() (r bool, s string) {
					var strs []string
					thisPlayer := thisCard.GetParent().GetParent()
					n := rand.Intn(20)
					thisPlayer.MakeGemStone(n)
					strs = append(strs, fmt.Sprintf("【%s】獲得%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":gem_stone:"), thisPlayer.GetGemStone(), n))
					thisPlayer.GetCardPile().UsedCard("8")
					users.UsersList.SaveData(thisPlayer.GetUserID())
					thisPlayer.SaveData()
					return true, strings.Join(strs, "\n")
				}
			},
		},
	*/
	"9": CardOption{
		CardName:    "R「時間刪除 克里姆王」",
		DisplayName: "時間刪除 克里姆王",
		Class:       "R",
		CoreSet:     "9",
		CoolDown:    11,
		ReCoolDown:  11,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				n := 1
				if thisCard.GetLevel() >= 4 {
					n = 2
				}
				str := fmt.Sprintf("先手:%d%%機率卡片CD/凍結減少%d,CD%d", sp+thisCard.GetLevel()*2, n, thisCard.GetReCoolDown())
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
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 20
				n := rand.Intn(3) + 1
				if thisCard.GetLevel() >= 4 {
					n++
				}
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					strs = append(strs, fmt.Sprintf("【%s】克里姆王「時間刪除」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":hourglass_not_done:"), -n))
					for _, co := range thisPlayer.GetCardPile().GetCards() {
						if co.GetCoolDown() > 0 {
							co.MakeCoolDown(-n)
							//strs = append(strs, fmt.Sprintf("【%s】%s%s%d(%+d)", thisPlayer.GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":hourglass_not_done:"), co.GetCoolDown(), -n))
						}
						if co.GetFreeze() > 0 {
							co.MakeFreeze(-n)
							//strs = append(strs, fmt.Sprintf("【%s】%s%s%d(%+d)", thisPlayer.GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":Japanese_prohibited_button:"), co.GetFreeze(), -n))
						}
					}
				}
				users.UsersList.SaveData(thisPlayer.GetUserID())
				thisPlayer.SaveData()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"10": CardOption{
		CardName:    "R「遲鈍果實 銀狐福克西」",
		DisplayName: "遲鈍果實 銀狐福克西",
		Class:       "R",
		CoreSet:     "10",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				n := 1
				if thisCard.GetLevel() >= 4 {
					n = 2
				}
				str := fmt.Sprintf("先手:%d%%機率使1位玩家%d張卡片CD增加%d,CD%d", sp+thisCard.GetLevel()*2, n, 30, thisCard.GetReCoolDown())
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
				queue := g.GetQueue()
				tmp := rand.Perm(len(queue))
				if len(queue) < 2 {
					return false, ""
				}
				sp := 20
				n := 5
				if thisCard.GetLevel() >= 4 {
					n = 6
				}
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					strs = append(strs, fmt.Sprintf("【%s】銀狐福克西「遲鈍光線」", thisPlayer.GetDisplayName()))
					for i := 0; i < len(tmp); i++ {
						uid := queue[tmp[i]]
						if uid != thisPlayer.GetUserID() {
							strs = append(strs, fmt.Sprintf("%s%s(%+d)", users.UsersList.Data[uid].GetDisplayName(), emoji.Emoji(":hourglass_not_done:"), n))
							if cos := g.GetPlayer(uid).GetRandCards(n); len(cos) > 0 {
								for _, co := range cos {
									co.MakeCoolDown(5)
									//strs = append(strs, fmt.Sprintf("%s%s - %s%s%d(%+d)", emoji.Emoji(":hourglass_not_done:"), users.UsersList.Data[uid].GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":hourglass_not_done:"), co.GetCoolDown(), 5))
								}
								g.GetPlayer(uid).SaveData()
							}
							break
						}
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"16": CardOption{
		CardName:    "R「鋼之鍊金術師 愛德華·愛力克」",
		DisplayName: "鋼之鍊金術師 愛德華·愛力克",
		Class:       "R",
		CoreSet:     "16",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("先手:%d%%機率使炸彈轉成鑽石,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					boomCnt := g.GetInfoBoomCnt()
					toCnt := rand.Intn(boomCnt)
					g.MakeInfoBoomCnt(-toCnt)
					boomCnt -= toCnt
					strs = append(strs, fmt.Sprintf("【%s】愛德華「鍊成!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), boomCnt, -toCnt))
					thisPlayer.MakeGemStone(toCnt)
					gem := thisPlayer.GetGemStone()
					strs = append(strs, fmt.Sprintf("【%s】%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":gem_stone:"), gem, toCnt))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"17": CardOption{
		CardName:    "R「爆裂魔法 惠惠」",
		DisplayName: "爆裂魔法 惠惠",
		Class:       "R",
		CoreSet:     "17",
		CoolDown:    7,
		ReCoolDown:  7,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("後手:%d%%機率使炸彈增加,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 30
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					boomCnt := g.GetInfoBoomCnt()
					toCnt := rand.Intn(2) + 3
					g.MakeInfoBoomCnt(toCnt)
					boomCnt += toCnt
					strs = append(strs, fmt.Sprintf("【%s】惠惠「Explosion!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), boomCnt, toCnt))
					thisCard.ResetCoolDown()
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"18": CardOption{
		CardName:    "R「幫我撐10秒 桐人」",
		DisplayName: "幫我撐10秒 桐人",
		Class:       "R",
		CoreSet:     "18",
		CoolDown:    10,
		ReCoolDown:  10,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("引爆(攻擊):%d%%機率使炸彈增加,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					toCnt := rand.Intn(16) + 1
					g.MakeInfoBoomCnt(toCnt)
					boomCnt := g.GetInfoBoomCnt()
					strs = append(strs, fmt.Sprintf("【%s】桐人「星爆氣流斬!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), boomCnt, toCnt))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"21": CardOption{
		CardName:    "R「超電磁砲 御坂美琴」",
		DisplayName: "超電磁砲 御坂美琴",
		Class:       "R",
		CoreSet:     "18",
		CoolDown:    9,
		ReCoolDown:  9,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 20
				str := fmt.Sprintf("引爆(攻擊):%d%%機率使炸彈增加,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 20
				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					toCnt := 5
					g.MakeInfoBoomCnt(toCnt)
					boomCnt := g.GetInfoBoomCnt()
					strs = append(strs, fmt.Sprintf("【%s】御坂美琴「超電磁砲!」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), boomCnt, toCnt))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
