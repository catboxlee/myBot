package boomgame1

import (
	"fmt"
	"myBot/emoji"
	"myBot/helper"
	"myBot/users"
	"strings"
)

func (b *GameType) scene(phase string) {
	switch b.data.scene {
	case 1:
		b.scene1(phase)
	default:
		b.scene0(phase)
	}
}

func (b *GameType) scene0(phase string) {
	switch phase {
	case "start":
	case "show":
		if b.current == b.data.hit {
			texts = append(texts, fmt.Sprintf("%s %s %d", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":umbrella:"), b.data.hit))
		} else {
			texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.data.min), emoji.Emoji(":closed_umbrella:"), helper.Min(100, b.data.max)))
		}
	case "end":
		if _, exist := b.data.rank[users.LineUser.UserProfile.UserID]; exist {
			b.data.rank[users.LineUser.UserProfile.UserID].Boom++
		} else {
			b.data.rank[users.LineUser.UserProfile.UserID] = &rankType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName, Boom: 1}
		}
	default:
	}
}

func (b *GameType) scene1(phase string) {
	switch phase {
	case "start":
		texts = append(texts, fmt.Sprintf("%s%s%s 你在3萬5千英呎的高空上，一顆隱藏的炸彈即將引爆，找到唯一的降落傘。", emoji.Emoji(":balloon:"), emoji.Emoji(":house:"), emoji.Emoji(":balloon:")))
	case "show":
		if b.current == b.data.hit {
			var str []string
			for _, u := range b.data.players {
				if u.UserID == users.LineUser.UserProfile.UserID {
					str = append(str, fmt.Sprintf("%s %s %d", u.DisplayName, emoji.Emoji(":umbrella:"), b.data.hit))
				} else {
					str = append(str, fmt.Sprintf("%s %s", u.DisplayName, emoji.Emoji(":boom:")))
				}
			}
			strings.Join(str, "\n")
			texts = append(texts, strings.Join(str, "\n"))
		} else {
			texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.data.min), emoji.Emoji(":closed_umbrella:"), helper.Min(100, b.data.max)))
		}
	case "end":
		for _, u := range b.data.players {
			if u.UserID != users.LineUser.UserProfile.UserID {
				if _, exist := b.data.rank[u.UserID]; exist {
					b.data.rank[u.UserID].Boom++
				} else {
					b.data.rank[u.UserID] = &rankType{UserID: u.UserID, DisplayName: u.DisplayName, Boom: 1}
				}
			}
		}
	default:
	}
}
