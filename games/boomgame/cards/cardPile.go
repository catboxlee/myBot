package cards

import (
	"fmt"
	"myBot/emoji"
	"myBot/games/boomgame/data"
	"myBot/games/boomgame/scheduler"
	"strings"
)

// CardPile ...
type CardPile struct {
	Parent scheduler.Player       `json:"-"`
	Cards  map[string]*CardOption `json:"cards"`
}

// BuildCard ...
func (c *CardPile) BuildCard(bc data.CardOption) *CardOption {
	nc := new(CardOption)
	nc.Parent = c
	nc.Own = c.Parent
	nc.CardName = bc.CardName
	nc.DisplayName = bc.DisplayName
	nc.Class = bc.Class
	nc.Desc = bc.Desc
	nc.CoolDown = bc.ReCoolDown
	nc.ReCoolDown = bc.ReCoolDown
	nc.CoreSet = bc.CoreSet
	nc.Quantity = 1
	nc.Unique = bc.Unique

	if bc.DescFunc != nil {
		nc.DescFunc = bc.DescFunc(nc)
	}

	if bc.OnMythosFunc != nil {
		nc.OnMythosFunc = bc.OnMythosFunc(nc)
	}

	if bc.OnPlayFunc != nil {
		nc.OnPlayFunc = bc.OnPlayFunc(nc)
	}

	if bc.OnHitFunc != nil {
		nc.OnHitFunc = bc.OnHitFunc(nc)
	}

	if bc.OnPassFunc != nil {
		nc.OnPassFunc = bc.OnPassFunc(nc)
	}

	if bc.OnShieldFunc != nil {
		nc.OnShieldFunc = bc.OnShieldFunc(nc)
	}

	if bc.OnAttackFunc != nil {
		nc.OnAttackFunc = bc.OnAttackFunc(nc)
	}
	if bc.OnMythosFunc != nil {
		nc.OnMythosFunc = bc.OnMythosFunc(nc)
	}

	if bc.OnMythosPassFunc != nil {
		nc.OnMythosPassFunc = bc.OnMythosPassFunc(nc)
	}

	if bc.OnSenteFunc != nil {
		nc.OnSenteFunc = bc.OnSenteFunc(nc)
	}

	if bc.OnPlayFunc != nil {
		nc.OnPlayFunc = bc.OnPlayFunc(nc)
	}
	if bc.OnHitFunc != nil {
		nc.OnHitFunc = bc.OnHitFunc(nc)
	}

	return nc
}

// Clear ...
func (c *CardPile) Clear() {
	c.Cards = make(map[string]*CardOption)
}

// GetTopParent ...
func (c *CardPile) GetTopParent() scheduler.Game {
	return c.GetParent().GetTopParent()
}

// GetParent ...
func (c *CardPile) GetParent() scheduler.Player {
	return c.Parent
}

// SetParent ...
func (c *CardPile) SetParent(po scheduler.Player) {
	c.Parent = po
}

// GetCards ...
func (c *CardPile) GetCards() map[string]scheduler.Card {
	tmp := make(map[string]scheduler.Card)
	for k, v := range c.Cards {
		tmp[k] = v
	}
	return tmp
}

// TakeCard ...
func (c *CardPile) TakeCard(cardID string) (s string) {
	if _, exist := c.Cards[cardID]; !exist {
		c.Cards[cardID] = c.BuildCard(data.CardData[cardID])
		s = fmt.Sprintf("%s%sLv%d", c.Cards[cardID].CardName, emoji.Emoji(":NEW_button:"), c.Cards[cardID].GetLevel())
	} else {
		if c.Cards[cardID].Unique {
			if 4 >= c.Cards[cardID].GetLevel()+1 {
				c.Cards[cardID].MakeLevel(1)
				s = fmt.Sprintf("%s%sLv%d", c.Cards[cardID].CardName, emoji.Emoji(":up:"), c.Cards[cardID].GetLevel())
			} else {
				s = fmt.Sprintf("%s%s%sLv%d", c.Cards[cardID].CardName, emoji.Emoji(":Japanese_no_vacancy_button:"), emoji.Emoji(":admission_tickets:"), c.Cards[cardID].GetLevel())
				switch c.Cards[cardID].Class {
				case "SSR":
					c.TakeCard("13")
				case "SR":
					c.TakeCard("14")
				case "R":
					c.TakeCard("15")
				}
			}
		} else {
			c.Cards[cardID].Quantity++
			s = fmt.Sprintf("%s(%d)", c.Cards[cardID].CardName, c.Cards[cardID].Quantity)
		}
	}
	return
}

// UsedCard ...
func (c *CardPile) UsedCard(cardID string) {
	c.Cards[cardID].Quantity--
	if c.Cards[cardID].Quantity < 0 {
		delete(c.Cards, cardID)
	}
}

// ViewCardsInfo ...
func (c *CardPile) ViewCardsInfo(desc bool) string {
	var strs []string
	var strs1 []string
	var strs2 []string
	for coid, co := range c.Cards {
		if co.OnPlayFunc != nil {
			strs2 = append(strs2, fmt.Sprintf("[%s]%s", coid, func() string {
				if desc {
					return co.ViewCardInfoWithDesc()
				}
				return co.ViewCardInfo()
			}()))
		} else {
			strs1 = append(strs1, fmt.Sprintf("[%s]%s", coid, func() string {
				if desc {
					return co.ViewCardInfoWithDesc()
				}
				return co.ViewCardInfo()
			}()))
		}
	}
	if len(strs1) > 0 {
		strs = append(strs, "[[被動]]")
		strs = append(strs, strings.Join(strs1, "\n"))
	}
	if len(strs2) > 0 {
		strs = append(strs, "[[主動]]")
		strs = append(strs, strings.Join(strs2, "\n"))
	}
	return strings.Join(strs, "\n")
}

// CreateMythosCard ...
func CreateMythosCard() []*CardOption {
	var tmpCards []*CardOption
	for _, v := range data.MythosCard {
		tmpCards = append(tmpCards, BuildMythosCard(v))
	}
	return tmpCards
}

// BuildMythosCard ...
func BuildMythosCard(bc data.CardOption) *CardOption {
	nc := new(CardOption)
	nc.CardName = bc.CardName
	nc.DisplayName = bc.DisplayName
	nc.Class = bc.Class
	nc.CoolDown = 0
	nc.ReCoolDown = bc.ReCoolDown
	nc.Desc = bc.Desc
	nc.CoreSet = bc.CoreSet
	nc.Quantity = 1
	nc.Unique = bc.Unique

	if bc.DescFunc != nil {
		nc.DescFunc = bc.DescFunc(nc)
	}

	if bc.OnPlayFunc != nil {
		nc.OnPlayFunc = bc.OnPlayFunc(nc)
	}

	if bc.OnHitFunc != nil {
		nc.OnHitFunc = bc.OnHitFunc(nc)
	}

	if bc.OnShieldFunc != nil {
		nc.OnShieldFunc = bc.OnShieldFunc(nc)
	}

	if bc.OnAttackFunc != nil {
		nc.OnAttackFunc = bc.OnAttackFunc(nc)
	}

	if bc.OnMythosFunc != nil {
		nc.OnMythosFunc = bc.OnMythosFunc(nc)
	}

	if bc.OnMythosPassFunc != nil {
		nc.OnMythosPassFunc = bc.OnMythosPassFunc(nc)
	}

	if bc.OnSenteFunc != nil {
		nc.OnSenteFunc = bc.OnSenteFunc(nc)
	}
	if bc.OnPassFunc != nil {
		nc.OnPassFunc = bc.OnPassFunc(nc)
	}

	return nc
}
