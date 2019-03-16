package boomgame1

import (
	"myBot/dice"
	"myBot/helper"
	"fmt"
	"strconv"
	"strings"
	. "myBot/emoji"
)

type boomType struct {
	hit     int
	current int
	min     int
	max     int
}

type rankType struct {
	userId   string
	displayName string
	boom     int
}

var boom boomType
var texts []string

func init() {
	boom.reset()
	boom.show()
}

// Run ...
func Run(input string) []string {
	texts = nil

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		boom.checkCommand(strings.TrimLeft(input, "/"))
	} else if x, err := strconv.Atoi(input); err == nil {
		// 數字 - 檢查炸彈
		boom.checkBoom(x)
	}
	return texts
}

func (b *boomType) checkCommand(input string) {
	switch input {
	case "reset":
		b.reset()
		b.show()
	case "rank":
		b.rank()
		b.show()

	}
}

func (b *boomType) checkBoom(x int) {
	if x > b.min && x < b.max {
		b.current = x
		switch {
		case b.current == b.hit:
			
			b.show()
			b.rank()
			b.reset()
			b.show()
			b.save()
		case b.current < b.hit:
			b.min = b.current
			b.show()
		case b.current > b.hit:
			b.max = b.current
			b.show()
		}
	}
}

func (b *boomType) rank() {

}

func (b *boomType) reset() {
	hit := dice.Dice
	hit.Roll("1d100")
	b.hit = hit.N
	b.current = 0
	b.min = 0
	b.max = 101
}

func (b *boomType) show() {
	if b.current == b.hit {
		texts = append(texts, fmt.Sprintf("%s %d",Emoji(":collision:"), b.hit))
	} else {
		texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.min), Emoji(":bomb:"), helper.Min(100, b.max)))
	}
}

func (b *boomType) save() {

}
