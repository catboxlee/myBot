package boomgame

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"myBot/emoji"
	"myBot/games/boomgame/cards"
	"myBot/games/boomgame/data"
	"myBot/games/boomgame/players"
	"myBot/games/boomgame/scheduler"
	"myBot/helper"
	"myBot/mydb"
	"strings"
)

// GameType ..
type GameType struct {
	players.Players
	sourceid    string `db:"sourceid"`
	season      int    `db:"season"`
	rank        map[string]*rankType
	Info        *InfoType           `json:"scene_info"`
	mythosCards []*cards.CardOption `json:"-"`
}

// InfoType ...
type InfoType struct {
	Hit           int      `json:"hit"`
	Current       int      `json:"current"`
	Min           int      `json:"min"`
	Max           int      `json:"max"`
	Turn          int      `json:"turn"`
	Inning        int      `json:"inning"`
	BoomCnt       int      `json:"boomcnt"`
	CurrentUserID string   `json:"currentuserid"`
	Queue         []string `json:"queue"`
	PlayQueue     []string `json:"playqueue"`
}

type rankType struct {
	UserID    string `json:"userid"`
	Boom      int    `json:"boom"`
	PlayTimes int    `json:"playertimes"`
	WinTimes  int    `json:"wintimes"`
}

const (
	seasonBoomCount    = 100
	singleGameBonusGem = 25
	seasonGameBonusGem = 390
)

// Boom ...
var Boom = make(map[string]*GameType)
var texts []string

// Start ...
func (g *GameType) Start() {
	log.Println("<Boom.Start()>")
	g.reset()
}

func (g *GameType) startPhase() {
	var strs []string
	strs = append(strs, fmt.Sprintf("[%s終極密碼3.0]", emoji.Emoji(":bomb:")))
	strs = append(strs, g.showGameInfo())
	texts = append(texts, strings.Join(strs, "\n"))
}

// Command ...
func (g *GameType) Command(input string, currentID string) []string {
	texts = nil
	g.Info.CurrentUserID = currentID

	if len(input) > 0 {
		g.Players.CheckPlayerExist(g.Info.CurrentUserID)

		// 去除前後空格
		input = strings.Trim(input, " ")

		// 檢查字首
		switch input[0:1] {
		case "/": // Game Command
			g.checkCommand(strings.TrimLeft(input, "/"), currentID)
		default: // Player Phase
			g.GamePhase(input)
		}
	}
	log.Println("texts", texts)
	return texts
}

func (g *GameType) checkCommand(input string, currentID string) (r string) {
	s := strings.Fields(input)
	input = strings.ToLower(input)
	if len(s) > 0 {
		switch s[0] {
		case "reset":
			log.Println(fmt.Sprintf("Command<reset>: %s, %s", s[0], input))
			g.reset()
			g.startPhase()
		case "rank":
			log.Println(fmt.Sprintf("Command<reset>: %s, %s", s[0], input))
			g.showRank()
			g.Show()
		case "resetrank":
			log.Println(fmt.Sprintf("Command<resetRank>: %s, %s", s[0], input))
			g.resetRank()
		case "v":
			log.Println(fmt.Sprintf("Command<v>: %s, %s", s[0], input))
			texts = append(texts, g.Players.Player(currentID).ViewInfo())
		case "u":
			log.Println(fmt.Sprintf("Command<v>: %s, %s", s[0], input))
			g.PlayerCard(input, currentID)
		case "gacha":
			log.Println(fmt.Sprintf("Command<v>: %s, %s", s[0], input))
			texts = append(texts, g.GaCha(input, currentID))
		case "listcard":
			log.Println(fmt.Sprintf("Command<v>: %s, %s", s[0], input))
			var send []string
			for k, v := range data.CardData {
				send = append(send, fmt.Sprintf("%s,%s,%s", k, v.CoreSet, v.CardName))
			}
			texts = append(texts, strings.Join(send, "\n"))
		default:
			log.Println(fmt.Sprintf("Command<v>: %s, %s", s[0], input))
			log.Println(fmt.Sprintf("Command not found: %s", input))
		}
	}
	return
}

// Show ...
func (g *GameType) Show() {
	texts = append(texts, g.showGameInfo())
}

func (g *GameType) reset() {
	g.resetInfo()
	// Assemble and shuffle the player decks.
	log.Println("reset OK.")
}

func (g *GameType) resetRank() {
	g.resetInfo()
	g.rank = make(map[string]*rankType)
	g.startPhase()

	log.Println("resetRank OK.")
}

// PlayerCard ...
func (g *GameType) PlayerCard(input string, currentID string) {
	var strs []string
	s := strings.Fields(input)
	if len(s) == 1 {
		strs = append(strs, g.Player(currentID).GetDisplayName())
		strs = append(strs, "【使用卡牌】")
		for k, co := range g.Players.Player(currentID).CardPile.Cards {
			if co.OnPlayFunc != nil {
				strs = append(strs, fmt.Sprintf("<%s>%s", k, co.ViewCardInfo()))
			}
		}
		strs = append(strs, "[[指令]]\n/u <卡片編號>")
		texts = append(texts, strings.Join(strs, "\n"))
		return
	}
	if co, exist := g.Players.Player(currentID).CardPile.Cards[s[1]]; exist {
		if co.OnPlayFunc != nil {
			if r, str := co.OnPlayFunc(); r {
				if len(str) > 0 {
					if len(texts) > 0 {
						texts[0] = str + "\n" + texts[0]
					} else {
						texts = append(texts, str)
					}
				}
			}
		}
	}
}

// GamePhase ...
func (g *GameType) GamePhase(input string) {
	g.runPhase(input)
}

func (g *GameType) mysterPhase() {

}

