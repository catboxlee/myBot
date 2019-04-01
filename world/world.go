package world

import (
	"database/sql"
	"log"
	"myBot/mydb"
)

// WorldType ...
type WorldType struct {
	Game int
	Bank int
	Pot  int
	Pot2 int
}

// World ...
var World WorldType

func init() {
	World.loadWorldData()
}

func (w *WorldType) loadWorldData() {

	row := mydb.Db.QueryRow("SELECT game, bank, pot, pot2 FROM world_info limit 1")
	// defer row.Close()
	var data WorldType
	switch err := row.Scan(&data.Game, &data.Bank, &data.Pot, &data.Pot2); err {
	case sql.ErrNoRows:
		log.Println("world - No rows were returned")
	case nil:
		w.Game = data.Game
		w.Bank = data.Bank
		w.Pot = data.Pot
		w.Pot2 = data.Pot2
		log.Println("World data load.")
		//log.Println(w.Game)
	default:
		checkError(err)
	}
}

func (w *WorldType) SaveWorldData() {
	log.Println("save world data")
	//mydb.Db.QueryRow("update world_info set game = $1, bank = $2, pot = $3, pot2 = $4", w.Game, w.Bank, w.Pot, w.Pot2)
	stmt, err := mydb.Db.Prepare("update world_info set game = $1, bank = $2, pot = $3, pot2 = $4")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(w.Game, w.Bank, w.Pot, w.Pot2)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
