package lang

import "my/emoji"

var langCodeMap = map[string]string{
	"boom":       "Boom",
	"Combat":     emoji.Emoji(":fist_oncoming:"),
	"Health":     emoji.Emoji(":anatomical_heart:"),
	"Horror":     emoji.Emoji(":brain:"),
	"ActionDone": emoji.Emoji(":hourglass_done:"),
	"Asset":      "支援",
	"Event":      "事件",
	"Enemy":      "敵人",
	"Weakness":   "弱點",
	"Ally":       "伙伴",
	"Item":       "物品",
	"Weapon":     "武器",
	"Skill":      "技能",
	"Gun":        "槍械",
	//
	".45 Automate": ".45自動手槍",
	"Guard dog":    "護衛犬",
	"Bandage":      "繃帶",
	"Steal":        "順手牽羊",
}

// Lang ...
func Lang(input string) (s string) {
	if _, exist := langCodeMap[input]; exist {
		s = langCodeMap[input]
	}
	return
}
