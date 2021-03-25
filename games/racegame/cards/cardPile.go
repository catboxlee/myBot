package cards

import (
	"fmt"
	"log"
	"math/rand"
	"myBot/games/racegame/data"
	"myBot/games/racegame/scheduler"
	"strings"
)

// CardPile ...
type CardPile struct {
	Parent scheduler.Game         `json:"-"`
	Cards  map[string]*CardOption `json:"cards"`
	Deck   []string
}

// BuildCard ...
func (c *CardPile) BuildCard(bc data.CardOption) *CardOption {
	nc := new(CardOption)
	nc.Parent = c
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

	if bc.OnSpeedLimitFunc != nil {
		nc.OnSpeedLimitFunc = bc.OnSpeedLimitFunc(nc)
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

	if bc.OnMythosPassFunc != nil {
		nc.OnMythosPassFunc = bc.OnMythosPassFunc(nc)
	}

	if bc.OnEffectFunc != nil {
		nc.OnEffectFunc = bc.OnEffectFunc(nc)
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
	return c.GetParent()
}

// GetParent ...
func (c *CardPile) GetParent() scheduler.Game {
	return c.Parent
}

// SetParent ...
func (c *CardPile) SetParent(g scheduler.Game) {
	c.Parent = g
}

// GetCards ...
func (c *CardPile) GetCards() map[string]scheduler.Card {
	tmp := make(map[string]scheduler.Card)
	for k, v := range c.Cards {
		tmp[k] = v
	}
	return tmp
}

// PopCard ...
func (c *CardPile) PopCard() string {
	coid := c.Deck[0]
	c.Deck = c.Deck[1:]
	return coid
}

// UsedCard ...
func (c *CardPile) UsedCard(cardID string) {
	c.Cards[cardID].Quantity--
	if c.Cards[cardID].Quantity < 0 {
		delete(c.Cards, cardID)
	}
}

// ViewCardsInfo ...
func (c *CardPile) ViewCardsInfo() string {
	var strs []string
	for coid, co := range c.Cards {
		if co.OnPlayFunc != nil {
			log.Println(co)
			strs = append(strs, fmt.Sprintf("[%s]%s", coid, co.ViewCardInfo()))
		}
	}
	return strings.Join(strs, "\n")
}

// CreateMythosCard ...
func (c *CardPile) CreateMythosCard() {
	c.Cards = make(map[string]*CardOption)
	c.Deck = nil
	var tmp []string
	for id, v := range data.CardData {
		c.Cards[id] = BuildMythosCard(v)
		if v.Class != "SSR" {
			for i := 0; i < v.Quantity; i++ {
				tmp = append(tmp, id)
			}
		}
	}
	rnd := rand.Perm(len(tmp))
	for _, v := range rnd {
		c.Deck = append(c.Deck, tmp[v])
	}

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
	nc.Quantity = bc.Quantity
	nc.Unique = bc.Unique

	if bc.DescFunc != nil {
		nc.DescFunc = bc.DescFunc(nc)
	}

	if bc.OnPlayFunc != nil {
		nc.OnPlayFunc = bc.OnPlayFunc(nc)
	}

	if bc.OnAttackFunc != nil {
		nc.OnAttackFunc = bc.OnAttackFunc(nc)
	}
	if bc.OnEffectFunc != nil {
		nc.OnEffectFunc = bc.OnEffectFunc(nc)
	}
	if bc.OnSpeedLimitFunc != nil {
		nc.OnSpeedLimitFunc = bc.OnSpeedLimitFunc(nc)
	}

	if bc.OnHitFunc != nil {
		nc.OnHitFunc = bc.OnHitFunc(nc)
	}

	if bc.OnShieldFunc != nil {
		nc.OnShieldFunc = bc.OnShieldFunc(nc)
	}

	if bc.OnMythosPassFunc != nil {
		nc.OnMythosPassFunc = bc.OnMythosPassFunc(nc)
	}

	if bc.OnPassFunc != nil {
		nc.OnPassFunc = bc.OnPassFunc(nc)
	}

	return nc
}
