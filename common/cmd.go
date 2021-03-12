package common

import (
	"myBot/world"
	"regexp"
	"strconv"
	"strings"
)

// Cmd ...
func Cmd(sourceID string, input string) []string {

	var texts []string
	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		cmdLine := strings.TrimSpace(strings.TrimLeft(input, "/"))

		re := regexp.MustCompile(`(^\w+)\s?(.*)`)
		matches := re.FindStringSubmatch(cmdLine)

		if len(matches) > 1 {
			switch matches[1] {
			case "game":
				texts = changeGame(sourceID, matches[2])
			}
		}
		return texts
	}
	return texts
}

func changeGame(sourceID string, s string) []string {
	if n, err := strconv.Atoi(s); err == nil {
		switch n {
		default:
			world.ConfigsData[sourceID].Game = n
			world.ConfigsData[sourceID].UpdateConfigData()
			return []string{"切換遊戲：[1]終極密碼3.0"}
		}
	}
	return []string{"[1]終極密碼3.0"}
}
