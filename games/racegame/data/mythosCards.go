package data

import (
	"fmt"
	"myBot/emoji"
	"myBot/games/racegame/scheduler"
	"strings"
)

// MythosCard ...
var MythosCard = map[string]CardOption{
	"gandalf": CardOption{
		CardName:    fmt.Sprintf("%s「」", emoji.Emoji(":ghost:")),
		DisplayName: "",
		Class:       "R",
		CoreSet:     "gandalf",
		CoolDown:    13,
		ReCoolDown:  13,
		Unique:      true,
		DescFunc: func(thisCard scheduler.Card) func() string {
			return func() string {
				str := fmt.Sprintf("CD%d", thisCard.GetReCoolDown())
				thisCard.SetDesc(str)
				return str
			}
		},
		OnMythosPassFunc: func(thisCard scheduler.Card) func(scheduler.Game) (r bool, s string) {
			return func(g scheduler.Game) (r bool, s string) {
				var strs []string
				s = strings.Join(strs, "\n")
				return r, s
			}
		},
	},
}
