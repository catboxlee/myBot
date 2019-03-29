package world

import (
	"database/sql"
	"myBot/mydb"
)

// WorldType ...
type WorldType struct {
	Bank int
}

// World ...
var World *WorldType

func init() {
	World.loadWorldData()
}

func (w *WorldType) loadWorldData() {

	rows, err := mydb.Db.Query("SELECT bank FROM world_info")
	checkError(err)
	defer rows.Close()
	var data WorldType
	for rows.Next() {
		switch err := rows.Scan(&data.Bank); err {
		case sql.ErrNoRows:
			//fmt.Println("No rows were returned")
		case nil:
			w = &data
		default:
			checkError(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
