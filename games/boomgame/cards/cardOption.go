package cards

import (
	"fmt"
	"myBot/emoji"
	"myBot/games/boomgame/data"
	"myBot/games/boomgame/scheduler"
	"myBot/helper"
	"strings"
)

// CardOption ...
type CardOption struct {
	ID               int              `json:"id"`
	CardName         string           `json:"-"`
	DisplayName      string           `json:"-"`
	Parent           scheduler.Cards  `json:"-"`
	Own              scheduler.Player `json:"-"`
	Class            string           `json:"class"`
	Level            int              `json:"level"`
	TriggerTimes     int              `json:"triggertimes"`
	UnTriggerTimes   int              `json:"untriggertimes"`
	CoolDown         int              `json:"cooldown"`
	ReCoolDown       int              `json:"-"`
	Desc             string           `json:"desc"`
	Property         `json:"property"`
	Quantity         int                                 `json:"quantity"` // 數量
	Unique           bool                                `json:"unique"`   // 唯一
	Set              string                              `json:"set"`      // 牌組
	CoreSet          string                              `json:"coreset"`
	ActivateFunc     func(scheduler.Card)                `json:"-"`
	DescFunc         func() string                       `json:"-"`
	OnPlayFunc       func() (bool, string)               `json:"-"`
	OnHitFunc        func() (bool, string)               `json:"-"`
	OnSenteFunc      func() (bool, string)               `json:"-"`
	OnPassFunc       func() (bool, string)               `json:"-"`
	OnMythosFunc     func(scheduler.Game) (bool, string) `json:"-"`
	OnMythosPassFunc func(scheduler.Game) (bool, string) `json:"-"`
	OnShieldFunc     func() (bool, string)               `json:"-"`
	OnAttackFunc     func() (bool, string)               `json:"-"`
}

// Property ...
type Property struct {
	Freeze            int `json:"freeze"`
	LuckyBuff         int `json:"luckybuff"`
	LuckyBuffDuration int `json:"luckybuffduration"`
}

// GenerateCard ...
func (co *CardOption) GenerateCard(c scheduler.Cards, po scheduler.Player, bc data.CardOption) {
	co.Parent = c
	co.Own = po
	co.Class = bc.Class
	co.CardName = bc.CardName
	co.DisplayName = bc.DisplayName
	co.CoolDown = 0
	co.ReCoolDown = bc.ReCoolDown
	co.Desc = bc.Desc
	co.Unique = bc.Unique
	co.CoreSet = bc.CoreSet

	if bc.DescFunc != nil {
		co.DescFunc = bc.DescFunc(co)
	}

	if bc.OnMythosFunc != nil {
		co.OnMythosFunc = bc.OnMythosFunc(co)
	}

	if bc.OnMythosPassFunc != nil {
		co.OnMythosPassFunc = bc.OnMythosPassFunc(co)
	}

	if bc.OnPlayFunc != nil {
		co.OnPlayFunc = bc.OnPlayFunc(co)
	}

	if bc.OnHitFunc != nil {
		co.OnHitFunc = bc.OnHitFunc(co)
	}
	if bc.OnShieldFunc != nil {
		co.OnShieldFunc = bc.OnShieldFunc(co)
	}

	if bc.OnAttackFunc != nil {
		co.OnAttackFunc = bc.OnAttackFunc(co)
	}
	if bc.OnSenteFunc != nil {
		co.OnSenteFunc = bc.OnSenteFunc(co)
	}

	if bc.OnPassFunc != nil {
		co.OnPassFunc = bc.OnPassFunc(co)
	}

}

// GetTopParent ...
func (co *CardOption) GetTopParent() scheduler.Game {
	if co.Parent == nil {
		return nil
	}
	return co.GetParent().GetTopParent()
}

// GetParent ...
func (co *CardOption) GetParent() scheduler.Cards {
	if co.Parent == nil {
		return nil
	}
	return co.Parent
}

// GetDisplayName ...
func (co *CardOption) GetDisplayName() string {
	return fmt.Sprintf("%s", co.CardName)
}

// GetLevel ...
func (co *CardOption) GetLevel() int {
	return co.Level
}

// MakeLevel ...
func (co *CardOption) MakeLevel(n int) {
	co.Level = helper.Min(4, helper.Max(0, co.Level+n))
}

// SetDesc ...
func (co *CardOption) SetDesc(s string) {
	co.Desc = s
}

// GetCoolDown ...
func (co *CardOption) GetCoolDown() int {
	return co.CoolDown
}

// MakeCoolDown ...
func (co *CardOption) MakeCoolDown(n int) {
	co.CoolDown = helper.Max(co.CoolDown+n, 0)
}

// GetReCoolDown ...
func (co *CardOption) GetReCoolDown() int {
	return co.ReCoolDown
}

// ResetCoolDown ...
func (co *CardOption) ResetCoolDown() {
	co.CoolDown = co.ReCoolDown
}

// SetCoolDown ...
func (co *CardOption) SetCoolDown(n int) {
	co.CoolDown = n
}

// DoCoolDown ...
func (co *CardOption) DoCoolDown() {
	if co.CoolDown > 0 {
		co.CoolDown--
	}
}

// GetFreeze ...
func (co *CardOption) GetFreeze() int {
	return co.Property.Freeze
}

// MakeFreeze ...
func (co *CardOption) MakeFreeze(n int) {
	co.Property.Freeze = helper.Max(0, co.Property.Freeze+n)
}

// DoFreeze ...
func (co *CardOption) DoFreeze() {
	if co.Property.Freeze > 0 {
		co.Property.Freeze--
	}
}

// GetLuckyBuff ...
func (co *CardOption) GetLuckyBuff() int {
	return co.Property.LuckyBuff
}

// MakeLuckyBuff ...
func (co *CardOption) MakeLuckyBuff(n int) {
	co.Property.LuckyBuff = helper.Max(0, co.Property.LuckyBuff+n)
}

// GetLuckyBuffDuration ...
func (co *CardOption) GetLuckyBuffDuration() int {
	return co.Property.LuckyBuffDuration
}

// ViewCardInfo ...
func (co *CardOption) ViewCardInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s%s%s\n-%s", func() string {
		if co.GetCoolDown() > 0 {
			return fmt.Sprintf("(%s%d)", emoji.Emoji(":hourglass_not_done:"), co.GetCoolDown())
		}
		return ""
	}(), co.GetDisplayName(), func() string {
		if co.Unique {
			return fmt.Sprintf("LV%d", co.GetLevel())
		}
		return fmt.Sprintf("(%d)", co.Quantity)
	}(), func() string {
		if co.DescFunc != nil {
			return co.DescFunc()
		}
		return ""
	}()))
	return strings.Join(strs, "\n")
}
