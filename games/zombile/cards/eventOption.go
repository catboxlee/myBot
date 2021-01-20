package cards

import "myBot/games/zombile/power"

type eventOption struct {
	//attribute []struct{}
	equipped bool
}

func (eo *eventOption) onPlay(p power.PlayerIF, co *CardOption) (bool, string) {
	return false, ""
}

func (eo *eventOption) activate() (r string, e bool) {
	return
}
func (eo *eventOption) viewCardInfo() (r string) {
	//r = fmt.Sprintf("Uses: %s %d", ao.usesItem, ao.quantity)
	return
}
