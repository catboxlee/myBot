package data

import (
	"fmt"
	"myBot/emoji"
	"myBot/games/racegame/scheduler"
	"strings"
)

// LimitedCard ...
var LimitedCard = map[string]CardOption{
	"1": CardOption{
		CardName:    "「一騎絕塵」",
		DisplayName: "一騎絕塵",
		Class:       "SSR",
		CoreSet:     "pacer",
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				var strs []string
				str := strings.Join(strs, "\n")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 1, 0)
				if property.GetTurn() == 1 {
					property.MakeDice(0, 1, 3)
					strs = append(strs, fmt.Sprintf("%s「一騎絕塵」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), 3))
				} else if thisPlayer.GetTopParent().GetRanking(thisPlayer.GetUserID()) == 0 {
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"2": CardOption{
		CardName:    "「赤兔」",
		DisplayName: "赤兔",
		Class:       "SSR",
		CoreSet:     "pacer",
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				var strs []string
				str := strings.Join(strs, "\n")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, 1)
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"3": CardOption{
		CardName:    "「攻城車」",
		DisplayName: "攻城車",
		Class:       "SSR",
		CoreSet:     "pacer",
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				var strs []string
				str := strings.Join(strs, "\n")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 2, 0)
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"4": CardOption{
		CardName:    "「刺客」",
		DisplayName: "刺客",
		Class:       "SSR",
		CoreSet:     "assassin",
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				g := thisPlayer.GetTopParent()
				property := thisPlayer.GetProperty()
				if property.TotalMove >= (g.GetMeter() - g.GetMeter()/2) {
					move := 2
					property.MakeDice(0, 2, move)
					strs = append(strs, fmt.Sprintf("%s「刺客」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), move))
					if thisPlayer.GetTopParent().GetRanking(thisPlayer.GetUserID()) != 0 && property.TotalMove >= (g.GetMeter()-g.GetMeter()/3) {
						move++
						property.MakeDice(0, 1, 1)
						strs = append(strs, fmt.Sprintf("%s「刺客:背刺」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), 1))
					}
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"5": CardOption{
		CardName:    "「毒牙」",
		DisplayName: "毒牙",
		Class:       "SSR",
		CoreSet:     "snake",
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				g := thisPlayer.GetTopParent()
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 1, 0)
				if thisPlayer.GetTurn() > 1 {
					for _, userID := range g.GetQueue() {
						if userID != thisPlayer.GetUserID() {
							if property.TotalMove >= g.GetPlayer(userID).GetProperty().TotalMove && g.GetPlayer(userID).GetProperty().TotalMove >= property.TotalMove-5 {
								property.MakeDice(0, 0, 1)
								strs = append(strs, fmt.Sprintf("%s「毒牙」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), 1))
								g.GetPlayer(userID).AddDeBuff("speed_down1")
								strs = append(strs, fmt.Sprintf("%s「減速1」", g.GetPlayer(userID).GetDisplayName()))
							}
						}
					}
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
}
