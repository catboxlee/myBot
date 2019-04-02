package common

import (
	"fmt"
	"log"
	"myBot/emoji"
	"myBot/mydb"
	"myBot/users"
	"myBot/world"
	"regexp"
	"strconv"
	"strings"
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
			case "bonusday":
				texts = bonusDay(matches[2])
			}
		}
		return texts
	}
	return texts
}

func changeGame(s string) []string {
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
		case 3:
			world.World.Game = n
			world.World.SaveWorldData()
			return []string{"切換遊戲：[3]射龍門祥師版"}
		case 4:
			world.World.Game = n
			world.World.SaveWorldData()
			return []string{"切換遊戲：[4]骰子PK測試版"}
		}
	}
	return []string{"[1]終極密碼\n[2]射龍門\n[3]射龍門祥師版\n[4]骰子PK測試版"}
}
func setBank(s string) []string {
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
	return []string{fmt.Sprintf("%s 銀行\n%s%d\n[1]借款\n[2]還款", emoji.Emoji(":bank:"), emoji.Emoji(":money_bag:"), world.World.Bank)}
}

func bonusDay(s string) []string {

	if n, err := strconv.Atoi(s); err == nil {
		stmt, err := mydb.Db.Prepare("update users set money = $1 where money < $1")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec(n)
		if err != nil {
			log.Fatal(err)
		}
		stmt.Close()
		for i := range users.UsersList.Data {
			if users.UsersList.Data[i].Money < n {
				users.UsersList.Data[i].Money = n
			}
		}
		return []string{"ok"}
	}
	return []string{"Input Error"}
}
