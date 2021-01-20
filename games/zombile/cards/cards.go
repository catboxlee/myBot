package cards

import (
	"fmt"
	"myBot/games/zombile/power"
	"strings"
)

// Cards ...
type Cards struct {
	Card []*CardOption
}

// ClearCards ...
func (c *Cards) ClearCards() {
	c.Card = nil
}

// CreateCardDeck ...
func (c *Cards) CreateCardDeck() {
	for sourceIndex, val := range baseCards {
		for i := 0; i < val.pics; i++ {
			c.Card = append(c.Card, newCard(val.card, len(c.Card), sourceIndex))
		}
	}
}

// CreateNewCard ...
func (c *Cards) CreateNewCard(sourceIndex int, toDeck *Cards) {
	val := baseCards[sourceIndex]
	c.Card = append(toDeck.Card, newCard(val.card, len(toDeck.Card), sourceIndex))
}

// NewCard ...
func newCard(bc DefaultCardOption, idx int, sourceIndex int) *CardOption {
	nc := new(CardOption)
	nc.id = idx
	nc.sourceID = sourceIndex
	nc.cost = bc.cost
	nc.Info = bc.Info
	nc.cardName = bc.cardName
	nc.displayname = bc.displayname
	nc.cardType = bc.cardType
	nc.CardTraits = bc.CardTraits
	nc.desc = bc.desc
	nc.equipped = bc.equipped
	nc.usesOption = bc.usesOption
	nc.actionTimes = bc.actionTimes

	// 重建method
	//newasset := bc.method.(*assetOption)
	//newmethod := *newasset
	//nc.method = &newmethod
	//nc.Method = bc.Method.clone()
	if bc.ActivateFunc != nil {
		nc.ActivateFunc = bc.ActivateFunc(nc)
	}
	if bc.OnDisplayFunc != nil {
		nc.OnDisplayFunc = bc.OnDisplayFunc(nc)
	}
	if bc.OnMysterFunc != nil {
		nc.OnMysterFunc = bc.OnMysterFunc(nc)
	}
	if bc.OnHealthDamageAfterFunc != nil {
		nc.OnHealthDamageAfterFunc = bc.OnHealthDamageAfterFunc(nc)
	}
	if bc.OnHorrorDamageAfterFunc != nil {
		nc.OnHorrorDamageAfterFunc = bc.OnHorrorDamageAfterFunc(nc)
	}
	if bc.OnPlayerHealthHurtAfterFunc != nil {
		nc.OnPlayerHealthHurtAfterFunc = bc.OnPlayerHealthHurtAfterFunc(nc)
	}

	return nc
}

// GetAllCards ...
func (c *Cards) GetAllCards() []*CardOption {
	return c.Card
}

// GetCardsCount ...
func (c *Cards) GetCardsCount() int {
	return len(c.Card)
}

// PopCards ...
func (c *Cards) PopCards() (nc *CardOption) {
	if len(c.Card) > 0 {
		nc = c.Card[0]
		c.Card = c.Card[1:]
	}
	return
}

// RemoveCards ...
func (c *Cards) RemoveCards(n int) {
	if len(c.Card) > 0 {
		c.Card = append(c.Card[:n], c.Card[n+1:]...)
	}
}

// TakeCard ...
func (c *Cards) TakeCard(thisPlayer power.PlayerIF, val power.CardIF) {
	val.(*CardOption).OwnPlayer = thisPlayer
	c.Card = append(c.Card, val.(*CardOption))
}

// ViewAllCardsSimpleCheckList ...
func (c *Cards) ViewAllCardsSimpleCheckList() string {
	var strs []string
	if len(c.Card) > 0 {
		for id, val := range c.Card {
			strs = append(strs, fmt.Sprintf("(%d) %s", id, val.viewCardSimpleListInfo()))
		}
	} else {
		strs = append(strs, "Empty...")
	}
	return strings.Join(strs, "\n")
}

// ViewAllCardsCheckList ...
func (c *Cards) ViewAllCardsCheckList() string {
	var strs []string
	if len(c.Card) > 0 {
		for id, val := range c.Card {
			strs = append(strs, fmt.Sprintf("(%d) %s", id, val.viewCardListInfo()))
		}
	} else {
		strs = append(strs, "Empty...")
	}
	return strings.Join(strs, "\n")
}

// ViewCardFullInfo ...
func (c *Cards) ViewCardFullInfo(index int) (s string) {
	if index < len(c.Card) {
		s = c.Card[index].viewCardFullInfo()
	}
	return
}
