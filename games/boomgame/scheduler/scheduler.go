package scheduler

// Game ...
type Game interface {
	GetQueue() []string
	GetPlayQueue() []string
	AddPlayQueue(string)
	OnPlay()
	GetHit() (int, int, int)
	SetHit(int, int, int)
	GetInfoCurrent() int
	SetInfoCurrent(int)
	SetInfoRange(int, int)
	GetInfoBoomCnt() int
	MakeInfoBoomCnt(int)
	GetRankBoomCnt(string) int
	SetRankBoomCnt(string, int)
	MakeRankBoomCnt(string, int)
	GamePhase(string)
	GetPlayer(string) Player
	Show()
	GaCha(string, string) string
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
	GetCardPile() Cards
	GetRandCards(int) []Card
	SaveData()
}

// Cards ...
type Cards interface {
	GetTopParent() Game
	GetParent() Player
	GetCards() map[string]Card
	TakeCard(string) string
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
