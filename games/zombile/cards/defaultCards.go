package cards

import (
	"fmt"
	"myBot/emoji"
	"myBot/games/zombile/cards/dice"
	"myBot/games/zombile/power"
	"myBot/helper"
	"strings"
)

// DefaultCardOption ...
type DefaultCardOption struct {
	// general
	id   int
	cost int
	Info
	cardName                    string
	displayname                 string
	cardType                    cardTypeEnum
	CardTraits                  []cardTraitsEnum
	desc                        string
	equipped                    bool // 是否已上場
	usesOption                       // 秏材
	actionTimes                 int
	ActivateFunc                func(*CardOption) func(power.PlayerIF, power.FightIF) string
	OnDisplayFunc               func(*CardOption) func(...interface{}) string
	OnMysterFunc                func(*CardOption) func(...interface{}) string
	OnHealthDamageAfterFunc     func(*CardOption) func(...interface{}) string
	OnHorrorDamageAfterFunc     func(*CardOption) func(...interface{}) string
	OnHealthHurtAfterFunc       func(*CardOption) func(...interface{}) string
	OnHorrorHurtAfterFunc       func(*CardOption) func(...interface{}) string
	OnPlayerHealthHurtAfterFunc func(*CardOption) func(...interface{}) string
}

var baseCards = []struct {
	card DefaultCardOption
	pics int
}{
	{
		DefaultCardOption{
			cost:        3,
			Info:        Info{Damage: 1},
			cardName:    ".45 Automate",
			displayname: ".45 Automate",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.item, CardTraitsValue.weapon, CardTraitsValue.gun},
			desc:        "攻擊: 消秏1彈藥，對目標造成的傷害+1",
			equipped:    false,
			usesOption:  usesOption{true, "彈藥", 4, 1},
			actionTimes: 1,
			ActivateFunc: func(thisCard *CardOption) func(power.PlayerIF, power.FightIF) string {
				return func(targetPlayer power.PlayerIF, target power.FightIF) (r string) {
					var strs []string
					if thisCard.actionTimes == 0 {
						return thisCard.GetDisplayNameWithBracket() + "本回合無法再使用."
					}

					thisPlayer := thisCard.OwnPlayer

					if !thisCard.getEquipped() {
						if target != nil {
							strs = append(strs, fmt.Sprintf("%s將%s扔向了%s.", thisPlayer.GetDisplayNameWithBracket(), thisCard.GetDisplayNameWithBracket(), target.GetDisplayNameWithBracket()))
							strs = append(strs, thisPlayer.Attack(target, power.Damage{Atk: 1, Hor: 0}))
							if targetPlayer == target {
								power.Power.MoveCards(thisPlayer, thisCard, targetPlayer)
								strs = append(strs, fmt.Sprintf("%s獲得%s.", targetPlayer.GetDisplayNameWithBracket(), thisCard.GetDisplayNameWithBracket()))
							}
							return strings.Join(strs, "\n")
						}

						if thisCard.makeEquipped(true) {
							return fmt.Sprintf("[%s]裝備<%s>", thisPlayer.GetDisplayName(), thisCard.GetDisplayName())
						}
					}

					if thisCard.isUses() {
						if ok, s := thisCard.checkUses(); !ok {
							strs = append(strs, s)
							return strings.Join(strs, "\n")
						} else if s := thisCard.spendUses(thisPlayer); len(s) > 0 {
							strs = append(strs, s)
						}
					}

					var targetName string

					if target != nil {
						targetName = target.GetDisplayNameWithBracket()
					} else {
						targetName = "自己"
						target = thisPlayer
					}
					strs = append(strs, fmt.Sprintf("%s舉起%s對%s扣下了板機.", thisPlayer.GetDisplayNameWithBracket(), thisCard.GetDisplayNameWithBracket(), targetName))
					strs = append(strs, thisPlayer.Attack(target, power.Damage{Atk: 2, Hor: 0}))
					thisCard.actionTimes--

					return strings.Join(strs, "\n")
				}
			},
		},
		2,
	},
	{
		DefaultCardOption{
			cost:        3,
			Info:        Info{2, 2, 1, 1, 0},
			cardName:    "Guard dog",
			displayname: "護衛犬",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.ally},
			desc:        "我方玩家受到敵人攻擊時，對該敵人進行攻擊",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			ActivateFunc: func(thisCard *CardOption) func(power.PlayerIF, power.FightIF) string {
				return func(targetPlayer power.PlayerIF, target power.FightIF) (r string) {
					var strs []string

					thisPlayer := thisCard.OwnPlayer

					if !thisCard.getEquipped() {
						if target != nil {
							strs = append(strs, fmt.Sprintf("%s將%s扔向了%s.", thisPlayer.GetDisplayNameWithBracket(), thisCard.GetDisplayNameWithBracket(), target.GetDisplayNameWithBracket()))
							strs = append(strs, thisPlayer.Attack(target, power.Damage{Atk: 1, Hor: 0}))
							strs = append(strs, thisPlayer.Attack(thisCard, power.Damage{Atk: 1, Hor: 0}))
							if targetPlayer == target {
								power.Power.MoveCards(thisPlayer, thisCard, targetPlayer)
								strs = append(strs, fmt.Sprintf("%s獲得%s.", targetPlayer.GetDisplayNameWithBracket(), thisCard.GetDisplayNameWithBracket()))
							}
							return strings.Join(strs, "\n")
						}

						if thisCard.makeEquipped(true) {
							return fmt.Sprintf("<%s>成為你的伙伴.", thisCard.GetDisplayName())
						}
					}

					return strings.Join(strs, "\n")
				}
			},
			OnPlayerHealthHurtAfterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					var strs []string
					if thisCard.getEquipped() {
						if args[0] != thisCard && args[0] != thisCard.OwnPlayer {
							strs = append(strs, fmt.Sprintf("%s攻擊%s.", thisCard.GetDisplayNameWithBracket(), args[0].(power.FightIF).GetDisplayNameWithBracket()))
							strs = append(strs, thisCard.Attack(args[0].(power.FightIF), power.Damage{Atk: 1, Hor: 0}))
						} else {
							strs = append(strs, fmt.Sprintf("%s???.", thisCard.GetDisplayNameWithBracket()))
						}
					}
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
		},
		2,
	},
	{
		DefaultCardOption{
			cost:        4,
			Info:        Info{},
			cardName:    "Bandage",
			displayname: "急救包紮",
			cardType:    cardTypeValue.event,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.item, CardTraitsValue.skill},
			desc:        "使用: 消秏1繃帶，回復生命1",
			equipped:    false,
			usesOption:  usesOption{true, "繃帶", 3, 1},
			ActivateFunc: func(thisCard *CardOption) func(power.PlayerIF, power.FightIF) string {
				return func(targetPlayer power.PlayerIF, target power.FightIF) (r string) {
					var strs []string

					thisPlayer := thisCard.OwnPlayer

					if thisCard.isUses() {
						if ok, s := thisCard.checkUses(); !ok {
							strs = append(strs, s)
							return strings.Join(strs, "\n")
						} else if s := thisCard.spendUses(thisPlayer); len(s) > 0 {
							strs = append(strs, s)
						}
					}

					var targetName string
					if target != nil {
						targetName = target.GetDisplayNameWithBracket()
					} else {
						targetName = "自己"
						target = thisPlayer
					}
					strs = append(strs, fmt.Sprintf("%s對%s使用%s.", thisPlayer.GetDisplayNameWithBracket(), targetName, thisCard.GetDisplayNameWithBracket()))
					strs = append(strs, thisPlayer.Heal(target, power.Damage{Atk: 1, Hor: 0}))

					if thisCard.getQuantity() <= 0 {
						strs = append(strs, fmt.Sprintf("%s%s已秏盡, 移除%s.", thisPlayer.GetDisplayNameWithBracket(), thisCard.getUsesItem(), thisCard.GetDisplayNameWithBracket()))
						thisPlayer.RemoveCards(thisCard)
					}

					return strings.Join(strs, "\n")
				}
			},
		},
		2,
	},
	{
		DefaultCardOption{
			cost:        3,
			Info:        Info{},
			cardName:    "Ammo",
			displayname: "備用彈藥",
			cardType:    cardTypeValue.event,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.item, CardTraitsValue.skill},
			desc:        "行動: 對已裝備的槍械類武器補充彈藥+2",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			ActivateFunc: func(thisCard *CardOption) func(power.PlayerIF, power.FightIF) string {
				return func(targetPlayer power.PlayerIF, target power.FightIF) (r string) {
					r = "目標錯誤, 無法使用."
					thisPlayer := thisCard.OwnPlayer

					switch target.(type) {
					case (*CardOption):
						if target.(*CardOption).getEquipped() && target.(*CardOption).isUses() {
							if exists, _ := helper.InArray(CardTraitsValue.gun, target.(*CardOption).CardTraits); exists {
								if e, s := target.(*CardOption).MakeUses(thisPlayer, 2); e {
									r = s
									thisPlayer.RemoveCards(thisCard)
								}
							}
						}
					}

					return
				}
			},
		},
		2,
	},
	{
		DefaultCardOption{
			cost:        0,
			Info:        Info{4, 4, 3, 1, 1},
			cardName:    "Joker",
			displayname: emoji.Emoji(":ghost:") + "天花板上的小丑",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.Enemy},
			desc:        "隱蔽: 小丑在非戰鬥回合無法被攻擊.",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			actionTimes: 1,
			OnDisplayFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					var strs []string
					if thisCard.makeEquipped(true) {
						//strs = append(strs, fmt.Sprintf("天花板上傳來了奇怪聲響."))
					}
					strs = append(strs, fmt.Sprintf("%s盯上了%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: 0, Hor: 1}))
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
			OnMysterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					if thisCard.actionTimes == 0 {
						return ""
					}
					var strs []string
					if !thisCard.getEquipped() {
						return thisCard.OnDisplayFunc()
					}
					d := dice.Dice
					ap := power.Power.GetPlayersSequence()
					d.Roll(fmt.Sprintf("1d%d", len(ap)))
					if ap[d.Hit-1] == thisCard.OwnPlayer.GetUserID() {
						strs = append(strs, fmt.Sprintf("%s攻擊%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
						strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: thisCard.Info.Damage, Hor: thisCard.Info.Horror}))
					} else {
						moveTo := power.Power.GetPlayer(ap[d.Hit-1])
						power.Power.MoveCards(thisCard.OwnPlayer, thisCard, moveTo)
						strs = append(strs, fmt.Sprintf("%s轉移目標..", thisCard.GetDisplayNameWithBracket()))
						strs = append(strs, fmt.Sprintf("%s盯上了%s.", thisCard.GetDisplayNameWithBracket(), moveTo.GetDisplayNameWithBracket()))
						strs = append(strs, thisCard.Attack(moveTo, power.Damage{Atk: 0, Hor: 1}))
					}
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
		},
		1,
	},
	{
		DefaultCardOption{
			cost:        0,
			Info:        Info{3, 3, 1, 1, 1},
			cardName:    "Vampire",
			displayname: emoji.Emoji(":ghost:") + "吸血鬼",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.Enemy},
			desc:        "吸血: 每次攻擊回復自身生命1",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			actionTimes: 1,
			OnDisplayFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					var strs []string
					if thisCard.makeEquipped(true) {
						//strs = append(strs, fmt.Sprintf("%s自動裝備.", thisCard.GetDisplayNameWithBracket()))
					}
					strs = append(strs, fmt.Sprintf("%s盯上了%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: 0, Hor: 1}))
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
			OnMysterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					if thisCard.actionTimes == 0 {
						return ""
					}
					var strs []string
					if !thisCard.getEquipped() {
						return thisCard.OnDisplayFunc()
					}
					strs = append(strs, fmt.Sprintf("%s攻擊%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: thisCard.Info.Damage, Hor: thisCard.Info.Horror}))
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
			OnHealthDamageAfterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					thisCard.MakeHealth(1)
					return fmt.Sprintf("<%s>回復1", thisCard.GetDisplayName())
				}
			},
		},
		2,
	},
	{
		DefaultCardOption{
			cost:        0,
			Info:        Info{5, 5, 1, 2, 1},
			cardName:    "Vampire",
			displayname: emoji.Emoji(":ghost:") + "吸血鬼。艾德華",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.Enemy},
			desc:        "吸血: 每次攻擊回復自身生命1",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			actionTimes: 1,
			OnDisplayFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					var strs []string
					if thisCard.makeEquipped(true) {
						//strs = append(strs, fmt.Sprintf("%s自動裝備.", thisCard.GetDisplayNameWithBracket()))
					}
					strs = append(strs, fmt.Sprintf("%s盯上了%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: 0, Hor: 1}))
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
			OnMysterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					if thisCard.actionTimes == 0 {
						return ""
					}
					var strs []string
					if !thisCard.getEquipped() {
						return thisCard.OnDisplayFunc()
					}
					strs = append(strs, fmt.Sprintf("%s攻擊%s.", thisCard.GetDisplayNameWithBracket(), thisCard.OwnPlayer.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(thisCard.OwnPlayer, power.Damage{Atk: thisCard.Info.Damage, Hor: thisCard.Info.Horror}))
					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
			OnHealthDamageAfterFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					thisCard.MakeHealth(1)
					return fmt.Sprintf("<%s>回復1", thisCard.GetDisplayName())
				}
			},
		},
		1,
	},
	{
		DefaultCardOption{
			cost:        0,
			Info:        Info{Health: 4, HealthMax: 4, Combat: 1, Damage: 1},
			cardName:    "Vampire Hunter",
			displayname: "吸血鬼獵人",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.ally},
			desc:        "",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			actionTimes: 1,
			ActivateFunc: func(thisCard *CardOption) func(power.PlayerIF, power.FightIF) string {
				return func(targetPlayer power.PlayerIF, target power.FightIF) (r string) {
					if thisCard.actionTimes == 0 {
						return ""
					}

					var strs []string
					thisPlayer := thisCard.OwnPlayer

					if !thisCard.getEquipped() {
						if thisCard.makeEquipped(true) {
							thisCard.actionTimes--
							return fmt.Sprintf("<%s>成為你的伙伴.", thisCard.GetDisplayName())
						}
					}
					if target == nil {
						target = thisPlayer
					}

					strs = append(strs, fmt.Sprintf("%s攻擊%s.", thisCard.GetDisplayNameWithBracket(), target.GetDisplayNameWithBracket()))
					strs = append(strs, thisCard.Attack(target, power.Damage{Atk: thisCard.Info.Damage, Hor: thisCard.Info.Horror}))

					thisCard.actionTimes--
					return strings.Join(strs, "\n")
				}
			},
		},
		1,
	},
	{
		DefaultCardOption{
			cost:        0,
			Info:        Info{},
			cardName:    "Psychasthenia", // Mental weakness
			displayname: "精神衰弱",
			cardType:    cardTypeValue.asset,
			CardTraits:  []cardTraitsEnum{CardTraitsValue.weakness},
			desc:        "受到恐懼傷害時, 恐懼額外+1.",
			equipped:    false,
			usesOption:  usesOption{uses: false},
			OnDisplayFunc: func(thisCard *CardOption) func(args ...interface{}) string {
				return func(args ...interface{}) (r string) {
					var strs []string
					if thisCard.makeEquipped(true) {
						strs = append(strs, fmt.Sprintf("%s自動裝備.", thisCard.GetDisplayNameWithBracket()))
					}
					return strings.Join(strs, "\n")
				}
			},
			OnHorrorDamageAfterFunc: func(thisCard *CardOption) func(ags ...interface{}) string {
				return func(ags ...interface{}) (r string) {
					thisCard.OwnPlayer.MakeHorror(1)
					r = fmt.Sprintf("[%s]<%s>恐懼+1(%d)", thisCard.OwnPlayer.GetDisplayName(), thisCard.GetDisplayName(), thisCard.OwnPlayer.GetHorror())
					return
				}
			},
		},
		1,
	},
}
