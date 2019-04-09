package boomgame1

import (
	"database/sql"
	"encoding/json"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/mydb"
	"myBot/users"
	"sort"

	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// GameType ...
type GameType struct {
	current int
	data    *gameDataType
}

type gameDataType struct {
	sourceid string `db:"sourceid"`
	hit      int    `db:"hit"`
	min      int    `db:"min"`
	max      int    `db:"max"`
	season   int    `db:"season"`
	scene    int    `db:"scene"`
	rank     map[string]*rankType
	players  map[string]*playerType
}

type rankType struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"displayname"`
	Boom        int    `json:"boom"`
}

type playerType struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"displayname"`
}

// Boom ...
var Boom = make(map[string]*GameType)
var texts []string

// Run ...
func (b *GameType) Run(input string) []string {
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

func (b *GameType) checkCommand(input string) {
	switch input {
	case "reset":
		b.reset()
		b.show()
		b.data.updateData()
	case "rank":
		b.rank()
		b.show()
	case "resetRank":
		b.resetRank()
		//rank.saveRank()
		b.data.updateData()
	case "boom":
		texts = append(texts, fmt.Sprintf("%s %s自爆", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":bomb:")))
		b.gameOver()
		b.rank()
		b.checkBoomKing()
		b.reset()
		b.show()
	}
}

func (b *GameType) checkBoom(x int) {
	if x > b.data.min && x < b.data.max {
		b.current = x
		switch {
		case b.current == b.data.hit:
			b.show()
			b.gameOver()
			b.rank()
			b.checkBoomKing()
			b.reset()
			b.show()
		case b.current < b.data.hit:
			b.data.min = b.current
			b.show()
		case b.current > b.data.hit:
			b.data.max = b.current
			b.show()
		}
		b.data.updateData()
	}
}

func (b *GameType) reset() {
	boomNumber := &dice.Dice
	boomNumber.Roll("1d100")
	b.data.hit = boomNumber.Hit
	b.current = 0
	b.data.min = 0
	b.data.max = 101
	b.data.scene = 0
	b.scene("start")
}

func (b *GameType) show() {
	b.scene("show")
}

func (b *GameType) gameOver() {
	b.scene("end")
}

func (b *GameType) checkBoomKing() {
	var text string
	boomKing := false
	for _, v := range b.data.rank {
		if v.Boom >= 100 {
			boomKing = true
			text += fmt.Sprintf("\n%s %sx%d", v.DisplayName, emoji.Emoji(":boom:"), v.Boom)
		}
	}
	if boomKing {
		texts = append(texts, fmt.Sprintf("%s S%d 爆爆王 %s%s", emoji.Emoji(":confetti_ball:"), b.data.season, emoji.Emoji(":confetti_ball:"), text))
		b.data.season++
		b.resetRank()
	}
}

func (b *GameType) rank() {
	text := fmt.Sprintf("爆爆王 S%d Rank：", b.data.season)
	values := make([]*rankType, 0, len(b.data.rank))
	for _, v := range b.data.rank {
		values = append(values, v)
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Boom > values[j].Boom
	})
	for _, v := range values {
		text += fmt.Sprintf("\n%s %s x %d", v.DisplayName, emoji.Emoji(":boom:"), v.Boom)
	}
	texts = append(texts, text)
}

func (b *GameType) resetRank() {
	b.data.rank = make(map[string]*rankType)
}

// CheckExistData ...
func CheckExistData(SourceID string) {
	if _, exist := Boom[SourceID]; !exist {
		loadData(SourceID)
	}
}

func loadData(SourceID string) {
	row := mydb.Db.QueryRow("SELECT sourceid, hit, min, max, season, rank FROM boom_game where sourceid = $1 limit 1", SourceID)
	var data gameDataType
	var rank json.RawMessage
	switch err := row.Scan(&data.sourceid, &data.hit, &data.min, &data.max, &data.season, &rank); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned")
		Boom[SourceID] = &GameType{}
		Boom[SourceID].data = &gameDataType{}
		Boom[SourceID].data.sourceid = SourceID
		Boom[SourceID].data.rank = make(map[string]*rankType)
		Boom[SourceID].reset()
		Boom[SourceID].data.addData()
	case nil:
		Boom[SourceID] = &GameType{}
		Boom[SourceID].data = &gameDataType{}
		Boom[SourceID].data.sourceid = data.sourceid
		Boom[SourceID].data.hit = data.hit
		Boom[SourceID].data.min = data.min
		Boom[SourceID].data.max = data.max
		Boom[SourceID].data.season = data.season
		Boom[SourceID].data.rank = make(map[string]*rankType)
		json.Unmarshal(rank, &Boom[SourceID].data.rank)
		log.Println("Boom data load.", Boom[SourceID].data)
		//Boom[SourceID].data.updateData()
	default:
		checkError(err)
	}
}

func (b *gameDataType) addData() {
	stmt, err := mydb.Db.Prepare("insert into boom_game (sourceid, hit, min, max, season, rank) values ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Fatal(err)
	}
	rank, err := json.Marshal(b.rank)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(b.sourceid, b.hit, b.min, b.max, b.season, rank)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
	log.Println("Boom Data Create...")
}

func (b *gameDataType) updateData() {
	stmt, err := mydb.Db.Prepare("update boom_game set hit = $2, min = $3, max = $4, season = $5, rank = $6 where sourceid = $1")
	if err != nil {
		log.Fatal(err)
	}

	//b.rank[users.LineUser.UserProfile.UserID] = &rankType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName, Boom: 1}

	rank, err := json.Marshal(b.rank)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(b.sourceid, b.hit, b.min, b.max, b.season, rank)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
	log.Println("Boom Data Update...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
