package world

import (
	"database/sql"
	"log"
	"myBot/mydb"
)

// CfgType ...
type CfgType struct {
	SourceID string
	Game     int
}

// ConfigsData ...
var ConfigsData = make(map[string]*CfgType)

func init() {

}

// LoadConfigData ...
func LoadConfigData(sourceid string) {
	if _, exist := ConfigsData[sourceid]; exist {
		return
	}
	rows, err := mydb.Db.Query("SELECT sourceid, game FROM base_config where sourceid = $1 limit 1", sourceid)
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var data CfgType
		switch err := rows.Scan(&data.SourceID, &data.Game); err {
		case sql.ErrNoRows:
			log.Println("No rows were returned")
		case nil:
			ConfigsData[data.SourceID] = &data
			log.Println("Config data load.", ConfigsData[data.SourceID])
		default:
			checkError(err)
		}
	}
	if _, exist := ConfigsData[sourceid]; !exist {
		ConfigsData[sourceid] = &CfgType{sourceid, 1}
		ConfigsData[sourceid].addNewConfigData()
	}
}

func (c *CfgType) addNewConfigData() {
	stmt, err := mydb.Db.Prepare("insert into base_config (sourceid, game) values ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(c.SourceID, c.Game)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

// UpdateConfigData ...
func (c *CfgType) UpdateConfigData() {
	stmt, err := mydb.Db.Prepare("update base_config set game = $1 where sourceid = $2")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(c.Game, c.SourceID)
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
