package world

import (
	"database/sql"
	"myBot/mydb"
	"log"
)

// WorldType ...
type WorldType struct {
	Game int
	Bank int
}

// World ...
var World WorldType

func init() {
	World.loadWorldData()
}

func (w *WorldType) loadWorldData() {

	row := mydb.Db.QueryRow("SELECT game, bank FROM world_info limit 1")
	// defer row.Close()
	var data WorldType
		switch err := row.Scan(&data.Game, &data.Bank); err {
		case sql.ErrNoRows:
			log.Println("world - No rows were returned")
		case nil:
			w = data
			log.Println("World data load.")
			//log.Println(w.Game)
		default:
			checkError(err)
		}
}

func (w *WorldType) SaveWorldData() {
	log.Println("save world data")
	mydb.Db.QueryRow("update world_info set game = $1, bank = $2", w.Game, w.Bank)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
