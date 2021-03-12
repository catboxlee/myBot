package boomgame

import (
	"fmt"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/data"
	"myBot/users"
	"strconv"
	"strings"
)

// GaCha ...
func (g *GameType) GaCha(input string, currentID string) string {
	var strs []string
	gem1 := 250
	gem10 := 2500
	n := 0
	s := strings.Fields(input)
	thisPlayer := g.Player(currentID)

	if len(s) > 1 {
		if x, err := strconv.Atoi(s[1]); err == nil {
			n = x
		}
	}
	switch n {
	case 1:
		thisPlayer.MakeGemStone(-gem1)
		strs = append(strs, fmt.Sprintf("<<%s 單抽>>%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":gem_stone:"), thisPlayer.GetGemStone(), -gem1))
		strs = append(strs, g.doGaCha(1, currentID))
	case 10:
		thisPlayer.MakeGemStone(-gem10)
		strs = append(strs, fmt.Sprintf("<<%s 10抽>>%s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":gem_stone:"), thisPlayer.GetGemStone(), -gem1))
		strs = append(strs, g.doGaCha(10, currentID))
	default:
		strs = append(strs, g.viewGaChaInfo())
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) doGaCha(n int, currentID string) string {
	var strs []string
	var ssrData []string
	var srData []string
	var rData []string
	for _, val := range data.GachaCardData {
		switch val.Class {
		case "SSR":
			ssrData = append(ssrData, val.CoreSet)
		case "SR":
			srData = append(srData, val.CoreSet)
		case "R":
			rData = append(rData, val.CoreSet)
		}
	}
	thisPlayer := g.Player(currentID)
	isSR := false
	for i := 0; i < n; i++ {
		tmp := rand.Perm(100)[0]
		switch {
		case tmp < 3:
			lucky := rand.Perm(len(ssrData))[0]
			isSR = true
			strs = append(strs, fmt.Sprintf("%s", thisPlayer.TakeCard(ssrData[lucky])))

		case tmp < 15:
			lucky := rand.Perm(len(srData))[0]
			isSR = true
			strs = append(strs, fmt.Sprintf("%s", thisPlayer.TakeCard(srData[lucky])))
		default:
			if i == 9 && isSR == false {
				lucky := rand.Perm(len(srData))[0]
				isSR = true
				strs = append(strs, fmt.Sprintf("%s", thisPlayer.TakeCard(srData[lucky])))
			} else {
				lucky := rand.Perm(len(rData))[0]
				strs = append(strs, fmt.Sprintf("%s", thisPlayer.TakeCard(rData[lucky])))
			}
		}
	}
	users.UsersList.SaveData(thisPlayer.GetUserID())
	thisPlayer.SaveData()
	return strings.Join(strs, "\n")
}

func (g *GameType) viewGaChaInfo() string {
	var strs []string
	var ssrData []string
	var srData []string
	var rData []string
	for _, val := range data.GachaCardData {
		switch val.Class {
		case "SSR":
			ssrData = append(ssrData, val.CoreSet)
		case "SR":
			srData = append(srData, val.CoreSet)
		case "R":
			rData = append(rData, val.CoreSet)
		}
	}
	strs = append(strs, "[[轉蛋舉辦中]]")
	strs = append(strs, "SSR出現率3%")
	for _, v := range ssrData {
		strs = append(strs, fmt.Sprintf("<%s>", data.GachaCardData[v].CardName))
	}
	strs = append(strs, "SR出現率12%")
	for _, v := range srData {
		strs = append(strs, fmt.Sprintf("<%s>", data.GachaCardData[v].CardName))
	}
	strs = append(strs, "R出現率85%")
	for _, v := range rData {
		strs = append(strs, fmt.Sprintf("<%s>", data.GachaCardData[v].CardName))
	}
	strs = append(strs, "[[指令]]")
	strs = append(strs, "/gacha 1 : 單抽")
	strs = append(strs, "/gacha 10 : 10抽(1張SR以上)")
	return strings.Join(strs, "\n")
}
