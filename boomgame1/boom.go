package boomgame1

import (
	"fmt"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/user"
	"reflect"
	"strconv"
	"strings"
)

type boomType struct {
	hit     int
	current int
	min     int
	max     int
}

type rankType struct {
	userID      string
	displayName string
	boom        int
}

// Boom ...
var Boom boomType
var texts []string

// Run ...
func (b *boomType) Run(input string) []string {
	if b.hit == 0 {
		b.reset()
	}
	texts = nil

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		b.checkCommand(strings.TrimLeft(input, "/"))
	} else if x, err := strconv.Atoi(input); err == nil {
		// 數字 - 檢查炸彈
		b.checkBoom(x)
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
			b.boomUser()
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

func (b *boomType) boomUser() {
	if userInfo := user.GetSenderInfo(); userInfo != nil {
		texts = append(texts, fmt.Sprintf("%s", reflect.TypeOf(userInfo).String()))
	}
}

func (b *boomType) rank() {

}

func (b *boomType) reset() {
	hit := &dice.Dice
	hit.Roll("1d100")
	b.hit = hit.N
	b.current = 0
	b.min = 0
	b.max = 101
}

func (b *boomType) show() {
	if b.current == b.hit {
		texts = append(texts, fmt.Sprintf("%s %d", emoji.Emoji(":collision:"), b.hit))
	} else {
		texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.min), emoji.Emoji(":bomb:"), helper.Min(100, b.max)))
	}
}

func (b *boomType) save() {

}

func loadProfile(userID string) {

}
