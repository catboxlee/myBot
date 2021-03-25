package players

import (
	"encoding/json"
	"myBot/games/racegame/scheduler"
	"myBot/games/racegame/typeset"
)

// Players ...
type Players struct {
	Parent   scheduler.Game
	SourceID string
	Data     map[string]*PlayerOption
}

// GetTopParent ...
func (p *Players) GetTopParent() scheduler.Game {
	return p.GetParent()
}

// GetParent ...
func (p *Players) GetParent() scheduler.Game {
	return p.Parent
}

// ClearPlayers ...
func (p *Players) ClearPlayers() {
	p.Data = make(map[string]*PlayerOption)
}

// LoadPlayersData ...
func (p *Players) LoadPlayersData(g scheduler.Game, sourceid string) {
	p.Parent = g
	p.SourceID = sourceid
	p.Data = make(map[string]*PlayerOption)

}

type sqlPlayerOption struct {
	UserID   string
	SourceID string
	Titles   json.RawMessage
	CardPile json.RawMessage
	ItemPile json.RawMessage
	Property json.RawMessage
}

func sqlDataTest() []sqlPlayerOption {
	datas := []sqlPlayerOption{
		{
			UserID:   "ID0",
			SourceID: "0",
			CardPile: []byte(`[]`),
		},
		{
			UserID:   "ID1",
			SourceID: "0",
			CardPile: []byte(`[]`),
		},
	}
	return datas
}

func (p *Players) buildPlayer(spo sqlPlayerOption) *PlayerOption {
	np := new(PlayerOption)
	np.Parent = p
	np.SourceID = p.SourceID
	np.UserID = spo.UserID
	//np.CardPile.Cards = make(map[string]*cards.CardOption)
	json.Unmarshal(spo.CardPile, &np.CardPile)
	np.Property = new(typeset.Property)
	json.Unmarshal(spo.Property, &np.Property)
	return np
}

// CheckPlayerExist ...
func (p *Players) CheckPlayerExist(userid string) {
	if _, exist := p.Data[userid]; !exist {
		p.Data[userid] = p.newPlayer(userid)
		p.Data[userid].SaveData()
	}
}

func (p *Players) newPlayer(userid string) *PlayerOption {
	data := sqlPlayerOption{
		UserID:   userid,
		SourceID: p.SourceID,
		CardPile: []byte(`{}`),
	}
	return p.buildPlayer(data)
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
