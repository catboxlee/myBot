package cards

import (
	"math/rand"
)

// Shuffle 洗牌
func Shuffle(vals []*CardOption) {
	//rand.Seed(time.Now().UnixNano())
	for len(vals) > 0 {
		n := len(vals)                                          // 陣列長度
		randIndex := rand.Intn(n)                               // 取隨機index
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1] // 將最後一張牌和第randIndex張牌互換
		vals = vals[:n-1]
	}
}
