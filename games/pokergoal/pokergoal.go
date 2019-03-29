package pokergoal

import (
	"fmt"
	"log"
	"math/rand"
	"myBot/emoji"
	"myBot/helper"
	"myBot/world"
	"myBot/users"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type gameType struct {
	pot         int
	antes       int
	deck        []cardType
	discardPile []cardType
	players     map[string]*playerType
}

type playerType struct {
	UserID      string
	DisplayName string
	bets        int
	cards       []cardType
}

type cardType struct {
	suit   int
	number int
}

// Pokergoal ...
var Pokergoal gameType
var texts []string
var cardFaces = struct {
	numbers []string
	suits   []string
}{
	[]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"},
	[]string{emoji.Emoji(":spade_suit:"), emoji.Emoji(":heart_suit:"), emoji.Emoji(":diamond_suit:"), emoji.Emoji(":club_suit:")},
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	Pokergoal.createDeck()

	Shuffle(Pokergoal.deck)
	Pokergoal.players = make(map[string]*playerType)
	Pokergoal.antes = 1
	//Pokergoal.pot = world.World.Bank
	//Shuffle(p.deck)
}

func (p *gameType) Run(input string) []string {
	texts = nil
	var text string
	if p.pot == 0 {
		p.pot = 10
	}

	//p.showAllCard()

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		p.checkCommand(strings.ToLower(strings.TrimLeft(input, "/")))
		return texts
	} else if strings.HasPrefix(input, "-") {
		// 棄牌
		if _, exist := p.players[users.LineUser.UserProfile.UserID]; exist {
			if len(p.players[users.LineUser.UserProfile.UserID].cards) > 0 {
				p.discardPile = append(p.discardPile, p.players[users.LineUser.UserProfile.UserID].cards...)
				p.players[users.LineUser.UserProfile.UserID].cards = nil
				p.pot += p.players[users.LineUser.UserProfile.UserID].bets
				p.players[users.LineUser.UserProfile.UserID].bets = 0
				texts = append(texts, "棄牌")
				p.showPot()
				return texts
			}
		}
	} else if strings.HasPrefix(input, "+") {
		// New Player
		if _, exist := p.players[users.LineUser.UserProfile.UserID]; !exist {
			p.players[users.LineUser.UserProfile.UserID] = &playerType{}
			p.players[users.LineUser.UserProfile.UserID].UserID = users.LineUser.UserProfile.UserID
			p.players[users.LineUser.UserProfile.UserID].DisplayName = users.LineUser.UserProfile.DisplayName
			p.players[users.LineUser.UserProfile.UserID].bets = 0
			log.Println("新玩家")
		}
		currentPlayer := p.players[users.LineUser.UserProfile.UserID]

		if len(currentPlayer.cards) < 2 {
			p.pot++
			users.UsersList.Data[users.LineUser.UserProfile.UserID].Money--
			users.LineUser.SaveUserData()
			p.showPot()
			text = fmt.Sprintf("%s 下注：%s %d\n", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":money_bag:"), p.antes)
			text += p.dealGate(currentPlayer)
			text += fmt.Sprintf("\n剩餘資金：%d",  users.UsersList.Data[users.LineUser.UserProfile.UserID].Money)
			if len(currentPlayer.cards) < 2 {
				text += fmt.Sprintf("\n可加注：+%d ~ +%d (預設+0)", 0, helper.Min(users.UsersList.Data[users.LineUser.UserProfile.UserID].Money, p.pot))
			}
			texts = append(texts, text)
			
		} else {
			// 下注
			re := regexp.MustCompile(`^\+(\d+)`)
			matches := re.FindStringSubmatch(input)
			bet := 0
			if len(matches) > 1 {
				// 喊注
				if bet, _ = strconv.Atoi(matches[1]); bet < 0 {
					texts = append(texts, fmt.Sprintf("%s 加注金額錯誤: %d", users.LineUser.UserProfile.DisplayName, bet))
					return texts
				}
			}
			bet = helper.Min(helper.Max(0, bet), helper.Min(users.UsersList.Data[users.LineUser.UserProfile.UserID].Money, p.pot))
			currentPlayer.bets = bet
			texts = append(texts, fmt.Sprintf("%s 加注: %s %d", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":money_bag:"), helper.Max(0, bet)))
			p.hit(currentPlayer)
			p.showPot()
		}
	}
	return texts
}

// 指令
func (p *gameType) checkCommand(input string) {
	switch input {
	case "reset":
		p.createDeck()
		Shuffle(p.deck)
	case "shuffle":
		Shuffle(p.deck)
	case "pot":
		p.showPot()
	}
}

func (p *gameType) showAllCard() {
	for _, v := range p.deck {
		fmt.Print(fmt.Sprintf("%s ", convCard(v)))
	}
	fmt.Println()
}

