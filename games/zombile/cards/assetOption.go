package cards

import (
	"fmt"
	"my/games/zombile/power"
)

type assetOption struct {
	//attribute []struct{}
	equipped   bool // 是否已上場
	usesOption      // 秏材
}

func (ao *assetOption) onPlay(currentPlayer power.PlayerIF, co *CardOption) (bool, string) {
	var s string

	if ao.getEquipped() == false {
		ao.setEquipped(true)
		return false, fmt.Sprintf("[%s]EQ<%s>", currentPlayer.GetDisplayName(), co.GetDisplayName())
	}

	return true, s
}

// getEquipped 是否已上場
func (ao *assetOption) getEquipped() bool {
	return ao.equipped
}

// setEquipped 標記是否已上場
func (ao *assetOption) setEquipped(b bool) {
	ao.equipped = b
}

// activate 啟動卡牌

func (ao *assetOption) viewCardInfo() (r string) {
	if ao.uses {
		r = fmt.Sprintf("Uses: %s %d", ao.usesItem, ao.quantity)
	}
	return
}
