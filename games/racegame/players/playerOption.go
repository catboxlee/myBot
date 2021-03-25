package players

import (
	"encoding/json"
	"fmt"
	"log"
	"myBot/emoji"
	"myBot/games/racegame/data"
	"myBot/games/racegame/scheduler"
	"myBot/games/racegame/typeset"
	"myBot/helper"
	"myBot/users"
	"strings"
)

// PlayerOption ...
type PlayerOption struct {
	Parent            scheduler.Players `json:"-,omitempty"`
	UserID            string            `db:"userid" json:"user_id,omitempty"`
	SourceID          string            `db:"sourceid" json:"source_id,omitempty"`
	CardPile          []string          `json:"cardpile"`
	*typeset.Property `json:"property,omitempty"`
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
	return fmt.Sprintf("%s%s", emoji.Emoji(":horse:"), users.UsersList.User(po.UserID).GetDisplayName())
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

// GetProperty ...
func (po *PlayerOption) GetProperty() *typeset.Property {
	return po.Property
}

// GetCardPile ...
func (po *PlayerOption) GetCardPile() []string {
	return po.CardPile
}

// RemoveCardPile ...
func (po *PlayerOption) RemoveAllDeBuff() (isRemoe bool, str string) {
	var strs []string
	for _, val := range po.DeBuff {
		strs = append(strs, fmt.Sprintf("移除%s", data.CardData[val].CardName))
	}
	po.DeBuff = nil
	str = strings.Join(strs, "\n")
	return
}

// RemoveCardPile ...
func (po *PlayerOption) RemoveDeBuff(ids ...string) (isRemoe bool, str string) {
	var tmp []string
	var strs []string
	for _, val := range po.DeBuff {
		exist, _ := helper.InArray(val, ids)
		if !exist {
			tmp = append(tmp, val)
		} else {
			strs = append(strs, fmt.Sprintf("移除%s", data.CardData[val].CardName))
			log.Println(strs)
			isRemoe = true
		}
	}
	po.DeBuff = tmp
	str = strings.Join(strs, "\n")
	return
}

// TakeCard ...
func (po *PlayerOption) TakeCard(id string) {
	po.CardPile = append(po.CardPile, id)
}

// ViewInfo ...
func (po *PlayerOption) ViewInfo() string {
	var strs []string
	strs = append(strs, fmt.Sprintf("%s", po.GetDisplayName()))
	strs = append(strs, fmt.Sprintf("【資產】\n%s%d\n%s%d", emoji.Emoji(":money_bag:"), po.GetMoney(), emoji.Emoji(":gem_stone:"), po.GetGemStone()))
	strs = append(strs, fmt.Sprintf("【擁有卡片】"))
	strs = append(strs, "[[指令]]")
	strs = append(strs, "使用卡牌: /u <卡片編號>")
	strs = append(strs, "轉蛋: /gacha")
	return strings.Join(strs, "\n")
}

// ViewCardsInfo ...
func (po *PlayerOption) ViewCardsInfo() string {
	var strs []string
	g := po.GetTopParent()
	for i, coid := range po.CardPile {
		strs = append(strs, fmt.Sprintf("%d.%s", i+1, g.ViewCardsInfo(coid)))
	}
	return strings.Join(strs, "\n")
}

func (po *PlayerOption) addData() {
	cardpile, err := json.Marshal(po.CardPile)
	checkError(err)
	property, err := json.Marshal(po.Property)
	checkError(err)

	log.Println("PlayerOption data insert:", po.UserID, po.SourceID, string(cardpile), string(property))
	/*
		stmt, err := mydb.Db.Prepare("insert into boomplayer(userid, sourceid, titles, cardpile, itempile, property) values($1, $2, $3, $4, $5, $6)")
		checkError(err)

		_, err = stmt.Exec(po.UserID, po.SourceID, titles, cardpile, itempile, property)
		checkError(err)

		stmt.Close()
	*/
}

// SaveData ...
func (po *PlayerOption) SaveData() {
	cardpile, err := json.Marshal(po.CardPile)
	checkError(err)
	property, err := json.Marshal(po.Property)
	checkError(err)

	log.Println("PlayerOption.updateData()", po.UserID, po.SourceID, string(cardpile), string(property))

	// SQL update
	/*
		stmt, err := mydb.Db.Prepare("update boomplayer set titles = $3, cardpile = $4, itempile = $5, property = $6 where userid = $1 and sourceid = $2")
		checkError(err)

		_, err = stmt.Exec(po.UserID, po.SourceID, titles, cardpile, itempile, property)
		checkError(err)

		stmt.Close()
	*/
}
