package power

import (
	"strings"
)

// Powers ...
type Powers struct {
	GameIF
}

// GameIF ...
type GameIF interface {
	GetPlayer(string) PlayerIF
	DisCards(string, int)
	GetPlayersSequence() []string
	MoveCards(PlayerIF, CardIF, PlayerIF) string
}

// PlayerIF ...
type PlayerIF interface {
	GetUserID() string
	GetDisplayName() string
	GetHealth() int
	GetHorror() int
	SetHealth(int)
	MakeHealth(int) bool
	MakeHorror(int)
	GotHurt(FightIF, Damage) string
	GotHeal(FightIF, Damage) string
	GotHorror(int) string
	GetCard(int) CardIF
	GetCardsCount() int
	TakeCard(PlayerIF, CardIF)
	RemoveCards(interface{})
	Attack(FightIF, Damage) string
	Heal(FightIF, Damage) string
	GetDisplayNameWithBracket() string
}

// CardIF ...
type CardIF interface {
	GetDisplayName() string
	GetDisplayNameWithBracket() string
	Attack(FightIF, Damage) string
	Heal(FightIF, Damage) string
	GotHurt(FightIF, Damage) string
	GotHeal(FightIF, Damage) string
}

// Damage ...
type Damage struct {
	Atk           int
	Hor           int
	DamageSuccess func() string
	DamageFunc    map[string]func() string
}

// FightIF ...
type FightIF interface {
	GetDisplayName() string
	GetDisplayNameWithBracket() string
	Attack(FightIF, Damage) string
	Heal(FightIF, Damage) string
	GotHurt(FightIF, Damage) string
	GotHeal(FightIF, Damage) string
}

// Power ...
var Power = new(Powers)

// NewPower ...
func NewPower(g GameIF) GameIF {
	return g
}

// Fight ...
func (p *Powers) Fight(current FightIF, targetPlayer FightIF, dmg Damage) string {
	var strs []string
	if s := targetPlayer.GotHurt(current, dmg); len(s) > 0 {
		strs = append(strs, s)
	}
	return strings.Join(strs, "\n")
}
