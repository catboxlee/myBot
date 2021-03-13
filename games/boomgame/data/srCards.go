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

// SRCard ...
var SRCard = map[string]CardOption{
	"6": CardOption{
		CardName:    "SR「無駄無駄 世界」",
		DisplayName: "無駄無駄 世界",
		Class:       "SR",
		CoreSet:     "6",
		CoolDown:    9,
		ReCoolDown:  9,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				n := 1
				if thisCard.GetLevel() >= 4 {
					n = 2
				}
				str := fmt.Sprintf("先手:%d%%機率凍結1位玩家%d張卡片(%d回合),CD%d", sp+thisCard.GetLevel()*2, n, n, thisCard.GetReCoolDown())
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
									co.MakeFreeze(n)
									strs = append(strs, fmt.Sprintf("%s%s - %s%s%d(%+d)", emoji.Emoji(":Japanese_prohibited_button:"), users.UsersList.Data[uid].GetDisplayName(), co.GetDisplayName(), emoji.Emoji(":Japanese_prohibited_button:"), co.GetFreeze(), n))
								}
								g.GetPlayer(uid).SaveData()
							}
						}
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"7": CardOption{
		CardName:    "SR「我是不會騙人的 騙人布」",
		DisplayName: "我是不會騙人的 騙人布",
		Class:       "SR",
		CoreSet:     "7",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
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
				if thisCard.GetCoolDown() > 0 {
					return false, "技能CD中..."
				}
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
						strs = append(strs, fmt.Sprintf("【%s】騙人布「鉛星!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":bullseye:"), x))
						g.GamePhase(strconv.Itoa(x))
					}
				} else {
					strs = append(strs, fmt.Sprintf("【%s】騙人布「我說的話連我自己都不信」", thisPlayer.GetDisplayName()))
					strs = append(strs, fmt.Sprintf("【%s】騙人布「臭蛋星!」%s%d", thisPlayer.GetDisplayName(), emoji.Emoji(":bullseye:"), hit))
					g.GamePhase(strconv.Itoa(hit))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"20": CardOption{
		CardName:    "SR「我是不會輸的 埼玉」",
		DisplayName: "我是不會輸的 埼玉",
		Class:       "SR",
		CoreSet:     "20",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				str := fmt.Sprintf("先手:%d%%機率消除炸彈,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
					shiled := rand.Intn(boomCnt)
					g.MakeInfoBoomCnt(-shiled)
					strs = append(strs, fmt.Sprintf("【%s】埼玉「我是不會輸的，地球，由我來守護」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), g.GetInfoBoomCnt(), -shiled))
					if diceRoll <= thisCard.GetLevel() {
						shiled = boomCnt
						g.MakeInfoBoomCnt(-shiled)
						strs = append(strs, fmt.Sprintf("【%s】埼玉「沒有一拳解決不了的事，如果有，就兩拳」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), g.GetInfoBoomCnt(), -shiled))
					}
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"21": CardOption{
		CardName:    fmt.Sprintf("SR「（ ﾟ∇ﾟ)つ%s 娜美」", emoji.Emoji(":gem_stone:")),
		DisplayName: fmt.Sprintf("（ ﾟ∇ﾟ)つ%s 娜美", emoji.Emoji(":gem_stone:")),
		Class:       "SR",
		CoreSet:     "21",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				str := fmt.Sprintf("後手:%d%%機率獲得鑽石,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				thisPlayer := thisCard.GetParent().GetParent()
				sp := 30
				diceRoll := rand.Intn(100)
				if diceRoll < sp+thisCard.GetLevel()*2 {
					toCnt := rand.Intn(100)
					if diceRoll == 0 {
						toCnt = 250
					}
					thisPlayer.MakeGemStone(toCnt)
					gem := thisPlayer.GetGemStone()
					strs = append(strs, fmt.Sprintf("【%s】娜美「我喜歡錢和橘子」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":gem_stone:"), gem, toCnt))
				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"22": CardOption{
		CardName:    "SR「瘋狂鑽石 東方杖助」",
		DisplayName: "瘋狂鑽石 東方杖助",
		Class:       "SR",
		CoreSet:     "22",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 30
				str := fmt.Sprintf("後手:%d%%機率消除炸彈恢復傷害,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
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
				diceRoll := rand.Intn(100)
				boomCnt := g.GetInfoBoomCnt()
				rankBoomCnt := g.GetRankBoomCnt(thisPlayer.GetUserID())
				if diceRoll < sp+thisCard.GetLevel()*2 {
					toCnt := rand.Intn(int(boomCnt/2) + 1)
					if diceRoll == 0 {
						toCnt = boomCnt
					}
					toCnt = helper.Min(toCnt, rankBoomCnt)
					g.MakeInfoBoomCnt(-toCnt)
					g.MakeRankBoomCnt(thisPlayer.GetUserID(), -toCnt)
					strs = append(strs, fmt.Sprintf("【%s】瘋狂鑽石「嘟啦啦啦啦啦啦啦」%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":bomb:"), g.GetInfoBoomCnt(), -toCnt))
					strs = append(strs, fmt.Sprintf("【%s】%s%d(%s%d)", thisPlayer.GetDisplayName(), emoji.Emoji(":collision:"), g.GetRankBoomCnt(thisPlayer.GetUserID()), emoji.Emoji(":sparkling_heart:"), toCnt))
				}
				//thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
