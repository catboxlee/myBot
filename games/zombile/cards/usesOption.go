package cards

import "fmt"

type usesOption struct {
	//supplies string
	uses     bool
	usesItem string // ammo, supplies
	quantity int
	spend    int
}

func (uo *usesOption) isUses() bool {
	return uo.uses
}

func (uo *usesOption) getQuantity() int {
	return uo.quantity
}

func (uo *usesOption) getUsesItem() string {
	return uo.usesItem
}

func (uo *usesOption) spendUses() {
	uo.quantity -= uo.spend
}

func (uo *usesOption) checkUses() (bool, string) {
	if uo.isUses() {
		if 0 > uo.quantity-uo.spend {
			return false, fmt.Sprintf("%s不足", uo.getUsesItem())
		}
	}
	return true, ""
}

func (uo *usesOption) MakeUses(n int) {
	uo.quantity += n
}
