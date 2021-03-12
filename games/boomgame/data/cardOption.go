package data

import "myBot/games/boomgame/scheduler"

// CardOption ...
type CardOption struct {
	Cost             int // 成本
	CardName         string
	DisplayName      string
	Own              scheduler.Player
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
	DescFunc         func(scheduler.Card) func() string
	OnMythosFunc     func(scheduler.Card) func(scheduler.Game) (bool, string)
	OnMythosPassFunc func(scheduler.Card) func(scheduler.Game) (bool, string)
	OnPlayFunc       func(scheduler.Card) func() (bool, string)
	OnHitFunc        func(scheduler.Card) func() (bool, string)
	OnSenteFunc      func(scheduler.Card) func() (bool, string)
	OnPassFunc       func(scheduler.Card) func() (bool, string)
	OnShieldFunc     func(scheduler.Card) func() (bool, string)
	OnAttackFunc     func(scheduler.Card) func() (bool, string)
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
