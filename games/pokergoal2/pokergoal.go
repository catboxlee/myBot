package pokergoal2

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
	gate1		cardType
	gate2		cardType
	deck        []cardType
	discardPile []cardType
	players     map[string]*playerType
}

type playerType struct {
	UserID      string
	DisplayName string
	bets        int
	ball 		cardType
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
	log.Println("pokergoal init")
	Pokergoal.pot = world.World.Pot2
	//Shuffle(p.deck)
}

func (p *gameType) Run(input string) []string {
	texts = nil
	var text string

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		p.checkCommand(strings.ToLower(strings.TrimLeft(input, "/")))
		return texts
	} else if strings.HasPrefix(input, "-") {
		// 開牌
		p.dealGate()
		p.showPot()
		return texts
		
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
		
		if currentPlayer.bets <= 0 {
			// 下注
			re := regexp.MustCompile(`^\+(\d+)`)
			matches := re.FindStringSubmatch(input)
			bets := 1
			if len(matches) > 1 {
				if bet, err := strconv.Atoi(matches[1]); err == nil {
					bets = helper.Max(bets, bet)
				}
			}
			currentPlayer.bets = bets
			p.pot += bets
			users.UsersList.Data[users.LineUser.UserProfile.UserID].Money -= bets
			//users.LineUser.SaveUserData()
			// 拿牌
			currentPlayer.ball = p.deal()
			text = fmt.Sprintf("%s 下注：%s%d (%s%d)", users.LineUser.UserProfile.DisplayName, emoji.Emoji(":money_bag:"), bets, emoji.Emoji(":money_bag:"),  users.UsersList.Data[users.LineUser.UserProfile.UserID].Money)
			text += fmt.Sprintf("\n目前獎池：%s%d(%+d)", emoji.Emoji(":money_bag:"), p.pot, bets)
			text += fmt.Sprintf("\n%s", convCard(currentPlayer.ball))
			texts = append(texts, text)
			
		} else {
			text = fmt.Sprintf("%s %s", currentPlayer.DisplayName, convCard(currentPlayer.ball))
			texts = append(texts, text)
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

func (p *gameType) dealGate() {
	// 發門柱
	p.gate1 = p.deal()
	p.gate2 = p.deal()
	str := fmt.Sprintf("%s %s %s", convCard(p.gate1), emoji.Emoji(":goal_net:"), convCard(p.gate2))
	if len(p.players) > 0 {
		for _,v := range p.players {
			str += p.hit(v)
		}
	}
}

// 要牌
func (p *gameType) hit(currentPlayer *playerType) string {
	// 結算
	str := fmt.Sprintf("\n%s %s", currentPlayer.DisplayName, convCard(currentPlayer.ball))
	bets := 0
	if currentPlayer.ball.number == p.gate1.number || currentPlayer.ball.number == p.gate2.number {
		// 撞柱
		bets = -(currentPlayer.bets)
		p.pot -= bets
		str += " 撞柱"
	} else if currentPlayer.ball.number < helper.Min(p.gate1.number, p.gate2.number) || currentPlayer.ball.number > helper.Max(p.gate1.number, p.gate2.number) {
		// 未入門
		//bets = -(currentPlayer.bets)
		//p.pot -= bets
		str += " 不中"
	} else {
		bets = currentPlayer.bets * 2
		p.pot -= bets
		str += " Goal!!!"
	}
	users.UsersList.Data[currentPlayer.UserID].Money += bets
	str += fmt.Sprintf("%s%d(%+d)", emoji.Emoji(":money_bag:"), users.UsersList.Data[currentPlayer.UserID].Money, bets)
	// 結算
	users.LineUser.SaveUserData()

	//texts = append(texts, str)
	p.endGame(currentPlayer, bets)
	return str
}

func (p *gameType) endGame(currentPlayer *playerType, bets int) (){
	world.World.Pot2 = p.pot
	world.World.SaveWorldData()
	
	// 清理桌面
	p.discardPile = append(p.discardPile, currentPlayer.ball)
	p.discardPile = append(p.discardPile, p.gate1)
	p.discardPile = append(p.discardPile, p.gate2)
	//p.gate = nil
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
