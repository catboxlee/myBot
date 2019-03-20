package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Users ...
type Users struct {
	Users []User `json:"users"`
}

// User ...
type User struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
	Items       Items  `json:"Items"`
}

// Items ...
type Items struct {
	Item1 string `json:"item1"`
	Item2 string `json:"item2"`
}

type onlineType struct {
	userID      string
	displayName string
}

var bot *linebot.Client

// Event ...
var Event *linebot.Event
var onlineUsers map[string]onlineType

// LoadUsersData ...
func LoadUsersData() {
	// Open our jsonFile
	jsonFile, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var users Users

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)
}

// GetSenderInfo ...
func GetSenderInfo() interface{} {
	switch Event.Source.Type {
	case linebot.EventSourceTypeGroup:
		return Event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		return Event.Source.RoomID
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(Event.Source.UserID).Do(); err == nil {
			return senderProfile
		}
		//return event.Source.UserID
	}
	return nil
}

// GetSenderID - Get event sender's id
func GetSenderID(c context.Context, event *linebot.Event) string {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID
	case linebot.EventSourceTypeUser:
		return event.Source.UserID
	}
	log.Printf("Can not get sender id. type: %v\n", event.Source.Type)
	return ""
}

// GetSenderName ...
/* 送信者の表示名を取得する
 *
 * ユーザしか取得できないので、ルームおよびグループではidをそのまま返す
 * グループメンバーのUserIDの場合、そのユーザが直接Botと友だち登録していなければ取得できない
 */
func GetSenderName(c context.Context, bot *linebot.Client, from string) string {
	if len(from) == 0 {
		log.Println(c, "Parameter `mid` was not specified.")
		return from
	}
	if from[0:1] == "U" {
		senderProfile, err := bot.GetProfile(from).Do()
		if err != nil {
			log.Println(c, "Error occurred at get sender profile. from: %v, err: %v", from, err)
			return from
		}
		return senderProfile.DisplayName
	}
	return from
}
