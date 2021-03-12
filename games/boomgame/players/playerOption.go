package players

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"myBot/emoji"
	"myBot/games/boomgame/cards"
	"myBot/games/boomgame/scheduler"
	"myBot/helper"
	"myBot/mydb"
	"myBot/users"
	"strings"
)

// PlayerOption ...
type PlayerOption struct {
	Parent         scheduler.Players
	UserID         string   `db:"userid"`
	SourceID       string   `db:"sourceid"`
	Titles         []string `json:"titles"`
	cards.CardPile `json:"cardpile"`
	CoolDownData   map[string]int `json:"cooldowndata"`
	AchievementType
	Property `json:"property"`
}

// Property ...
type Property struct {
	OneShot       int `json:"oneshot"`
	WinningStreak int `json:"winningstreak"`
	LosingStreak  int `json:"losingstreak"`
}

// AchievementType 成就
type AchievementType struct {
}

// GetTopParent ...
func (po *PlayerOption) GetTopParent() scheduler.Game {
	return po.GetParent().GetTopParent()
}

// GetParent ...
func (po *PlayerOption) GetParent() scheduler.Players {
	return po.Parent
}

// GetUserID ...
func (po *PlayerOption) GetUserID() string {
	return po.UserID
}

// GetDisplayName ...
func (po *PlayerOption) GetDisplayName() string {
	return users.UsersList.User(po.UserID).GetDisplayName()
}

// GetGemStone ...
func (po *PlayerOption) GetGemStone() int {
	return users.UsersList.User(po.UserID).GetGemStone()
}

// MakeGemStone ...
func (po *PlayerOption) MakeGemStone(n int) {
	users.UsersList.User(po.UserID).MakeGemStone(n)
}

// GetMoney ...
func (po *PlayerOption) GetMoney() int {
	return users.UsersList.User(po.UserID).GetMoney()
}

// GetCardPile ...
func (po *PlayerOption) GetCardPile() scheduler.Cards {
	return &po.CardPile
}

// GetRandCards ...
func (po *PlayerOption) GetRandCards(n int) []scheduler.Card {
	var IDs []string
	var cos []scheduler.Card
	if len(po.Cards) > 0 {
		for cardID := range po.Cards {
			IDs = append(IDs, cardID)
		}
		n = helper.Min(n, len(po.Cards))
		tmp := rand.Perm(len(IDs))
		for i := 0; i < n; i++ {
			cos = append(cos, po.Cards[IDs[tmp[i]]])
		}
	}
	return cos
}

// ViewInfo ...
func (po *PlayerOption) ViewInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s", po.GetDisplayName()))
	strs = append(strs, fmt.Sprintf("【資產】\n%s%d\n%s%d", emoji.Emoji(":money_bag:"), po.GetMoney(), emoji.Emoji(":gem_stone:"), po.GetGemStone()))
	strs = append(strs, fmt.Sprintf("【擁有卡片】"))
	strs = append(strs, po.CardPile.ViewCardsInfo())
	strs = append(strs, "[[指令]]")
	strs = append(strs, "使用卡牌: /u <卡片編號>")
	strs = append(strs, "轉蛋: /gacha")
	return strings.Join(strs, "\n")
}

func (po *PlayerOption) addData() {
	titles, err := json.Marshal(po.Titles)
	checkError(err)
	cardpile, err := json.Marshal(po.CardPile)
	checkError(err)
	property, err := json.Marshal(po.Property)
	checkError(err)

	log.Println("PlayerOption data insert:", po.UserID, po.SourceID, string(titles), string(cardpile), string(property))
	stmt, err := mydb.Db.Prepare("insert into boomplayer(userid, sourceid, titles, cardpile, itempile, property) values($1, $2, $3, $4, $5)")
	checkError(err)

	_, err = stmt.Exec(po.UserID, po.SourceID, titles, cardpile, property)
	checkError(err)

	stmt.Close()
}

// SaveData ...
func (po *PlayerOption) SaveData() {
	titles, err := json.Marshal(po.Titles)
	checkError(err)
	cardpile, err := json.Marshal(po.CardPile)
	checkError(err)
	property, err := json.Marshal(po.Property)
	checkError(err)

	log.Println("PlayerOption.updateData()", po.UserID, po.SourceID, string(titles), string(cardpile), string(property))

	// SQL update
	stmt, err := mydb.Db.Prepare("update boomplayer set titles = $3, cardpile = $4,  property = $5 where userid = $1 and sourceid = $2")
	checkError(err)

	_, err = stmt.Exec(po.UserID, po.SourceID, titles, cardpile, property)
	checkError(err)

	stmt.Close()
}
