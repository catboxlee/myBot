package data

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/racegame/scheduler"
	"myBot/helper"
	"strings"
)

// ResidentCard ...
var ResidentCard = map[string]CardOption{
	"green_shell": CardOption{
		CardName:    fmt.Sprintf("「%s綠色龜殼」", emoji.Emoji(":dizzy:")),
		DisplayName: "green_shell",
		CoreSet:     "green_shell",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.SetStop(true)
				strs = append(strs, fmt.Sprintf("%s撞上綠色龜殼,此回合暫停", thisPlayer.GetDisplayName()))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				g := thisPlayer.GetTopParent()
				myRanking := g.GetRanking(thisPlayer.GetUserID())
				if myRanking != 0 {
					arr := g.GetRankingArray()
					targetPlayer := g.GetPlayer(arr[myRanking-1])
					property := targetPlayer.GetProperty()
					property.AddDeBuff("green_shell")
					strs = append(strs, fmt.Sprintf("%s對%s使用綠色龜殼", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))
				} else {
					strs = append(strs, fmt.Sprintf("%s丟棄綠色龜殼", thisPlayer.GetDisplayName()))
				}

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"red_shell": CardOption{
		CardName:    fmt.Sprintf("「%s紅色龜殼」", emoji.Emoji(":dizzy:")),
		DisplayName: "red_shell",
		CoreSet:     "red_shell",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.SetStop(true)
				strs = append(strs, fmt.Sprintf("%s撞上紅色龜殼,此回合暫停", thisPlayer.GetDisplayName()))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				myRanking := g.GetRanking(thisPlayer.GetUserID())
				if myRanking != 0 {
					arr := g.GetRankingArray()
					if args != nil {
						if exist, rk := helper.InArray(args.GetUserID(), arr); exist {
							if rk < myRanking {
								targetPlayer = args
							}
						}
					}
					if targetPlayer == nil {
						targetPlayer = g.GetPlayer(arr[myRanking-1])
					}
					property := targetPlayer.GetProperty()
					property.AddDeBuff("red_shell")
					strs = append(strs, fmt.Sprintf("%s對%s使用紅色龜殼", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))
				} else {
					strs = append(strs, fmt.Sprintf("%s丟棄紅色龜殼", thisPlayer.GetDisplayName()))
				}

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"banana": CardOption{
		CardName:    fmt.Sprintf("「%s香蕉皮」", emoji.Emoji(":dizzy:")),
		DisplayName: "banana",
		CoreSet:     "banana",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.SetStop(true)
				strs = append(strs, fmt.Sprintf("%s撞上香蕉皮,此回合暫停", thisPlayer.GetDisplayName()))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				g := thisPlayer.GetTopParent()
				myRanking := g.GetRanking(thisPlayer.GetUserID())
				if myRanking < len(g.GetQueue())-1 {
					arr := g.GetRankingArray()
					targetPlayer := g.GetPlayer(arr[myRanking+1])
					property := targetPlayer.GetProperty()
					property.AddDeBuff("banana")
					strs = append(strs, fmt.Sprintf("%s對%s使用香蕉皮", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))
				} else {
					strs = append(strs, fmt.Sprintf("%s丟棄香蕉皮", thisPlayer.GetDisplayName()))
				}

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"in-wheel_lift": CardOption{
		CardName:    "「水溝蓋跑法(移除:彎道,髮夾彎)」",
		DisplayName: "水溝蓋跑法",
		CoreSet:     "in-wheel_lift",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    1,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				strs = append(strs, fmt.Sprintf("%s使用%s", thisPlayer.GetDisplayName(), "水溝蓋跑法"))
				if _, str := thisPlayer.RemoveDeBuff("speed_limit4", "speed_limit2"); len(str) > 0 {
					strs = append(strs, str)
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_up1": CardOption{
		CardName:    "「加速1」",
		DisplayName: "加速1",
		CoreSet:     "speed_up1",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				move := 1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「加速1」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := 1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「加速1」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_up2": CardOption{
		CardName:    "「加速2」",
		DisplayName: "加速2",
		CoreSet:     "speed_up2",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := 2
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「加速2」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_down1": CardOption{
		CardName:    "「減速1」",
		DisplayName: "減速1",
		CoreSet:     "speed_down1",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				move := -1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「減速1」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := -1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「減速1」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_down2": CardOption{
		CardName:    "「減速2」",
		DisplayName: "減速2",
		CoreSet:     "speed_down2",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := -2
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「減速2」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"turbo": CardOption{
		CardName:    "「渦輪引擎」",
		DisplayName: "渦輪引擎",
		CoreSet:     "turbo",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := 1
				property := thisPlayer.GetProperty()
				property.MakeDice(1, 0, 0)
				strs = append(strs, fmt.Sprintf("%s「渦輪引擎」%sx%d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), move+1))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"the_star": CardOption{
		CardName:    "「無敵星星」",
		DisplayName: "無敵星星",
		CoreSet:     "the_star",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				move := 1
				property := thisPlayer.GetProperty()
				thisPlayer.RemoveAllDeBuff()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s使用「無敵星星」%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"shield": CardOption{
		CardName:    fmt.Sprintf("「%s保護殼(移除:龜殼,香蕉皮)」", emoji.Emoji(":shield:")),
		DisplayName: "保護殼",
		CoreSet:     "shield",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    4,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				strs = append(strs, fmt.Sprintf("%s使用%s", thisPlayer.GetDisplayName(), "保護殼"))
				if _, str := thisPlayer.RemoveDeBuff("green_shell", "red_shell", "banana"); len(str) > 0 {
					strs = append(strs, str)
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"bandage": CardOption{
		CardName:    fmt.Sprintf("「%sOK繃(移除:骨折)」", emoji.Emoji(":adhesive_bandage:")),
		DisplayName: "OK繃",
		CoreSet:     "bandage",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				strs = append(strs, fmt.Sprintf("%s使用%s", thisPlayer.GetDisplayName(), "OK繃"))
				if _, str := thisPlayer.RemoveDeBuff("broken"); len(str) > 0 {
					strs = append(strs, str)
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"full_power": CardOption{
		CardName:    fmt.Sprintf("「%s體力充沛(移除:體力不支)」", emoji.Emoji(":pill:")),
		DisplayName: "體力充沛",
		CoreSet:     "full_power",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				strs = append(strs, fmt.Sprintf("%s使用體力充沛", thisPlayer.GetDisplayName()))
				if _, str := thisPlayer.RemoveDeBuff("drowsy"); len(str) > 0 {
					strs = append(strs, str)
				}
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"pill": CardOption{
		CardName:    fmt.Sprintf("「%s興奮劑%s+1(移除:骨折,體力不支)」", emoji.Emoji(":pill:"), emoji.Emoji(":game_die:")),
		DisplayName: "興奮劑",
		CoreSet:     "pill",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				strs = append(strs, fmt.Sprintf("%s使用興奮劑%s+1", thisPlayer.GetDisplayName(), emoji.Emoji(":pill:")))
				if _, str := thisPlayer.RemoveDeBuff("broken", "drowsy"); len(str) > 0 {
					strs = append(strs, str)
				}
				thisPlayer.GetProperty().MakeDice(0, 0, 1)
				return true, strings.Join(strs, "\n")
			}
		},
	},
	"broken": CardOption{
		CardName:    fmt.Sprintf("「%s骨折」", emoji.Emoji("::")),
		DisplayName: "骨折",
		CoreSet:     "broken",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnAttackFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				property := thisPlayer.GetProperty()
				property.SetStop(true)
				strs = append(strs, fmt.Sprintf("%s骨折,此回合暫停", thisPlayer.GetDisplayName()))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("broken")
				strs = append(strs, fmt.Sprintf("%s對%s使用骨折", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"drowsy": CardOption{
		CardName:    fmt.Sprintf("「體力不支%s%+d」", emoji.Emoji(":footprints:"), -3),
		DisplayName: "體力不支",
		CoreSet:     "drowsy",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    3,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				move := -3
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「體力不支%s%+d」", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("drowsy")
				strs = append(strs, fmt.Sprintf("%s對%s使用體力不支", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"muddy1": CardOption{
		CardName:    fmt.Sprintf("「泥濘%s%+d」", emoji.Emoji(":footprints:"), -1),
		DisplayName: "泥濘1",
		CoreSet:     "muddy1",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				move := -1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「泥濘%s%+d」", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("muddy1")
				strs = append(strs, fmt.Sprintf("%s對%s使用泥濘", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"downhill1": CardOption{
		CardName:    fmt.Sprintf("「下坡%s%+d」", emoji.Emoji(":footprints:"), 1),
		DisplayName: "下坡1",
		CoreSet:     "downhill1",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnEffectFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				move := 1
				property := thisPlayer.GetProperty()
				property.MakeDice(0, 0, move)
				strs = append(strs, fmt.Sprintf("%s「下坡%s%+d」", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), move))
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = thisPlayer
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("downhill1")
				strs = append(strs, fmt.Sprintf("%s對%s使用下坡", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_need4": CardOption{
		CardName:    fmt.Sprintf("「大峽谷%s>=%d」", emoji.Emoji(":game_die:"), 4),
		DisplayName: "大峽谷4",
		CoreSet:     "speed_need4",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnSpeedLimitFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				if thisPlayer.GetProperty().DiceHit < 4 {
					thisPlayer.GetProperty().SetStop(true)
					strs = append(strs, fmt.Sprintf("%s「大峽谷4」飛越失敗,此回合暫停", thisPlayer.GetDisplayName()))
				} else {
					strs = append(strs, fmt.Sprintf("%s「大峽谷4」飛越成功", thisPlayer.GetDisplayName()))
				}
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("speed_need4")
				strs = append(strs, fmt.Sprintf("%s對%s使用大峽谷4", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_limit4": CardOption{
		CardName:    fmt.Sprintf("「彎道%s<=%d」", emoji.Emoji(":game_die:"), 4),
		DisplayName: "彎道4",
		CoreSet:     "speed_limit4",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnSpeedLimitFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				if thisPlayer.GetProperty().DiceHit > 4 {
					thisPlayer.GetProperty().SetStop(true)
					strs = append(strs, fmt.Sprintf("%s「彎道4」過彎失敗,此回合暫停", thisPlayer.GetDisplayName()))
				} else {
					strs = append(strs, fmt.Sprintf("%s「彎道4」過彎成功", thisPlayer.GetDisplayName()))
				}
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("speed_limit4")
				strs = append(strs, fmt.Sprintf("%s對%s使用彎道4", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
	"speed_limit2": CardOption{
		CardName:    fmt.Sprintf("「髮夾彎%s<=%d」", emoji.Emoji(":game_die:"), 2),
		DisplayName: "髮夾彎2",
		CoreSet:     "speed_limit2",
		CoolDown:    0,
		ReCoolDown:  0,
		Quantity:    2,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("")
				thisCard.SetDesc(str)
				return str
			}
		},
		OnSpeedLimitFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player) (r bool, s string) {
				var strs []string
				if thisPlayer.GetProperty().DiceHit > 2 {
					thisPlayer.GetProperty().SetStop(true)
					strs = append(strs, fmt.Sprintf("%s「髮夾彎」過彎失敗,此回合暫停", thisPlayer.GetDisplayName()))
				} else if thisPlayer.GetProperty().DiceHit == 2 {
					rnd := rand.Intn(6) + 1
					thisPlayer.GetProperty().Move += rnd
					strs = append(strs, fmt.Sprintf("%s「髮夾彎」甩尾過彎%s%+d", thisPlayer.GetDisplayName(), emoji.Emoji(":game_die:"), rnd))
				} else {
					strs = append(strs, fmt.Sprintf("%s「髮夾彎」過彎成功", thisPlayer.GetDisplayName()))
				}
				return true, strings.Join(strs, "\n")
			}
		},
		OnPlayFunc: func(thisCard scheduler.Card) func(thisPlayer scheduler.Player, args scheduler.Player) (bool, string) {
			return func(thisPlayer scheduler.Player, args scheduler.Player) (r bool, s string) {
				var strs []string
				var targetPlayer scheduler.Player
				g := thisPlayer.GetTopParent()
				arr := g.GetRankingArray()
				if args != nil {
					targetPlayer = args
				}
				if targetPlayer == nil {
					targetPlayer = g.GetPlayer(arr[0])
				}
				property := targetPlayer.GetProperty()
				property.AddDeBuff("speed_limit2")
				strs = append(strs, fmt.Sprintf("%s對%s使用髮夾彎2", thisPlayer.GetDisplayName(), targetPlayer.GetDisplayName()))

				return true, strings.Join(strs, "\n")
			}
		},
	},
}
