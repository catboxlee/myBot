package boomgame1

import (
	"database/sql"
	"fmt"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/mydb"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
)

type gameType struct {
	hit      int
	current  int
	min      int
	max      int
	season   int
	sourceID string
	//rank     rankType
}

type rankType struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Boom        int    `json:"boom"`
}

// Boom ...
var Boom gameType
var texts []string

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Run ...
func (b *gameType) Run(input string) []string {
	if b.hit == 0 {
		b.reset()
		//fmt.Println(Boom.hit)
	}
	if b.season == 0 {
		b.getInfo()
		//fmt.Println(Boom)
	}
	texts = nil

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		b.checkCommand(strings.TrimLeft(input, "/"))
		return texts
	}

	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		if x, err := strconv.Atoi(matches[1]); err == nil {
			// 數字 - 檢查炸彈
			b.checkBoom(x)
			return texts
		}
	}
	return texts
}

func (b *gameType) checkCommand(input string) {
	switch input {
	case "reset":
		b.reset()
		b.show()
	case "rank":
		b.rank()
		b.show()
	case "resetRank":
		b.resetRank()
		//rank.saveRank()
	}
}

func (b *gameType) getInfo() {
	row := mydb.Db.QueryRow("select season from boom_info limit 1")
	var season int
	switch err := row.Scan(&season); err {
	case sql.ErrNoRows:
		//fmt.Println("No rows were returned")
	case nil:
		b.season = season
	default:
		checkError(err)
	}
}

func (b *gameType) checkBoom(x int) {
	if x > b.min && x < b.max {
		b.current = x
		switch {
		case b.current == b.hit:
			b.show()
			b.addUserBoom()
			b.rank()
			b.checkBoomKing()
			b.reset()
			b.show()
		case b.current < b.hit:
			b.min = b.current
			b.show()
		case b.current > b.hit:
			b.max = b.current
			b.show()
		}
	}
}

func (b *gameType) reset() {
	hit := &dice.Dice
	hit.Roll("1d100")
	b.hit = hit.N
	b.current = 0
	b.min = 0
	b.max = 101
}

func (b *gameType) show() {
	if b.current == b.hit {
		texts = append(texts, fmt.Sprintf("%s %s %d", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":collision:"), b.hit))
	} else {
		texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.min), emoji.Emoji(":bomb:"), helper.Min(100, b.max)))
	}
}

func (b *gameType) addUserBoom() {

	query := `insert into boom_rank(userid, displayname, boom) values($1, $2, 1)
					on conflict(userid)
					do update set displayname = $2, boom = boom_rank.boom + 1`
	mydb.Db.QueryRow(query, users.LineUser.UserProfile.UserID, users.LineUser.UserProfile.DisplayName)
}

func (b *gameType) checkBoomKing() {
	var r rankType
	row := mydb.Db.QueryRow("select userid, displayname, boom from boom_rank where boom >= 100 limit 1")
	switch err := row.Scan(&r.UserID, &r.DisplayName, &r.Boom); err {
	case sql.ErrNoRows:
		//fmt.Println("No rows were returned")
	case nil:
		texts = append(texts, fmt.Sprintf("%s S%d 爆爆王：%s %s", emoji.Emoji(":confetti_ball:"), b.season, r.DisplayName, emoji.Emoji(":confetti_ball:")))

		mydb.Db.QueryRow("truncate table boom_rank")
		b.season++
		mydb.Db.QueryRow("update boom_info set season = $1", b.season)
	default:
		checkError(err)
	}
}

func (b *gameType) rank() {

	text := fmt.Sprintf("爆爆王 S%d Rank：", Boom.season)
	rows, err := mydb.Db.Query("SELECT userid, displayname, boom FROM boom_rank order by boom desc")
	checkError(err)
	defer rows.Close()
	for rows.Next() {
		var r rankType
		switch err := rows.Scan(&r.UserID, &r.DisplayName, &r.Boom); err {
		case sql.ErrNoRows:
			//fmt.Println("No rows were returned")
		case nil:
			text += fmt.Sprintf("\n%s %s x %d", r.DisplayName, emoji.Emoji(":boom:"), r.Boom)
		default:
			checkError(err)
		}
	}
	texts = append(texts, text)
}

func (r *rankType) resetRank() {

}
