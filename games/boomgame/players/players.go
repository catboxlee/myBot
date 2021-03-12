package players

import (
	"database/sql"
	"encoding/json"
	"log"
	"myBot/games/boomgame/cards"
	"myBot/games/boomgame/data"
	"myBot/games/boomgame/scheduler"
	"myBot/mydb"
)

// Players ...
type Players struct {
	Parent   scheduler.Game
	SourceID string
	Data     map[string]*PlayerOption
}

// Reset ...
func (p *Players) Reset(g scheduler.Game, SourceID string) {
	p.Parent = g
	p.SourceID = SourceID
	p.Data = make(map[string]*PlayerOption)
}

// GetTopParent ...
func (p *Players) GetTopParent() scheduler.Game {
	return p.GetParent()
}

// GetParent ...
func (p *Players) GetParent() scheduler.Game {
	return p.Parent
}

// LoadPlayersData ...
func (p *Players) LoadPlayersData(g scheduler.Game, sourceid string) {
	p.Parent = g
	p.SourceID = sourceid
	p.Data = make(map[string]*PlayerOption)

	rows, err := mydb.Db.Query("SELECT userid, sourceid,cardpile, property FROM boomplayer")
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var data PlayerOption
		switch err := rows.Scan(&data.UserID, &data.SourceID, &data.CardPile, &data.Property); err {
		case sql.ErrNoRows:
			log.Println("No rows were returned")
		case nil:
			p.Data[data.UserID] = p.buildPlayer(&data)
		default:
			checkError(err)
		}
	}

	for _, data := range p.Data {
		p.Data[data.UserID] = p.buildPlayer(data)
	}
	log.Println("Players data load.", p.Data)
}

type sqlPlayerOption struct {
	UserID   string
	SourceID string
	Titles   json.RawMessage
	CardPile json.RawMessage
	Property json.RawMessage
}

func (p *Players) buildPlayer(spo *PlayerOption) *PlayerOption {
	np := new(PlayerOption)
	np.Parent = p
	np.CardPile.SetParent(np)
	np.SourceID = p.SourceID
	np.UserID = spo.UserID
	//json.Unmarshal(spo.Titles, &np.Titles)
	np.Titles = spo.Titles
	//json.Unmarshal(spo.CardPile, &np.CardPile)
	np.CardPile.Cards = make(map[string]*cards.CardOption)
	if len(spo.CardPile.Cards) > 0 {
		np.CardPile = spo.CardPile
		for k, v := range np.CardPile.Cards {
			v.GenerateCard(&np.CardPile, np, data.CardData[k])
		}
	}
	//json.Unmarshal(spo.Property, &np.Property)
	np.Property = spo.Property
	return np
}

// CheckPlayerExist ...
func (p *Players) CheckPlayerExist(userid string) {
	if _, exist := p.Data[userid]; !exist {
		p.Data[userid] = p.newPlayer(userid)
		p.Data[userid].addData()
	}
}

func (p *Players) newPlayer(userid string) *PlayerOption {
	data := PlayerOption{
		UserID:   userid,
		SourceID: p.SourceID,
	}
	return p.buildPlayer(&data)
}

// Player ...
func (p *Players) Player(thisPlayerID string) *PlayerOption {
	if _, exist := p.Data[thisPlayerID]; exist {
		return p.Data[thisPlayerID]
	}
	return nil
}

// SaveData ...
func (p *Players) SaveData(thisPlayerID string) {
	if _, exist := p.Data[thisPlayerID]; exist {
		p.Data[thisPlayerID].SaveData()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
