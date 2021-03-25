package typeset

// Property ...
type Property struct {
	TotalMove int `json:"total_move,omitempty"`
	Move      int `json:"move,omitempty"`
	DFace     int
	DPlus     int
	DCnt      int
	PropertyDice
	Stop   bool
	Turn   int      `json:"turn,omitempty"`
	Buff   []string `json:"buff,omitempty"`
	DeBuff []string `json:"de_buff,omitempty"`
}

type PropertyDice struct {
	DiceCnt  int
	DiceFace int
	DicePlus int
	DiceHit  int
}

func (pr *Property) GetDice() (int, int, int) {
	return pr.DiceCnt, pr.DiceFace, pr.DicePlus
}

func (pr *Property) MakeDice(n int, m int, o int) {
	pr.DiceCnt += n
	pr.DiceFace += m
	pr.DicePlus += o
}

func (pr *Property) ResetDice() {
	pr.DiceCnt = 0
	pr.DiceFace = 0
	pr.DicePlus = 0
	pr.DiceHit = 0
	pr.Stop = false
}

func (pr *Property) GetTotalMove() int {
	return pr.TotalMove
}

func (pr *Property) MakeTotalMove(n int) {
	pr.TotalMove += n
}

func (pr *Property) GetMove() int {
	return pr.Move
}

func (pr *Property) MakeMove(n int) {
	pr.Move += n
}

func (pr *Property) SetStop(n bool) {
	pr.Stop = n
}

func (pr *Property) GetTurn() int {
	return pr.Turn
}

func (pr *Property) MakeTurn() {
	pr.Turn++
}

func (pr *Property) AddDeBuff(s string) {
	pr.DeBuff = append(pr.DeBuff, s)
}

func (pr *Property) ClearDeBuff() {
	pr.DeBuff = nil
}

func (pr *Property) AddBuff(s string) {
	pr.Buff = append(pr.Buff, s)
}

func (pr *Property) ClearBuff() {
	pr.Buff = nil
}
