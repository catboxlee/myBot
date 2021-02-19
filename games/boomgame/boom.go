package boomgame

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"myBot/emoji"
	"myBot/helper"
	"myBot/mydb"
	"myBot/users"
	"sort"
	"strings"
)

// GameType ...
type GameType struct {
	sourceid string `db:"sourceid"`
	season   int    `db:"season"`
	scene    int    `db:"scene"`
	rank     map[string]*rankType
	data     *gameDataType
}

type gameDataType struct {
	sceneInfo interface {
		show(g *GameType) string
		startPhase(g *GameType)
		runPhase(x string, g *GameType)
		reset()
		stage(g *GameType)
	}
	players playersType
}

type rankType struct {
	UserID      string `json:"userid"`
	DisplayName string `json:"displayname"`
	Boom        int    `json:"boom"`
}

type playersType struct {
	List  map[string]playerType `json:"list"`
	Queue []playerType          `json:"queue"`
}
type playerType struct {
	UserID        string `json:"userid"`
	DisplayName   string `json:"displayname"`
	SwallowReturn int    `json:"swallowreturn"`
}

// Boom ...
var Boom = make(map[string]*GameType)
var texts []string

// Command ...
func (b *GameType) Command(input string) []string {
	texts = nil
	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		b.checkCommand(strings.TrimLeft(input, "/"))
		return texts
	}

	b.runPhase(input)
	return texts
}

func (b *GameType) checkCommand(input string) {
	switch input {
	case "reset":
		b.reset()
		b.startPhase()
		b.updateData()
	case "resetRank":
		b.resetRank()
		b.season++
		b.updateData()
	case "rank":
		b.showRank()
	case "v":
		log.Println("call showStatus()")
		b.showStatus()
	case "gacha":
		b.gaCha(1)
	case "gacha10":
		b.gaCha(10)
	}
}

func (b *GameType) recordPlayers() {
	if exist, _ := helper.InArray(playerType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName}, b.data.players.Queue); !exist {
		b.data.players.Queue = append(b.data.players.Queue, playerType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName})
	}

	if _, exist := b.data.players.List[users.LineUser.UserProfile.UserID]; !exist {
		b.data.players.List[users.LineUser.UserProfile.UserID] = playerType{UserID: users.LineUser.UserProfile.UserID, DisplayName: users.LineUser.UserProfile.DisplayName, SwallowReturn: 0}
	}
}

func (b *GameType) showStatus() {
	var str []string
	str = append(str, fmt.Sprintf("*%s*", users.LineUser.UserProfile.DisplayName))
	str = append(str, fmt.Sprintf("【資產】"))
	str = append(str, fmt.Sprintf("%s %d", emoji.Emoji(":money_bag:"), users.UsersList.Data[users.LineUser.UserProfile.UserID].Money))
	str = append(str, fmt.Sprintf("%s %d", emoji.Emoji(":gem_stone:"), users.UsersList.Data[users.LineUser.UserProfile.UserID].GemStone))
	str = append(str, fmt.Sprintf("【擁有技能】"))
	if users.UsersList.Data[users.LineUser.UserProfile.UserID].SwallowReturn > 0 {
		str = append(str, fmt.Sprintf("燕返: %d%%+%d%%", users.UsersList.Data[users.LineUser.UserProfile.UserID].SwallowReturn, b.data.players.List[users.LineUser.UserProfile.UserID].SwallowReturn))
	} else {
		str = append(str, fmt.Sprintf("無"))
	}
	str = append(str, fmt.Sprintf("【擁有稱號】"))
	str = append(str, fmt.Sprintf("無"))
	str = append(str, fmt.Sprintf("【擁有道具】"))
	str = append(str, fmt.Sprintf("無"))
	texts = append(texts, strings.Join(str, "\n"))
}

func (b *GameType) show() {
	texts = append(texts, b.data.sceneInfo.show(b))
}

func (b *GameType) stage() {
	b.data.sceneInfo.stage(b)
}

func (b *GameType) startPhase() {
	b.data.sceneInfo.startPhase(b)
}

func (b *GameType) runPhase(input string) {
	log.Println(b.data.sceneInfo)
	b.data.sceneInfo.runPhase(input, b)
}

func (b *GameType) reset() {
	/*
		boomDice := &dice.Dice
		boomDice.Roll("1d6")
		b.scene = boomDice.Hit
	*/
	b.scene = 2
	b.setSceneInfo()
	b.data.sceneInfo.reset()
}

func (b *GameType) showRank() {
	text := fmt.Sprintf("爆爆王 S%d Rank：", b.season)
	values := make([]*rankType, 0, len(b.rank))
	for _, v := range b.rank {
		values = append(values, v)
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Boom > values[j].Boom
	})
	for _, v := range values {
		text += fmt.Sprintf("\n%s %s x %d", v.DisplayName, emoji.Emoji(":collision:"), v.Boom)
	}
	texts = append(texts, text)
}

func (b *GameType) checkRank() {
	var text string
	boomKing := false
	for _, v := range b.rank {
		if v.Boom >= 100 {
			boomKing = true
			text += fmt.Sprintf("\n%s %sx%d", v.DisplayName, emoji.Emoji(":collision:"), v.Boom)
		}
	}
	if boomKing {
		texts = append(texts, fmt.Sprintf("%s S%d 爆爆王 %s%s", emoji.Emoji(":confetti_ball:"), b.season, emoji.Emoji(":confetti_ball:"), text))
		b.season++
		b.resetRank()
	}
}

