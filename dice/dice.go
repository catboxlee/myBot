package dice

import (
	"math/rand"
	"regexp"
	"strconv"
)

// DiceType type
type diceType struct {
	Hit   int
	Rolls []int
}

// Dice ...
var Dice diceType

// Roll dice
func (d *diceType) Roll(s string) {
	//rand.Seed(time.Now().UnixNano())
	d.Hit = 0
	d.Rolls = nil

	re := regexp.MustCompile(`(\d*)d(\d*)\+?(\d*)`)
	matches := re.FindStringSubmatch(s)

	nbr, _ := strconv.Atoi(matches[1])
	sided, _ := strconv.Atoi(matches[2])
	modifiers, _ := strconv.Atoi(matches[3])

	for i := 0; i < nbr; i++ {
		d.Rolls = append(d.Rolls, rand.Intn(sided)+1)
		d.Hit += d.Rolls[i]
	}

	d.Hit += modifiers

}
