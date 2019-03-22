package boomgame1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"myBot/dice"
	"myBot/emoji"
	"myBot/helper"
	"myBot/user"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type gameType struct {
	hit     int
	current int
	min     int
	max     int
}

type rankType struct {
	Season   int `json:"season"`
	rankUser map[string]*rankUser
}

type rankUser struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Boom        int    `json:"boom"`
}

var rank rankType

// Boom ...
var Boom gameType
var texts []string

func init() {

	rank.rankUser = make(map[string]*rankUser)
	rank.loadRank()
}

// Run ...
func (b *gameType) Run(input string) []string {
	rank.Season = helper.Max(12, rank.Season)
	if b.hit == 0 {
		b.reset()
	}
	texts = nil

	if strings.HasPrefix(input, "/") {
		// 字串 - 執行指令
		b.checkCommand(strings.TrimLeft(input, "/"))
		return texts
	}

	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) > 1 {
		if x, err := strconv.Atoi(matches[1]); err == nil {
			// 數字 - 檢查炸彈
			b.checkBoom(x)
			return texts
		}
	}
	return texts
}

func (b *gameType) checkCommand(input string) {
	switch input {
	case "reset":
		b.reset()
		b.show()
	case "rank":
		rank.rank()
		b.show()
	case "resetRank":
		rank.resetRank()
		rank.saveRank()
	}
}

func (b *gameType) checkBoom(x int) {
	if x > b.min && x < b.max {
		b.current = x
		switch {
		case b.current == b.hit:
			b.show()
			rank.addUserBoom()
			rank.rank()
			rank.checkBoomKing()
			rank.saveRank()
			b.reset()
			b.show()
		case b.current < b.hit:
			b.min = b.current
			b.show()
		case b.current > b.hit:
			b.max = b.current
			b.show()
		}
	}
}

func (b *gameType) reset() {
	hit := &dice.Dice
	hit.Roll("1d100")
	b.hit = hit.N
	b.current = 0
	b.min = 0
	b.max = 101
}

func (b *gameType) show() {
	if b.current == b.hit {
		texts = append(texts, fmt.Sprintf("%s %s %d", user.LineUser.UserProfile.DisplayName, emoji.Emoji(":collision:"), b.hit))
	} else {
		texts = append(texts, fmt.Sprintf("%d - %s - %d", helper.Max(1, b.min), emoji.Emoji(":bomb:"), helper.Min(100, b.max)))
	}
}

func (r *rankType) addUserBoom() {

	if len(r.rankUser) == 0 {
		r.rankUser = map[string]*rankUser{user.LineUser.UserProfile.UserID: {user.LineUser.UserProfile.UserID, user.LineUser.UserProfile.DisplayName, 1}}
	} else if _, exist := r.rankUser[user.LineUser.UserProfile.UserID]; exist {
		r.rankUser[user.LineUser.UserProfile.UserID].Boom++
	} else {
		r.rankUser[user.LineUser.UserProfile.UserID] = &rankUser{user.LineUser.UserProfile.UserID, user.LineUser.UserProfile.DisplayName, 1}

	}
}

func (r *rankType) checkBoomKing() {
	if _, exist := r.rankUser[user.LineUser.UserProfile.UserID]; exist {
		if r.rankUser[user.LineUser.UserProfile.UserID].Boom >= 100 {
			texts = append(texts, fmt.Sprintf("%sS%d 爆爆王：%s%s", emoji.Emoji(":confetti_ball:"), r.Season, user.LineUser.UserProfile.DisplayName, emoji.Emoji(":confetti_ball:")))
			r.Season++
			r.rankUser = nil
		}
	} else {
		texts = append(texts, fmt.Sprintf("UserID:%s 不存在", user.LineUser.UserProfile.UserID))
	}
}

func (r *rankType) rank() {

	tmpRank := make([]*rankUser, 0, len(r.rankUser))
	for _, val := range r.rankUser {
		tmpRank = append(tmpRank, val)
	}

	sort.SliceStable(tmpRank, func(i, j int) bool {
		return tmpRank[i].Boom > tmpRank[j].Boom
	})

	text := fmt.Sprintf("爆爆王 S%d Rank：", r.Season)
	for _, v := range tmpRank {
		text += fmt.Sprintf("\n%s %s %d", v.DisplayName, emoji.Emoji(":collision:"), v.Boom)
	}
	texts = append(texts, text)
}

func (r *rankType) resetRank() {
	rank.rankUser = make(map[string]*rankUser)
}

func (r *rankType) saveRank() {

	jsonData, _ := json.Marshal(rank)

	// sanity check
	//fmt.Println(string(jsonData))

	// write to JSON file
	jsonFile, err := os.Create("savedata/common/boomRank.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	log.Println("JSON data written to ", jsonFile.Name())
}

func (r *rankType) loadRank() {
	// Open our jsonFile
	jsonFile, err := os.Open("savedata/common/boomRank.json")
	// if we os.Open returns an error then handle it
	if err != nil && os.IsNotExist(err) {
		//log.Println(err)
		jsonFile, _ = os.Create("savedata/common/boomRank.json")
		log.Println("JSON data create : ", jsonFile.Name())
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		if len(byteValue) > 0 {
			json.Unmarshal(byteValue, r)
		}
		log.Println("JSON data load : ", jsonFile.Name())
	}
	defer jsonFile.Close()
}
