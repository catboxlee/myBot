package racegame

import (
	"fmt"
	"log"
	"myBot/games/racegame/cards"
	"myBot/games/racegame/players"
	"myBot/games/racegame/scheduler"
	"myBot/helper"
	"sort"
	"strconv"
	"strings"
)

// GameType ..
type GameType struct {
	players.Players
	sourceid     string `db:"sourceid"`
	defaultMeter int
	Info         *InfoType      `json:"info"`
	mythosCards  cards.CardPile `json:"-"`
}

// InfoType ...
type InfoType struct {
	Phase         bool
	Turn          int
	Meter         int
	GameOver      bool
	currentUserID string
	Queue         []string
	PlayQueue     []string
}

const (
	seasonBoomCount    = 100
	singleGameBonusGem = 25
	seasonGameBonusGem = 390
)

// Race ...
var Race = make(map[string]*GameType)
var texts []string

// Start ...
func (g *GameType) Start() {
	log.Println("<Race.Start()>")
	g.reset()
}

func (g *GameType) startPhase() {
	var strs []string
	strs = append(strs, g.showGameInfo())
	texts = append(texts, strings.Join(strs, "\n"))
}

// Command ...
func (g *GameType) Command(input string, currentID string) []string {
	texts = nil

	if len(input) > 0 {
		g.Players.CheckPlayerExist(currentID)

		// 去除前後空格
		input = strings.Trim(input, " ")

		// 檢查字首
		switch input[0:1] {
		case "/": // Game Command
			g.checkCommand(strings.TrimLeft(input, "/"), currentID)
		case "+": // Player Join
			g.Info.currentUserID = currentID
			if s := g.playerJoin(currentID, strings.TrimLeft(input, "+")); len(s) > 0 {
				texts = append(texts, s)
				if len(g.Info.Queue) == 3 {
					g.Info.Phase = true
				}
				g.Show()
			}
		case "-": // Gamming Start
			g.Info.Phase = true
			helper.Shuffle(g.Info.Queue)
			g.Show()
		default: // Player Phase
			g.Info.currentUserID = currentID
			g.GamePhase(input)
		}
	}
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
			log.Println(fmt.Sprintf("Command<rank>: %s, %s", s[0], input))
			g.Show()
		case "set":
			log.Println(fmt.Sprintf("Command<set>: %s, %s", s[0], input))
			if len(s) > 1 {
				if x, err := strconv.Atoi(s[1]); err == nil {
					g.SetMeter(x)
					g.reset()
					g.Show()
				}
			}
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
	g.ClearPlayers()
	g.mythosCards.CreateMythosCard()
	g.resetInfo()
	// Assemble and shuffle the player decks.
	log.Println("reset OK.")
}

// GamePhase ...
func (g *GameType) GamePhase(input string) {
	g.runPhase(input)
}

func (g *GameType) mysterPhase() {

}

func (g *GameType) playerJoin(currentID string, input string) string {
	var strs []string
	input = strings.Trim(input, " ")
	if exists, _ := helper.InArray(currentID, g.Info.Queue); !exists {
		g.Info.Queue = append(g.Info.Queue, currentID)
		strs = append(strs, fmt.Sprintf("%s join.", g.Players.Data[currentID].GetDisplayName()))
	}

	if _, exist := g.Players.Data[currentID]; exist {
		if len(input) > 0 {
			if _, exist := g.mythosCards.Cards[input]; exist {
				g.Players.Data[currentID].Property.AddBuff(input)
			}
		}
	}
	return strings.Join(strs, "\n")
}

// OnPlay ...
func (g *GameType) OnPlay() {
	if exists, _ := helper.InArray(g.Info.currentUserID, g.Info.Queue); !exists {
		g.Info.Queue = append(g.Info.Queue, g.Info.currentUserID)
	}

	g.Info.PlayQueue = append(g.Info.PlayQueue, g.Info.currentUserID)

}

// GetPlayer ...
func (g *GameType) GetPlayer(id string) scheduler.Player {
	return g.Player(id)
}

// GetQueue ...
func (g *GameType) GetQueue() []string {
	return g.Info.Queue
}

// GetQueue ...
func (g *GameType) GetNextQueue() string {
	return g.Info.Queue[(g.Info.Turn)%len(g.Info.Queue)]
}

// GetPlayQueue ...
func (g *GameType) GetPlayQueue() []string {
	return g.Info.PlayQueue
}

// AddPlayQueue ...
func (g *GameType) AddPlayQueue(s string) {
	g.Info.PlayQueue = append(g.Info.PlayQueue, s)
}

// PopCard ...
func (g *GameType) PopCard(po scheduler.Player) {
	if len(g.mythosCards.Deck) > 1 {
		coid := g.mythosCards.PopCard()
		po.TakeCard(coid)
	}
}

// GetMeter ...
func (g *GameType) GetMeter() int {
	return g.Info.Meter
}

func (g *GameType) GetRankingArray() []string {
	values := make([]*players.PlayerOption, 0, len(g.Info.Queue))
	for _, v := range g.Info.Queue {
		values = append(values, g.Player(v))
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].TotalMove > values[j].TotalMove
	})
	var ranking []string
	for _, v := range values {
		ranking = append(ranking, v.GetUserID())
	}
	return ranking
}

func (g *GameType) getRaceSort() []*players.PlayerOption {
	values := make([]*players.PlayerOption, 0, len(g.Info.Queue))
	for _, v := range g.Info.Queue {
		values = append(values, g.Player(v))
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].TotalMove > values[j].TotalMove
	})
	return values
}

// GetRanking 目前名次
func (g *GameType) GetRanking(userID string) int {
	var tmp []string
	for _, po := range g.getRaceSort() {
		tmp = append(tmp, po.GetUserID())
	}
	if exists, i := helper.InArray(userID, tmp); exists {
		return i
	}
	return -1
}

// GetRanking 目前名次
func (g *GameType) ViewCardsInfo(cid string) string {
	if _, exist := g.mythosCards.Cards[cid]; exist {
		return g.mythosCards.Cards[cid].GetDisplayName()
	}
	return ""
}

// defaultMeter
func (g *GameType) SetMeter(n int) {
	g.defaultMeter = n
}

// CheckExistData ...
func CheckExistData(SourceID string) {
	if _, exist := Race[SourceID]; !exist {
		//loadData(SourceID)
		Race[SourceID] = &GameType{}
		//Boom[SourceID].Start()
		Race[SourceID].sourceid = SourceID
		Race[SourceID].defaultMeter = 30

		Race[SourceID].Players.LoadPlayersData(Race[SourceID], SourceID)

		Race[SourceID].mythosCards.Parent = Race[SourceID]
		Race[SourceID].reset()
		log.Println(Race[SourceID])

	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
