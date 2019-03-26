package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"myBot/mydb"
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
	Items       Items  `json:"items"`
}

// Items ...
type Items struct {
	Item1 string `json:"item1"`
	Item2 string `json:"item2"`
}

// CurrentUserProfile ...
type CurrentUserProfile struct {
	UserProfile *linebot.UserProfileResponse
	Event       *linebot.Event
}

// LineUser ...
var LineUser CurrentUserProfile

// GetSenderInfo ...
func (u *CurrentUserProfile) GetSenderInfo(bot *linebot.Client) {
	switch u.Event.Source.Type {
	case linebot.EventSourceTypeGroup:
		if senderProfile, err := bot.GetGroupMemberProfile(u.Event.Source.GroupID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
		}
	case linebot.EventSourceTypeRoom:
		if senderProfile, err := bot.GetRoomMemberProfile(u.Event.Source.RoomID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
		}
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
		}
		//return event.Source.UserID
	}
}

func (u *CurrentUserProfile) SaveUserData() {
	query := `insert into users(userid, displayname, money) values($1, $2, 0)
					on conflict(userid)
					do update set displayname = $2`
	mydb.Db.QueryRow(query, &u.UserProfile.UserID, &u.UserProfile.DisplayName)
}

func getJSON() {
	// Open our jsonFile
	jsonFile, err := os.Open("savedata/user.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	var users map[string]User
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &users)
	fmt.Println(users)
}

func setJSON() {

	user := User{"test1", "c", Items{}}
	users := map[string]User{"t1": user}
	jsonData, _ := json.Marshal(users)

	// sanity check
	fmt.Println(string(jsonData))

	// write to JSON file
	jsonFile, err := os.Create("savedata/user.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
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
