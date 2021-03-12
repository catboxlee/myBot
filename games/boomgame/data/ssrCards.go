package data

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/scheduler"
	"myBot/helper"
	"myBot/users"
	"strings"
)

// SSRCard ...
var SSRCard = map[string]CardOption{
	"2": CardOption{
		CardName:    "SSR「燕返 不二周助」",
		DisplayName: "燕返 不二周助",
		Class:       "SSR",
		CoreSet:     "2",
		CoolDown:    7,
		ReCoolDown:  7,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 35
				str := fmt.Sprintf("引爆(躲避):%d%%機率炸彈返回上一位玩家,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				if thisCard.GetLevel() >= 4 {
					str += fmt.Sprintf("\nCD時15%%機率發動")
				}
				thisCard.SetDesc(str)
				return str
			}
		},
		OnHitFunc: func(thisCard scheduler.Card) func() (r bool, s string) {
			return func() (r bool, s string) {
				var strs []string
				g := thisCard.GetTopParent()
				playQueue := g.GetPlayQueue()
				sp := 35
				thisPlayer := thisCard.GetParent().GetParent()
				doSkill := func() {
					toUserID := playQueue[helper.Max(0, len(playQueue)-2)]
					strs = append(strs, fmt.Sprintf("【%s】不二周助「燕返!」%s%s", thisPlayer.GetDisplayName(), emoji.Emoji(":right_arrow:"), users.UsersList.Data[toUserID].GetDisplayName()))
					g.AddPlayQueue(toUserID)
				}
				if thisCard.GetCoolDown() > 0 {
					if thisCard.GetLevel() >= 4 && rand.Intn(100) < 15 {
						strs = append(strs, fmt.Sprintf("【%s】不二周助「看起來很有趣的樣子」", thisPlayer.GetDisplayName()))
						doSkill()
						thisCard.ResetCoolDown()
						r = true
						s = strings.Join(strs, "\n")
						return r, s
					}
					return false, ""
				}

				if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
					doSkill()
					r = true
				}
				thisCard.ResetCoolDown()
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
	"3": CardOption{
		CardName:    "SSR「歐拉歐拉 白金之星」",
		DisplayName: "歐拉歐拉 白金之星",
		Class:       "SSR",
		CoreSet:     "3",
		CoolDown:    9,
		ReCoolDown:  9,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 35
				str := fmt.Sprintf("引爆(躲避):%d%%機率炸彈移至下一位玩家,CD%d", sp+thisCard.GetLevel()*2, thisCard.GetReCoolDown())
				if thisCard.GetLevel() >= 4 {
					str += fmt.Sprintf("\nCD時15%%機率發動")
				}
				thisCard.SetDesc(str)
				return str
			}
		},
		OnHitFunc: func(thisCard scheduler.Card) func() (r bool, s string) {
			return func() (r bool, s string) {
				var strs []string
				sp := 35
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				queue := g.GetQueue()
				playQueue := g.GetPlayQueue()

				lastUserID := playQueue[helper.Max(0, len(playQueue)-1)]
				if exist, i := helper.InArray(lastUserID, g.GetQueue()); exist {
					toUserID := queue[(i+1)%len(queue)]
					doSkill := func() {
						strs = append(strs, fmt.Sprintf("【%s】白金之星「歐拉歐拉歐拉歐拉」%s%s", thisPlayer.GetDisplayName(), emoji.Emoji(":right_arrow:"), users.UsersList.Data[toUserID].GetDisplayName()))
						g.OnPlay()
						g.AddPlayQueue(toUserID)
					}
					if thisCard.GetCoolDown() > 0 {
						if thisCard.GetLevel() >= 4 && rand.Intn(100) < 15 {
							strs = append(strs, fmt.Sprintf("【%s】空条承太郎「你失敗的原因只有一個,%s,就是你惹火了我!」", thisPlayer.GetDisplayName(), users.UsersList.Data[toUserID].GetDisplayName()))
							doSkill()
							thisCard.ResetCoolDown()
							r = true
							s = strings.Join(strs, "\n")
							return r, s
						}
						return false, ""
					}
					if rand.Intn(100) < sp+thisCard.GetLevel()*2 {
						doSkill()
						r = true
					}
				}
				thisCard.ResetCoolDown()
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
	"4": CardOption{
		CardName:    "SSR「白金之星 The World! 空条承太郎」",
		DisplayName: "白金之星 The World! 空条承太郎",
		Class:       "SSR",
		CoreSet:     "4",
		CoolDown:    20,
		ReCoolDown:  20,
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
				playQueue := g.GetPlayQueue()
				queue := g.GetQueue()
				lastUserID := playQueue[helper.Max(0, len(playQueue)-1)]
				if exist, i := helper.InArray(lastUserID, g.GetQueue()); exist {
					toUserID := queue[(i+1)%len(queue)]
					strs = append(strs, fmt.Sprintf("【%s】空条承太郎「札．瓦魯斗！」%s%s", thisPlayer.GetDisplayName(), emoji.Emoji(":right_arrow:"), users.UsersList.Data[toUserID].GetDisplayName()))
				}
				thisCard.ResetCoolDown()
				g.Show()
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"11": CardOption{
		CardName:    "SSR「Master，請下指令 瑪修」",
		DisplayName: "寶具 展開 瑪修",
		Class:       "SSR",
		CoreSet:     "11",
		CoolDown:    7,
		ReCoolDown:  7,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				sp := 35
				n := 10
				str := fmt.Sprintf("引爆(防禦):%d%%機率降低傷害%d%%,CD%d", sp, (thisCard.GetLevel()+1)*n, thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnShieldFunc: func(thisCard scheduler.Card) func() (bool, string) {
			say := []string{"寶具 展開", "現為脆弱的雪花之壁", "頌為堅城的雪花之壁", "時而朦朧的白堊之壁", "奮於裁斷的決意之盾", "悲壯的奮起之盾"}
			return func() (r bool, s string) {
				var strs []string
				if thisCard.GetCoolDown() > 0 {
					return false, ""
				}
				sp := 35
				g := thisCard.GetTopParent()
				thisPlayer := thisCard.GetParent().GetParent()
				diceRoll := rand.Intn(100)
				if diceRoll >= 99 {
					strs = append(strs, fmt.Sprintf("【%s】瑪修「前輩最差勁了！」", thisPlayer.GetDisplayName()))
					n := 20
					thisCard.MakeFreeze(n)
					strs = append(strs, fmt.Sprintf("%s%s - %s%s%d(%d)", emoji.Emoji(":Japanese_prohibited_button:"), thisPlayer.GetDisplayName(), thisCard.GetDisplayName(), emoji.Emoji(":Japanese_prohibited_button:"), thisCard.GetFreeze(), n))
				} else if diceRoll < sp+thisCard.GetLevel()*2 {
					n := 10
					if thisCard.GetLevel() >= 4 && rand.Intn(100) < 15 {
						strs = append(strs, fmt.Sprintf("【%s】瑪修「我——生——氣——了！」", thisPlayer.GetDisplayName()))
						g.MakeInfoBoomCnt(-g.GetInfoBoomCnt())
					} else {
						boomCnt := g.GetInfoBoomCnt()
						shiled := int(float64(boomCnt) * (float64((thisCard.GetLevel()+1)*n) / 100))
						g.MakeInfoBoomCnt(-shiled)
						strs = append(strs, fmt.Sprintf("【%s】瑪修「%s」%s%d(%+d)", thisPlayer.GetDisplayName(), say[thisCard.GetLevel()], emoji.Emoji(":collision:"), g.GetInfoBoomCnt(), -shiled))
					}

				}
				thisCard.ResetCoolDown()
				return true, strings.Join(strs, "\n")
			}
		},
	},
}