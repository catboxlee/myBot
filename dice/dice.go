package dice

import (
	"math/rand"
	"regexp"
	"strconv"
)

// DiceType type
type DiceType struct {
	N     int
	Rolls []int
}

var Dice DiceType

// Roll dice
func (d *DiceType) Roll(s string) {

	re := regexp.MustCompile(`(\d*)d(\d*)\+?(\d*)`)
	matches := re.FindStringSubmatch(s)

	nbr, _ := strconv.Atoi(matches[1])
	sided, _ := strconv.Atoi(matches[2])
	modifiers, _ := strconv.Atoi(matches[3])

	for i := 0; i < nbr; i++ {
		d.Rolls = append(d.Rolls, rand.Intn(sided)+1)
		d.N += d.Rolls[i]
	}

	d.N += modifiers

}
