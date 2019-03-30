package common

import (
	"regexp"
	"strings"
	"myBot/world"
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
			}
		}
		return texts
	}
	return texts
}

func changeGame(s string) {
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