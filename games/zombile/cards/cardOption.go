package cards

import (
	"fmt"
	"my/games/zombile/lang"
	"my/games/zombile/power"
	"my/helper"
	"strings"
)

// Info ...
type Info struct {
	Health    int
	HealthMax int
	Combat    int
	Damage    int
	Horror    int
}

// CardOption ...
type CardOption struct {
	// general
	id       int
	sourceID int
	cost     int
	Info
	cardName                    string
	displayname                 string
	cardType                    cardTypeEnum
	CardTraits                  []cardTraitsEnum
	desc                        string
	equipped                    bool // 是否已上場
	usesOption                       // 秏材
	OwnPlayer                   power.PlayerIF
	actionTimes                 int
	ActivateFunc                func(power.PlayerIF, power.FightIF) string
	OnDisplayFunc               func(...interface{}) string
	OnMysterFunc                func(...interface{}) string
	OnHealthDamageAfterFunc     func(...interface{}) string
	OnHorrorDamageAfterFunc     func(...interface{}) string
	OnPlayerHealthHurtAfterFunc func(...interface{}) string
	//OnMysterFunc func(power.PlayerIF, int) string
	//Method      cardMethodInterface
	//supplies    int
	//consume     consumeType
}

type cardMethodInterface interface {
	clone() cardMethodInterface
	onPlay(power.PlayerIF, *CardOption) (bool, string)
	viewCardInfo() string
}

func (ao *assetOption) clone() cardMethodInterface {
	n := *ao
	return &n
}

func (eo *eventOption) clone() cardMethodInterface {
	n := *eo
	return &n
}

type cardTypeEnum = string

var cardTypeValue = struct {
	asset cardTypeEnum
	event cardTypeEnum
}{
	"Asset",
	"Event",
}

type cardTraitsEnum = string

// CardTraitsValue ...
var CardTraitsValue = struct {
	ally     cardTraitsEnum
	minion   cardTraitsEnum
	item     cardTraitsEnum
	weapon   cardTraitsEnum
	skill    cardTraitsEnum
	Enemy    cardTraitsEnum
	gun      cardTraitsEnum
	weakness cardTraitsEnum
}{
	"Ally",
	"Minion",
	"Item",
	"Weapon",
	"Skill",
	"Enemy",
	"Gun",
	"Weakness",
}

// OnPlay ...
func (co *CardOption) OnPlay(currentPlayer power.PlayerIF, targetPlayer power.PlayerIF, target power.FightIF) string {
	var strs []string
	if co.ActivateFunc != nil {
		strs = append(strs, co.ActivateFunc(targetPlayer, target))
	}
	return strings.Join(strs, "\n")
}

// ResetActionTimes ...
func (co *CardOption) ResetActionTimes() {
	co.actionTimes = baseCards[co.sourceID].card.actionTimes
}

// getEquipped 是否已上場
func (co *CardOption) getEquipped() bool {
	return co.equipped
}

// setEquipped 設定是否已上場
func (co *CardOption) setEquipped(b bool) {
	co.equipped = b
}

// makeEquipped 上場
func (co *CardOption) makeEquipped(b bool) bool {
	if !co.getEquipped() {
		co.setEquipped(b)
		return true
	}
	return false
}

// checkUses ...
func (co *CardOption) checkUses() (bool, string) {
	if co.isUses() {
		if 0 > co.quantity-co.spend {
			return false, fmt.Sprintf("<%s>%s不足", co.GetDisplayName(), co.getUsesItem())
		}
	}
	return true, ""
}

// spendUses ...
func (co *CardOption) spendUses(p power.PlayerIF) (r string) {
	if co.isUses() {
		co.usesOption.spendUses()
		r = fmt.Sprintf("[%s]消秏%d%s", p.GetDisplayName(), co.usesOption.spend, co.usesOption.usesItem)
	}
	return
}

// MakeUses ...
func (co *CardOption) MakeUses(p power.PlayerIF, n int) (bool, string) {
	if co.isUses() == true {
		co.usesOption.MakeUses(n)
		return true, fmt.Sprintf("[%s]對%s補充%d%s", p.GetDisplayName(), co.GetDisplayName(), n, co.usesItem)
	}
	return false, "無法使用"
}

func (co *CardOption) viewCardSimpleListInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s%s%s%s",
		func() (r string) {
			if co.actionTimes == 0 && baseCards[co.sourceID].card.actionTimes > 0 {
				r = lang.Lang("ActionDone")
			}
			return
		}(),
		func() string {
			if co.getEquipped() {
				return "*"
			}
			return ""
		}(),
		co.GetDisplayName(),
		func() string {
			if co.uses {
				return fmt.Sprintf("(%d)", co.quantity)
			}
			return ""
		}()))
	if co.Info.Health != 0 {
		if exists, _ := helper.InArray(CardTraitsValue.Enemy, co.CardTraits); exists {
			strs = append(strs, fmt.Sprintf("(%s %d, %s %d, %s %d)", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat, lang.Lang("Horror"), co.Info.Horror))
		} else {
			strs = append(strs, fmt.Sprintf("(%s %d, %s %d)", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat))
		}
	}
	return strings.Join(strs, "")
}

