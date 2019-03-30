package common

import (
	"regexp"
	"strings"
	"strconv"
	"fmt"
	"myBot/world"
	"myBot/emoji"
)

// Cmd ...
func Cmd(input string) []string {

	var texts []string
	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		cmdLine := strings.TrimSpace(strings.TrimLeft(input, "/"))

		re := regexp.MustCompile(`(^\w+)\s?(.*)`)
		matches := re.FindStringSubmatch(cmdLine)

		if len(matches) > 1 {
			switch matches[1] {
				case "game":
					texts = changeGame(matches[2])
				case "bank":
					texts = setBank(matches[2])
			}
		}
		return texts
	}
	return texts
}

func changeGame(s string) []string{
	if n, err := strconv.Atoi(s); err == nil {
		switch n {
		case 1:
			world.World.Game = n
			world.World.SaveWorldData()
			return []string{"切換遊戲：[1]終極密碼"}
		case 2:
			world.World.Game = n
			world.World.SaveWorldData()
			return []string{"切換遊戲：[2]射龍門"}
		}
	}
	return []string{"[1]終極密碼\n[2]射龍門"}
}
func setBank(s string) []string{
	re := regexp.MustCompile(`(^\w+)\s?(.*)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		switch matches[1] {
		case "set":
			if n, err := strconv.Atoi(matches[2]); err == nil {
				world.World.Bank = n
				world.World.SaveWorldData()
				//pokergoal.Pokergoal.pot = n
				return []string{"ok"}
			}
		}
	}
	return []string{fmt.Sprintf("%s 銀行\n[1]借款\n[2]還款", emoji.Emoji(":bank:"))}
}