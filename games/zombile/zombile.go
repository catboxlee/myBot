package zombile

import (
	"fmt"
	"log"
	"math/rand"
	"myBot/games/zombile/cards"
	"myBot/games/zombile/players"
	"myBot/games/zombile/power"
	"myBot/helper"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GameType ..
type GameType struct {
	deck            cards.Cards
	discardPile     cards.Cards
	playersSequence []string
	players.Players
	phase           int
	turn            int
	currentPlayerID string
	simpleList      bool
}

const (
	readyPhase int = iota
	gamePhase
)

// Zombile ...
var Zombile = make(map[string]*GameType)
var texts []string

func init() {
	log.Println("<Zombile init>")
}

// Start ...
func (g *GameType) Start() {
	log.Println("<zombie.Start()>")
	g.simpleList = true
	g.reset()
}

// Command ...
func (g *GameType) Command(input string) []string {
	texts = nil
	if len(input) > 0 {
		switch input[0:1] {
		case "/":
			g.checkCommand(strings.TrimLeft(input, "/"))
		// Player join
		case "+":
			g.PlayerJoin()
			log.Println("p.players:", g.Players)
		// Opening
		case "-":
			if g.phase == gamePhase {
				texts = append(texts, "\n"+g.turnChange()+"\n"+g.ViewPlayersInfo())
			} else {
				g.phase = gamePhase
				// deal open cards
				// print all players info with cards
				texts = append(texts, "deal open cards\n"+g.dealCards(2))
				texts = append(texts, "\n"+g.turnChange()+"\n\n"+g.ViewPlayersInfo())
			}
		default:
			if len(g.Players.Player) > 0 && g.phase == gamePhase {
				if s := g.runPhase(input); len(s) > 0 {
					texts = append(texts, s+"\n\n"+g.getWhosTurn()+"\n"+g.ViewPlayersInfo())
				}
			}
		}
	}
	return texts
}

func (g *GameType) getWhosTurn() (s string) {
	g.currentPlayerID = g.playersSequence[g.turn%len(g.playersSequence)]
	if val := g.Players.GetPlayer(g.currentPlayerID); val != nil {
		if val.GetHealth() > 0 {
			s = fmt.Sprintf("Turn(%d): %s行動.", int(g.turn/len(g.playersSequence)), val.GetDisplayName())
		}
	}
	return
}

func (g *GameType) turnChange() string {
	var strs []string
	g.turn++
	g.currentPlayerID = g.playersSequence[g.turn%len(g.playersSequence)]
	if g.turn%len(g.playersSequence) == 0 {
		strs = append(strs, g.mysterPhase())
	}
	if val := g.Players.GetPlayer(g.currentPlayerID); val != nil {
		if val.GetHealth() > 0 {
			strs = append(strs, fmt.Sprintf("\nTurn(%d): %s行動.", int(g.turn/len(g.playersSequence)), val.GetDisplayName()))
			if g.deck.GetCardsCount() > 0 {
				strs = append(strs, g.DrawCards(g.currentPlayerID, 1))
			}
		} else {
			strs = append(strs, g.turnChange())
		}
	}
	return strings.Join(strs, "\n")
}

func (g *GameType) mysterPhase() string {
	var strs []string
	strs = append(strs, "Myster Phase")
	for _, poid := range g.playersSequence {
		for _, co := range g.Player[poid].GetAllCards() {
			if co.OnMysterFunc != nil {
				strs = append(strs, co.OnMysterFunc(g.Player[poid], nil))
			}
		}
	}
	for _, poid := range g.playersSequence {
		for _, co := range g.Player[poid].GetAllCards() {
			co.ResetActionTimes()
		}
	}
	return strings.Join(strs, "\n")
}

// GetPlayer ...
func (g *GameType) GetPlayer(playerID string) power.PlayerIF {
	return g.Players.GetPlayer(playerID)
}

func (g *GameType) checkCommand(input string) {
	re := regexp.MustCompile(`(\S+)\s*(.*)\s*$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 0 {
		switch matches[1] {
		case "reset":
			log.Println(fmt.Sprintf("command: %s", input))
			g.phase = readyPhase
			g.reset()
		case "listType":
			g.simpleList = !g.simpleList
		case "play":
		case "?":
			if len(matches[2]) > 0 {
				g.viewCardInfo(matches[2])
			} else {
				texts = append(texts, g.ViewPlayersInfo())
			}
		default:
			log.Println(fmt.Sprintf("command not found: %s", input))
		}
	}
}

func (g *GameType) reset() {
	// 初始化
	g.deck.ClearCards()        // 清空牌堆
	g.discardPile.ClearCards() // 清空棄牌堆
	g.ClearPlayers()           // 清空players
	g.turn = -1
	g.currentPlayerID = ""
	power.Power.GameIF = power.NewPower(g)
	log.Println("reset OK.")

	// Assemble and shuffle the player decks.
	// 建立牌庫
	g.deck.CreateCardDeck()
	// 洗牌
	Shuffle(g.deck.GetAllCards())

}

// runPhase ...
func (g *GameType) runPhase(input string) (s string) {
	re := regexp.MustCompile(`^(\d+)\s*(.*)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 2 {
		if x, err := strconv.Atoi(matches[1]); err == nil {
			s = g.activateCard(g.currentPlayerID, x, matches[2])
		}
	}
	return
}

func (g *GameType) activateCard(playerID string, cardIndex int, cmd string) string {
	return g.Players.ActivateCard(playerID, cardIndex, cmd)
}

func (g *GameType) viewCardInfo(input string) {
	if x, err := strconv.Atoi(input); err == nil {
		texts = append(texts, g.Player[g.currentPlayerID].ViewCardFullInfo(x))
	}
}

// dealCards ...
func (g *GameType) dealCards(n int) string {
	var strs []string
	for i := 0; i < n; i++ {
		for _, val := range g.playersSequence {
			if po := g.Players.GetPlayer(val); po != nil {
				if po.GetHealth() > 0 {
					strs = append(strs, g.DrawCards(val, 1))
				}
			}
		}
	}
	return strings.Join(strs, "\n")
}

// DisCards ...
func (g *GameType) DisCards(playerID string, cardIndex int) {
	if po := g.Players.GetPlayer(playerID); po != nil {
		if co := po.GetCard(cardIndex); co != nil {
			defer po.RemoveCards(cardIndex)
			g.discardPile.TakeCard(po, co)
		}
	}
}

// DrawCards (playerID string, 數量 int)
func (g *GameType) DrawCards(playerID string, n int) string {
	var strs []string
	if po := g.Players.GetPlayer(playerID).(*players.PlayerOption); po != nil {
		for i := 0; i < n; i++ {
			if nc := g.deck.PopCards(); nc != nil {
				po.TakeCard(po, nc)
				strs = append(strs, fmt.Sprintf("抽卡: [%s]獲得<%s>", po.GetDisplayName(), nc.GetDisplayName()))
				if nc.OnDisplayFunc != nil {
					strs = append(strs, nc.OnDisplayFunc(po, nil))
				}
			}
		}
	}
	return strings.Join(strs, "\n")
}

// MoveCards ...
func (g *GameType) MoveCards(fromPlayer power.PlayerIF, tCard power.CardIF, toPlayer power.PlayerIF) (s string) {
	if exists, i := helper.InArray(tCard.(*cards.CardOption), fromPlayer.(*players.PlayerOption).Card); exists {
		fromPlayer.RemoveCards(i)
		toPlayer.TakeCard(toPlayer, tCard)
	}
	return
}

// PlayerJoin ...
func (g *GameType) PlayerJoin() {
	log.Println("zombile::PlayerJoin():", users.LineUser.UserProfile.UserID)
	g.playersSequence = append(g.playersSequence, users.LineUser.UserProfile.UserID)
	if r, n := g.Players.PlayerJoin(users.LineUser.UserProfile.UserID, users.LineUser.UserProfile.DisplayName); r {
		texts = append(texts, n+" join.")
	}
	/*
		g.playersSequence = append(g.playersSequence, "p2")
		if r, n := g.Players.PlayerJoin("p2", "Player2"); r {
			texts = append(texts, n+" join.")
		}
	*/
}

// ViewPlayersInfo ...
func (g *GameType) ViewPlayersInfo() string {
	var strs []string
	for _, val := range g.playersSequence {
		if po := g.Players.GetPlayer(val).(*players.PlayerOption); po != nil {
			if g.simpleList {
				strs = append(strs, po.ViewPlayerSimpleInfo())
			} else {
				strs = append(strs, po.ViewPlayerInfo())
			}
		}
	}
	return strings.Join(strs, "\n")
}

// GetPlayersSequence ...
func (g *GameType) GetPlayersSequence() []string {
	return g.playersSequence
}

// CheckExistData ...
func CheckExistData(SourceID string) {
	if _, exist := Zombile[SourceID]; !exist {
		//loadData(SourceID)
		Zombile[SourceID] = &GameType{}
	}
}

// Shuffle 洗牌
func Shuffle(vals []*cards.CardOption) {
	rand.Seed(time.Now().UnixNano())
	for len(vals) > 0 {
		n := len(vals)                                          // 陣列長度
		randIndex := rand.Intn(n)                               // 取隨機index
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1] // 將最後一張牌和第randIndex張牌互換
		vals = vals[:n-1]
	}
}