// OnPlay ...
func (g *GameType) OnPlay() {
	if exists, _ := helper.InArray(g.Info.CurrentUserID, g.Info.Queue); !exists {
		g.Info.Queue = append(g.Info.Queue, g.Info.CurrentUserID)
	}

	g.Info.PlayQueue = append(g.Info.PlayQueue, g.Info.CurrentUserID)

	if _, exist := g.rank[g.Info.CurrentUserID]; !exist {
		g.rank[g.Info.CurrentUserID] = &rankType{g.Info.CurrentUserID, 0, 0, 0}
	}
}

// GetPlayer ...
func (g *GameType) GetPlayer(id string) scheduler.Player {
	return g.Player(id)
}

// GetQueue ...
func (g *GameType) GetQueue() []string {
	return g.Info.Queue
}

// GetPlayQueue ...
func (g *GameType) GetPlayQueue() []string {
	return g.Info.PlayQueue
}

// AddPlayQueue ...
func (g *GameType) AddPlayQueue(s string) {
	g.Info.PlayQueue = append(g.Info.PlayQueue, s)
}

// GetHit ...
func (g *GameType) GetHit() (int, int, int) {
	return g.Info.Hit, g.Info.Min, g.Info.Max
}

// SetHit ...
func (g *GameType) SetHit(hit int, min int, max int) {
	g.Info.Hit = hit
	g.Info.Min = min
	g.Info.Max = max
}

// SetInfoRange ...
func (g *GameType) SetInfoRange(n int, m int) {
	g.Info.Min = n
	g.Info.Max = m
}

// GetInfoCurrent ...
func (g *GameType) GetInfoCurrent() int {
	return g.Info.Current
}

// SetInfoCurrent ...
func (g *GameType) SetInfoCurrent(n int) {
	g.Info.Current = n
}

// GetInfoBoomCnt ...
func (g *GameType) GetInfoBoomCnt() int {
	return g.Info.BoomCnt
}

// MakeInfoBoomCnt ...
func (g *GameType) MakeInfoBoomCnt(n int) {
	g.Info.BoomCnt = helper.Max(0, g.Info.BoomCnt+n)
}

// GetRankBoomCnt ...
func (g *GameType) GetRankBoomCnt(id string) int {
	return g.rank[id].Boom
}

// SetRankBoomCnt ...
func (g *GameType) SetRankBoomCnt(id string, n int) {
	g.rank[id].Boom = n
}

// MakeRankBoomCnt ...
func (g *GameType) MakeRankBoomCnt(id string, n int) {
	g.rank[id].Boom = helper.Max(0, g.rank[id].Boom+n)

}

func (g *GameType) doCoolDown(id string) {
	for _, co := range g.Player(id).Cards {
		co.DoCoolDown()
		co.DoFreeze()
	}
	for _, co := range g.mythosCards {
		co.DoCoolDown()
		co.DoFreeze()
	}
}

// CheckExistData ...
func CheckExistData(SourceID string) {
	if _, exist := Boom[SourceID]; !exist {
		//loadData(SourceID)
		Boom[SourceID] = &GameType{}
		loadData(SourceID)
		log.Println(Boom[SourceID])

		Boom[SourceID].mythosCards = cards.CreateMythosCard()
	}
}

// GetSourceID ...
func (g *GameType) GetSourceID() string {
	return g.sourceid
}

func loadData(SourceID string) {
	row := mydb.Db.QueryRow("SELECT sourceid, season, scene_info, rank FROM boom_game where sourceid = $1 limit 1", SourceID)
	var sourceid string
	var season int
	var sceneInfo json.RawMessage
	var rank json.RawMessage
	switch err := row.Scan(&sourceid, &season, &sceneInfo, &rank); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned")
		Boom[SourceID] = &GameType{}
		Boom[SourceID].sourceid = SourceID
		Boom[SourceID].season = 1
		Boom[SourceID].rank = make(map[string]*rankType)
		Boom[SourceID].Players.LoadPlayersData(Boom[SourceID], SourceID)
		Boom[SourceID].reset()
		Boom[SourceID].addData()
	case nil:
		Boom[SourceID] = &GameType{}
		Boom[SourceID].sourceid = sourceid
		Boom[SourceID].season = season
		Boom[SourceID].rank = make(map[string]*rankType)
		json.Unmarshal(sceneInfo, &Boom[SourceID].Info)
		json.Unmarshal(rank, &Boom[SourceID].rank)
		Boom[SourceID].Players.LoadPlayersData(Boom[SourceID], SourceID)
		log.Println("Boom data load.")
		//Boom[SourceID].data.updateData()
	default:
		checkError(err)
	}
}

func (g *GameType) addData() {
	rank, err := json.Marshal(g.rank)
	checkError(err)

	sceneInfo, err := json.Marshal(g.Info)
	checkError(err)

	stmt, err := mydb.Db.Prepare("insert into boom_game (sourceid, season, scene_info, rank) values ($1, $2, $3, $4)")
	checkError(err)

	_, err = stmt.Exec(g.sourceid, g.season, sceneInfo, rank)
	checkError(err)

	stmt.Close()
	log.Println("Boom Data Create...")
}

func (g *GameType) updateData() {
	rank, err := json.Marshal(g.rank)
	checkError(err)

	sceneInfo, err := json.Marshal(g.Info)
	checkError(err)

	stmt, err := mydb.Db.Prepare("update boom_game set season = $2,  scene_info = $3, rank = $4 where sourceid = $1")
	checkError(err)

	_, err = stmt.Exec(g.sourceid, g.season, sceneInfo, rank)
	checkError(err)

	stmt.Close()
	log.Println("Boom Data Update...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