func (co *CardOption) viewCardListInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s%s%s",
		func() (r string) {
			if co.actionTimes == 0 && baseCards[co.sourceID].card.actionTimes > 0 {
				r = lang.Lang("ActionDone")
			}
			return
		}(),
		func() string {
			if exists, _ := helper.InArray(CardTraitsValue.Enemy, co.CardTraits); exists {
				return "!"
			}
			if co.getEquipped() {
				return "*"
			}
			return ""
		}(), co.GetDisplayName()))
	if len(co.CardTraits) > 0 {
		strs = append(strs, fmt.Sprintf("%s", strings.Join(
			func() (traits []string) {
				for _, trait := range co.CardTraits {
					traits = append(traits, lang.Lang(trait))
				}
				return
			}(), ",")))
	}
	if co.Info.Health != 0 {
		if exists, _ := helper.InArray(CardTraitsValue.Enemy, co.CardTraits); exists {
			strs = append(strs, fmt.Sprintf("%s %d, %s %d, %s %d", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat, lang.Lang("Horror"), co.Info.Horror))
		} else {
			strs = append(strs, fmt.Sprintf("%s %d, %s %d", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat))
		}
	}
	if co.uses {
		strs = append(strs, fmt.Sprintf("Uses: %s %d", co.usesItem, co.quantity))
	}
	if len(co.desc) > 0 {
		strs = append(strs, fmt.Sprintf("%s", co.desc))
	}
	return strings.Join(strs, "\n")
}

func (co *CardOption) viewCardFullInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s%s", func() string {
		if exists, _ := helper.InArray(CardTraitsValue.Enemy, co.CardTraits); exists {
			return "!"
		}
		if co.getEquipped() {
			return "*"
		}
		return ""
	}(), co.GetDisplayName()))
	strs = append(strs, fmt.Sprintf("*%s*", lang.Lang(co.cardType)))

	if len(co.CardTraits) > 0 {
		strs = append(strs, fmt.Sprintf("-%s-", strings.Join(
			func() (traits []string) {
				for _, trait := range co.CardTraits {
					traits = append(traits, lang.Lang(trait))
				}
				return
			}(), ",")))
	}

	if co.Info.Health != 0 {
		strs = append(strs, func() string {
			if exists, _ := helper.InArray(CardTraitsValue.Enemy, co.CardTraits); exists {
				return fmt.Sprintf("%s %d, %s %d, %s %d", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat, lang.Lang("Horror"), co.Info.Horror)
			}
			return fmt.Sprintf("%s %d, %s %d", lang.Lang("Health"), co.Info.Health, lang.Lang("Combat"), co.Info.Combat)
		}())
	}

	if co.uses {
		strs = append(strs, fmt.Sprintf("Uses: %s %d", co.usesItem, co.quantity))
	}

	strs = append(strs, fmt.Sprintf("-%s-", co.desc))
	strs = append(strs, fmt.Sprintf("Cost: %d", co.cost))
	strs = append(strs, fmt.Sprintf("Core Set #%d", co.id))
	return strings.Join(strs, "\n")
}

// GetDisplayName ...
func (co *CardOption) GetDisplayName() string {
	if s := lang.Lang(co.displayname); len(s) > 0 {
		return s
	}
	return co.displayname
}

// GetDisplayNameWithBracket ...
func (co *CardOption) GetDisplayNameWithBracket() string {
	return "<" + co.GetDisplayName() + ">"
}

// GetHealth ...
func (co *CardOption) GetHealth() int {
	return co.Health
}

// SetHealth ...
func (co *CardOption) SetHealth(n int) {
	co.Health = helper.Min(helper.Max(n, 0), co.HealthMax)
}

// MakeHealth ...
func (co *CardOption) MakeHealth(n int) (r bool) {
	x := co.GetHealth()
	co.SetHealth(co.Health + n)
	if x != co.GetHealth() {
		r = true
	}
	return
}

// GetHorror ...
func (co *CardOption) GetHorror() int {
	return co.Horror
}

// SetHorror ...
func (co *CardOption) SetHorror(n int) {
	co.Horror = helper.Max(n, 0)
}

// MakeHorror ...
func (co *CardOption) MakeHorror(n int) {
	co.SetHorror(co.Horror + n)
}

// GotHurt ...
func (co *CardOption) GotHurt(from power.FightIF, dmg power.Damage) string {
	var strs []string
	if dmg.Atk != 0 {
		co.MakeHealth(-dmg.Atk)
		strs = append(strs, fmt.Sprintf("<%s>受到傷害%d(%d)", co.GetDisplayName(), helper.Abs(dmg.Atk), co.GetHealth()))
		if co.GetHealth() <= 0 {
			strs = append(strs, fmt.Sprintf("<%s> is dead.", co.GetDisplayName()))
			co.OwnPlayer.RemoveCards(co)
			return strings.Join(strs, "\n")
		}
	}

	return strings.Join(strs, "\n")
}

// GotHeal ...
func (co *CardOption) GotHeal(from power.FightIF, dmg power.Damage) string {
	var strs []string
	if dmg.Atk != 0 {
		if co.MakeHealth(dmg.Atk) {
			strs = append(strs, fmt.Sprintf("<%s>生命%+d(%d)", co.GetDisplayName(), dmg.Atk, co.GetHealth()))
		}
		if dmg.Hor != 0 {
			co.MakeHorror(-dmg.Hor)
			strs = append(strs, fmt.Sprintf("<%s>恐懼%+d(%d)", co.GetDisplayName(), -dmg.Hor, co.GetHorror()))
		}
	}
	return strings.Join(strs, "\n")
}

// Attack ...
func (co *CardOption) Attack(target power.FightIF, dmg power.Damage) string {
	r := target.GotHurt(co, dmg)
	return r

}

// Heal ...
func (co *CardOption) Heal(target power.FightIF, dmg power.Damage) string {
	r := target.GotHeal(co, dmg)
	return r
}
