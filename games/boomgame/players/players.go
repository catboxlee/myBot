package players

import (
	"encoding/json"
	"log"
	"my4/games/boomgame/data"
	"my4/games/boomgame/scheduler"
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

// LoadPlayersData ...
func (p *Players) LoadPlayersData(g scheduler.Game, sourceid string) {
	p.Parent = g
	p.SourceID = sourceid
	p.Data = make(map[string]*PlayerOption)

	// SQL load
	for _, data := range sqlDataTest() {
		p.Data[data.UserID] = p.buildPlayer(data)
	}

	log.Println("Players data load.", p.Data)
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
			CardPile: []byte(`{"cards":{"19":{"id":0,"level":0,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"sogeking"},"3":{"id":0,"level":0,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"tsubamegaeshi"},"4":{"id":0,"level":0,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"tsubamegaeshi"}}}`),
			Property: []byte(`{"winningstreak":5,"losingstreak":0}`),
		},
		{
			UserID:   "ID1",
			SourceID: "0",
			CardPile: []byte(`{"cards":{"2":{"id":0,"level":2,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"sogeking"},"3":{"id":0,"level":0,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"tsubamegaeshi"},"5":{"id":0,"level":0,"triggertimes":0,"untriggertimes":0,"quantity":0,"set":"","coreset":"tsubamegaeshi"}}}`),
			Property: []byte(`{"winningstreak":3,"losingstreak":3}`),
		},
	}
	return datas
}

func (p *Players) buildPlayer(spo sqlPlayerOption) *PlayerOption {
	np := new(PlayerOption)
	np.Parent = p
	np.SourceID = p.SourceID
	np.UserID = spo.UserID
	json.Unmarshal(spo.Titles, &np.Titles)
	//np.CardPile.Cards = make(map[string]*cards.CardOption)
	json.Unmarshal(spo.CardPile, &np.CardPile)
	np.CardPile.SetParent(np)
	for k, v := range np.CardPile.Cards {
		v.GenerateCard(&np.CardPile, np, data.CardData[k])
	}
	json.Unmarshal(spo.ItemPile, &np.ItemPile)
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
		CardPile: []byte(`{"cards":{}}`),
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
