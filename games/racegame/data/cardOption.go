package data

import (
	"myBot/games/racegame/scheduler"
)

// CardOption ...
type CardOption struct {
	Cost             int // 成本
	CardName         string
	DisplayName      string
	Class            string
	Level            int
	TriggerTimes     int
	UnTriggerTimes   int
	CoolDown         int
	ReCoolDown       int
	Desc             string
	Set              string // 牌組
	CoreSet          string
	Unique           bool
	Quantity         int
	DescFunc         func(scheduler.Card) func() string
	OnPlayFunc       func(scheduler.Card) func(scheduler.Player, scheduler.Player) (bool, string)
	OnAttackFunc     func(scheduler.Card) func(scheduler.Player) (bool, string)
	OnEffectFunc     func(scheduler.Card) func(scheduler.Player, bool) (bool, string)
	OnSpeedLimitFunc func(scheduler.Card) func(scheduler.Player) (bool, string)
	OnMythosPassFunc func(scheduler.Card) func(scheduler.Game) (bool, string)
	OnHitFunc        func(scheduler.Card) func() (bool, string)
	OnPassFunc       func(scheduler.Card) func() (bool, string)
	OnShieldFunc     func(scheduler.Card) func() (bool, string)
}

// CardStatus ...
type CardStatus struct {
	CanDrying        bool
	IsDrying         bool
	CanDeadBody      bool // corpse
	IsDead           bool
	IsIronWall       bool
	IsShield         int
	CanCounterAttack bool

	// Ability
	Vampire int // suckBlood
}
