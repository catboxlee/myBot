package players

import (
	"fmt"
	"my/games/zombile/cards"
	"my/games/zombile/lang"
	"my/games/zombile/power"
	"my/helper"
	"strings"
)

// PlayerOption ...
type PlayerOption struct {
	UserID       string
	DisplayName  string
	actionTimes  int
	attackTimes  int
	health       int
	healthMax    int
	horror       int
	horrorMax    int
	combat       int
	attackDamage int
	cards.Cards
}

// ViewPlayerSimpleInfo ...
func (po *PlayerOption) ViewPlayerSimpleInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("```%s```", po.DisplayName))
	strs = append(strs, fmt.Sprintf("%s %d/%d, %s %d/%d\n%s %d", lang.Lang("Health"), po.health, defaultHealth, lang.Lang("Horror"), po.horror, defaultHorror, lang.Lang("Combat"), po.combat))
	strs = append(strs, fmt.Sprintf("%s", po.Cards.ViewAllCardsSimpleCheckList()))
	return strings.Join(strs, "\n")
}

// ViewPlayerInfo ...
func (po *PlayerOption) ViewPlayerInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("```%s```", po.DisplayName))
	strs = append(strs, fmt.Sprintf("%s %d/%d, %s %d/%d\n%s %d", lang.Lang("Health"), po.health, defaultHealth, lang.Lang("Horror"), po.horror, defaultHorror, lang.Lang("Combat"), po.combat))
	strs = append(strs, fmt.Sprintf("%s", po.Cards.ViewAllCardsCheckList()))
	return strings.Join(strs, "\n")
}

// GetUserID ...
func (po *PlayerOption) GetUserID() string {
	return po.UserID
}

// GetDisplayName ...
func (po *PlayerOption) GetDisplayName() string {
	return po.DisplayName
}

// GetDisplayNameWithBracket ...
func (po *PlayerOption) GetDisplayNameWithBracket() string {
	return "[" + po.GetDisplayName() + "]"
}

// GetHealth ...
func (po *PlayerOption) GetHealth() int {
	return po.health
}

// SetHealth ...
func (po *PlayerOption) SetHealth(n int) {
	po.health = helper.Max(helper.Min(n, po.healthMax), 0)
}

// MakeHealth ...
func (po *PlayerOption) MakeHealth(n int) (r bool) {
	x := po.GetHealth()
	po.SetHealth(po.health + n)
	if x != po.GetHealth() {
		r = true
	}
	return
}

// GotHurt ...
func (po *PlayerOption) GotHurt(from power.FightIF, dmg power.Damage) string {
	var strs []string
	if dmg.Atk != 0 {
		if po.MakeHealth(-dmg.Atk) {
			strs = append(strs, fmt.Sprintf("[%s]生命%+d(%d)", po.GetDisplayName(), -dmg.Atk, po.GetHealth()))
			switch from.(type) {
			case *cards.CardOption:
				if from.(*cards.CardOption).OnHealthDamageAfterFunc != nil {
					if s := from.(*cards.CardOption).OnHealthDamageAfterFunc(); len(s) > 0 {
						strs = append(strs, s)
					}
				}
			}
			for _, val := range po.GetAllCards() {
				if val.OnPlayerHealthHurtAfterFunc != nil {
					if s := val.OnPlayerHealthHurtAfterFunc(from); len(s) > 0 {
						strs = append(strs, s)
					}
				}
			}
		} else {
			strs = append(strs, fmt.Sprintf("[%s]什麼事都沒有", po.GetDisplayName()))
		}
	}
	if dmg.Hor != 0 {
		po.MakeHorror(dmg.Hor)
		strs = append(strs, fmt.Sprintf("[%s]恐懼%+d(%d)", po.GetDisplayName(), dmg.Hor, po.GetHorror()))
		for _, val := range po.GetAllCards() {
			if val.OnHorrorDamageAfterFunc != nil {
				if s := val.OnHorrorDamageAfterFunc(); len(s) > 0 {
					strs = append(strs, s)
				}
			}
		}
	}
	if po.GetHealth() <= 0 {
		strs = append(strs, fmt.Sprintf("%s is dead.", po.GetDisplayName()))
		return strings.Join(strs, "\n")
	}
	if po.GetHorror() >= defaultHorror {
		strs = append(strs, fmt.Sprintf("%s失去了理智.", po.GetDisplayName()))
	}
	return strings.Join(strs, "\n")
}

// GotHeal ...
func (po *PlayerOption) GotHeal(from power.FightIF, dmg power.Damage) string {
	var strs []string
	if dmg.Atk != 0 {
		if po.MakeHealth(dmg.Atk) {
			strs = append(strs, fmt.Sprintf("[%s]生命%+d(%d)", po.GetDisplayName(), dmg.Atk, po.GetHealth()))
		}
		if dmg.Hor != 0 {
			po.MakeHorror(-dmg.Hor)
			strs = append(strs, fmt.Sprintf("[%s]恐懼%+d(%d)", po.GetDisplayName(), -dmg.Hor, po.GetHorror()))
		}
	}
	return strings.Join(strs, "\n")
}

// GetHorror ...
func (po *PlayerOption) GetHorror() int {
	return po.horror
}

// SetHorror ...
func (po *PlayerOption) SetHorror(n int) {
	po.horror = helper.Min(helper.Max(n, 0), po.horrorMax)
}

// MakeHorror ...
func (po *PlayerOption) MakeHorror(n int) {
	po.SetHorror(po.horror + n)
}

// GotHorror ...
func (po *PlayerOption) GotHorror(n int) (r string) {
	po.MakeHorror(n)
	r = fmt.Sprintf("[%s]受到恐懼%d(%d)", po.GetDisplayName(), n, po.GetHorror())
	if po.GetHorror() >= defaultHorror {
		r += fmt.Sprintf("\n%s失去了理智.", po.GetDisplayName())
	}
	return
}

// GetCard ...
func (po *PlayerOption) GetCard(n int) power.CardIF {
	if n < len(po.Card) {
		return po.Card[n]
	}
	return nil
}

// RemoveCards ...
func (po *PlayerOption) RemoveCards(n interface{}) {
	switch n.(type) {
	case int:
		po.Cards.RemoveCards(n.(int))
	default:
		if exists, i := helper.InArray(n, po.Card); exists {
			po.Cards.RemoveCards(i)
		}
	}
}

// Attack ...
func (po *PlayerOption) Attack(target power.FightIF, dmg power.Damage) string {
	//dmg.Atk += po.attackDamage
	r := target.GotHurt(po, dmg)
	return r
}

// Heal ...
func (po *PlayerOption) Heal(target power.FightIF, dmg power.Damage) string {
	r := target.GotHeal(po, dmg)
	return r
}