func (b *GameType) resetRank() {
	b.rank = make(map[string]*rankType)
}

func (b *GameType) gaCha(n int) {
	var strs []string
	money := n * 250
	if users.UsersList.Data[users.LineUser.UserProfile.UserID].GemStone-money >= 0 {
		users.UsersList.Data[users.LineUser.UserProfile.UserID].GemStone -= money
		strs = append(strs, fmt.Sprintf("【%s】%s%d(-%d)", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":gem_stone:"), users.UsersList.Data[users.LineUser.UserProfile.UserID].GemStone, money))
		if n == 10 {
			strs = append(strs, fmt.Sprintf("【%s】轉蛋10連抽", users.LineUser.UserProfile.DisplayName))
			strs = append(strs, b.doGaCha(n))
		} else {
			strs = append(strs, fmt.Sprintf("【%s】轉蛋單抽", users.LineUser.UserProfile.DisplayName))
			strs = append(strs, b.doGaCha(n))
		}
	} else {
		strs = append(strs, fmt.Sprintf("【%s】%s不足", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":gem_stone:")))
	}
	texts = append(texts, strings.Join(strs, "\n"))
}

func (b *GameType) doGaCha(n int) string {
	var gaShaPon = []struct {
		level string
		cnt   int
	}{
		{level: "ssr", cnt: 3},
		{level: "ssr", cnt: 12},
		{level: "r", cnt: 85},
	}
	var gaShaPons []string
	var strs []string

	for _, d := range gaShaPon {
		for i := 0; i < d.cnt; i++ {
			gaShaPons = append(gaShaPons, d.level)
		}
	}
	for i := 0; i < n; i++ {
		r := rand.Perm(len(gaShaPons))[0]
		switch gaShaPons[r] {
		case "ssr":
			users.UsersList.Data[users.LineUser.UserProfile.UserID].SwallowReturn++
			strs = append(strs, fmt.Sprint("SSR - 燕返(常駐)+1%"))
		case "sr":
			b.data.players.List[users.LineUser.UserProfile.UserID].SwallowReturn += 3
			strs = append(strs, fmt.Sprint("SR - 燕返+3%"))
		case "r":
			b.data.players.List[users.LineUser.UserProfile.UserID].SwallowReturn++
			strs = append(strs, fmt.Sprint("R - 燕返+1%"))
		default:
		}
	}
	return strings.Join(strs, "\n")
}

// CheckExistData ...
func CheckExistData(SourceID string) {
	if _, exist := Boom[SourceID]; !exist {
		loadData(SourceID)
	}
}

func getSourceID() {
	return
}

func loadData(SourceID string) {
	row := mydb.Db.QueryRow("SELECT sourceid, season, scene, scene_info, players, rank FROM boom_game where sourceid = $1 limit 1", SourceID)
	var sid string
	var season int
	var scene int
	var sceneInfo json.RawMessage
	var rank json.RawMessage
	var players json.RawMessage
	switch err := row.Scan(&sid, &season, &scene, &sceneInfo, &players, &rank); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned")
		Boom[SourceID] = &GameType{}
		Boom[SourceID].sourceid = SourceID
		Boom[SourceID].season = 1
		Boom[SourceID].scene = 0
		Boom[SourceID].rank = make(map[string]*rankType)
		Boom[SourceID].data = &gameDataType{}
		Boom[SourceID].data.players.List = make(map[string]playerType)
		Boom[SourceID].data.players.Queue = nil
		Boom[SourceID].reset()
		Boom[SourceID].addData()
	case nil:
		Boom[SourceID] = &GameType{}
		Boom[SourceID].sourceid = sid
		Boom[SourceID].season = season
		Boom[SourceID].scene = scene
		Boom[SourceID].rank = make(map[string]*rankType)
		Boom[SourceID].data = &gameDataType{}
		Boom[SourceID].setSceneInfo()
		json.Unmarshal(rank, &Boom[SourceID].rank)
		json.Unmarshal(sceneInfo, &Boom[SourceID].data.sceneInfo)
		json.Unmarshal(players, &Boom[SourceID].data.players)
		Boom[SourceID].stage()
		log.Println("Boom data load.")
		//Boom[SourceID].data.updateData()
	default:
		checkError(err)
	}
}

func (b *GameType) addData() {
	rank, err := json.Marshal(b.rank)
	checkError(err)

	sceneInfo, err := json.Marshal(b.data.sceneInfo)
	checkError(err)

	players, err := json.Marshal(b.data.players)
	checkError(err)

	stmt, err := mydb.Db.Prepare("insert into boom_game (sourceid, season, scene, scene_info, players, rank) values ($1, $2, $3, $4, $5, $6)")
	checkError(err)

	_, err = stmt.Exec(b.sourceid, b.season, b.scene, sceneInfo, players, rank)
	checkError(err)

	stmt.Close()
	log.Println("Boom Data Create...")
}

func (b *GameType) updateData() {
	rank, err := json.Marshal(b.rank)
	checkError(err)

	sceneInfo, err := json.Marshal(b.data.sceneInfo)
	checkError(err)

	players, err := json.Marshal(b.data.players)
	checkError(err)

	stmt, err := mydb.Db.Prepare("update boom_game set season = $2, scene = $3, scene_info = $4, players = $5, rank = $6 where sourceid = $1")
	checkError(err)

	_, err = stmt.Exec(b.sourceid, b.season, b.scene, sceneInfo, players, rank)
	checkError(err)

	stmt.Close()
	log.Println("Boom Data Update...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
