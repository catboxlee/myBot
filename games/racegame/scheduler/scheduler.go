package scheduler

import (
	"myBot/games/racegame/typeset"
)

// Game ...
type Game interface {
	GetQueue() []string
	GetPlayQueue() []string
	AddPlayQueue(string)
	GetMeter() int
	GetRankingArray() []string
	GetRanking(string) int
	PopCard(Player)
	OnPlay()
	GamePhase(string)
	GetPlayer(string) Player
	Show()
	ViewCardsInfo(string) string
}

// Players ...
type Players interface {
	GetTopParent() Game
	GetParent() Game
	SaveData(string)
}

// Player ...
type Player interface {
	GetTopParent() Game
	GetParent() Players
	GetUserID() string
	GetDisplayName() string
	GetGemStone() int
	MakeGemStone(int)
	GetProperty() *typeset.Property
	TakeCard(id string)
	RemoveCardPile(...string) (bool, string)
	SaveData()
}

// Cards ...
type Cards interface {
	GetTopParent() Game
	GetParent() Game
	GetCards() map[string]Card
	UsedCard(string)
}

// Card ...
type Card interface {
	GetTopParent() Game
	GetParent() Cards
	GetDisplayName() string
	GetLevel() int
	SetDesc(string)
	GetCoolDown() int
	MakeCoolDown(int)
	GetReCoolDown() int
	ResetCoolDown()
	SetCoolDown(int)
	GetFreeze() int
	MakeFreeze(n int)
}
