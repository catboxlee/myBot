package users

import (
	"log"
	"myBot/mydb"
)

// UserOption ...
type UserOption struct {
	UserID      string `db:"userid"`
	DisplayName string `db:"displayname"`
	Money       int    `db:"money"`
	GemStone    int    `db:"gemstone"`
}

// GetDisplayName ...
func (uo *UserOption) GetDisplayName() string {
	return uo.DisplayName
}

func (uo *UserOption) setDisplayName(name string) {
	uo.DisplayName = name
}

// GetMoney ...
func (uo *UserOption) GetMoney() int {
	return uo.Money
}

// GetGemStone ...
func (uo *UserOption) GetGemStone() int {
	return uo.GemStone
}

// MakeGemStone ...
func (uo *UserOption) MakeGemStone(n int) {
	uo.GemStone += n
}

// SetGemStone ...
func (uo *UserOption) SetGemStone(n int) {
	uo.GemStone = n
}

func (uo *UserOption) addData() {
	log.Println("User data insert:", uo.UserID)
	stmt, err := mydb.Db.Prepare("insert into users(userid, displayname, money, gemstone) values($1, $2, $3, $4)")
	checkError(err)

	_, err = stmt.Exec(uo.UserID, uo.DisplayName, uo.Money, uo.GemStone)
	checkError(err)

	stmt.Close()
}

func (uo *UserOption) saveData() {
	log.Println("User data update:", uo.UserID)
	// SQL update
	stmt, err := mydb.Db.Prepare("update users set money = $2, gemstone = $3 where userid = $1")
	checkError(err)

	_, err = stmt.Exec(uo.UserID, uo.Money, uo.GemStone)
	checkError(err)

	stmt.Close()
}