// 獎池
func (p *gameType) showPot() {
	text := "[射龍門]"
	text += fmt.Sprintf("\n獎池: %s %d", emoji.Emoji(":money_bag:"), p.pot)
	texts = append(texts, text)
}

func (p *gameType) dealGate(currentPlayer *playerType) string {
	// 發門柱
	currentPlayer.cards = append(currentPlayer.cards, p.deal())
	currentPlayer.cards = append(currentPlayer.cards, p.deal())
	if currentPlayer.cards[0].number == currentPlayer.cards[1].number {
		p.discardPile = append(p.discardPile, currentPlayer.cards...)
		text := fmt.Sprintf("%s %s %s\n同數字判定未完工，此局賭金充公", convCard(currentPlayer.cards[0]), emoji.Emoji(":goal_net:"), convCard(currentPlayer.cards[1]))
		currentPlayer.cards = nil
		currentPlayer.bets = 0
		return text
	}
	return fmt.Sprintf("%s %s %s", convCard(currentPlayer.cards[0]), emoji.Emoji(":goal_net:"), convCard(currentPlayer.cards[1]))
}

// 要牌
func (p *gameType) hit(currentPlayer *playerType) {
	// 發球
	currentPlayer.cards = append(currentPlayer.cards, p.deal())
	str := fmt.Sprintf("%s", convCard(currentPlayer.cards[2]))
	// 結算
	bets := 0
	if currentPlayer.cards[2].number == currentPlayer.cards[0].number || currentPlayer.cards[2].number == currentPlayer.cards[1].number {
		// 撞柱
		bets = -((currentPlayer.bets + p.antes) * 2)
		p.pot -= bets
		texts = append(texts, fmt.Sprintf("%s 撞柱", str))
	} else if currentPlayer.cards[2].number < helper.Min(currentPlayer.cards[0].number, currentPlayer.cards[1].number) || currentPlayer.cards[2].number > helper.Max(currentPlayer.cards[0].number, currentPlayer.cards[1].number) {
		// 未入門
		bets = -(currentPlayer.bets + p.antes)
		p.pot -= bets
		texts = append(texts, fmt.Sprintf("%s 不中", str))
	} else {
		bets = (currentPlayer.bets + p.antes)
		p.pot -= bets
		texts = append(texts, fmt.Sprintf("%s Goal!!!", str))
	}
	// 結算
	p.endGame(currentPlayer, bets)
}

func (p *gameType) endGame(currentPlayer *playerType, bets int) (){
	users.UsersList.Data[users.LineUser.UserProfile.UserID].Money += bets
	users.LineUser.SaveUserData()
	
	if p.pot <= 0 {
		p.pot = 10
		texts = append(texts, "補充獎池：10")
	}
	//world.World.Bank = p.pot
	//world.World.SaveWorldData()
	
	// 清理桌面
	p.discardPile = append(p.discardPile, currentPlayer.cards...)
	currentPlayer.cards = nil
	currentPlayer.bets = 0
	//log.Println(p.discardPile)
}

// 發牌
func (p *gameType) deal() cardType {
	if len(p.deck) <= 0 {
		p.reDeck()
	}
	dealCard := p.deck[0]
	p.deck = p.deck[1:]
	return dealCard
}

func (p *gameType) reDeck() {
	p.createDeck()
	Shuffle(p.deck)
	//Shuffle(p.discardPile)                      // 棄牌堆重洗
	//p.deck = append(p.deck, p.discardPile...) // 加入牌堆
	p.discardPile = nil // 清除棄牌堆
}

func (p *gameType) createDeck() {
	p.deck = nil
	//p.deck = make([]cardType, 52)
	for j := 0; j < 4; j++ {
		for i := 0; i < 13; i++ {
			//p.deck[(j*13)+i] = cardType{j, i}
			p.deck = append(p.deck, cardType{j, i})
		}
	}
}

// 牌面
func convCard(n cardType) string {
	//numbers := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	//suits := []string{emoji.Emoji(":spade_suit:"), emoji.Emoji(":heart_suit:"), emoji.Emoji(":diamond_suit:"), emoji.Emoji(":club_suit:")}
	return fmt.Sprintf("%s%s", cardFaces.suits[n.suit], cardFaces.numbers[n.number])
}

// Shuffle 洗牌
func Shuffle(vals []cardType) {
	rand.Seed(time.Now().UnixNano())
	for len(vals) > 0 { //根据牌面数组长度遍历
		n := len(vals)                                          //数组长度
		randIndex := rand.Intn(n)                               //得到随机index
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1] //最后一张牌和第randIndex张牌互换
		vals = vals[:n-1]
	}
}
