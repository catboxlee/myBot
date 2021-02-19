package users

import (
	"database/sql"
	"log"
	"myBot/mydb"

	"github.com/line/line-bot-sdk-go/linebot"
)

// UserDataType ...
type UserDataType struct {
	UserID        string
	DisplayName   string
	Money         int
	SwallowReturn int
}

// UsersType ...
type UsersType struct {
	Data map[string]*UserDataType
}

// CurrentUserProfile ...
type CurrentUserProfile struct {
	UserProfile *linebot.UserProfileResponse
	Event       *linebot.Event
}

// UsersList ...
var UsersList UsersType

// LineUser ...
var LineUser CurrentUserProfile

func init() {
	UsersList.loadUsersData()
}

func (u *UsersType) loadUsersData() {

	u.Data = make(map[string]*UserDataType)
	rows, err := mydb.Db.Query("SELECT userid, displayname, money, swallowreturn FROM users")
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var data UserDataType
		switch err := rows.Scan(&data.UserID, &data.DisplayName, &data.Money, &data.SwallowReturn); err {
		case sql.ErrNoRows:
			log.Println("No rows were returned")
		case nil:
			u.Data[data.UserID] = &data
		default:
			checkError(err)
		}
	}
	log.Println("Users data load.")
	//log.Println(u.Data)
}

// SaveUserData ...
func (u *CurrentUserProfile) SaveUserData() {

	stmt, err := mydb.Db.Prepare(`insert into users(userid, displayname, money, swallowreturn) values($1, $2, $3, $4)
	on conflict(userid)
	do update set displayname = $2, money = $3, swallowreturn = $4`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.UserProfile.UserID, u.UserProfile.DisplayName, UsersList.Data[u.UserProfile.UserID].Money, UsersList.Data[u.UserProfile.UserID].SwallowReturn)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()
}

func (u *CurrentUserProfile) checkUserExist() {
	if _, exist := UsersList.Data[u.UserProfile.UserID]; !exist {
		UsersList.Data[u.UserProfile.UserID] = &UserDataType{}
		UsersList.Data[u.UserProfile.UserID].UserID = u.UserProfile.UserID
		UsersList.Data[u.UserProfile.UserID].DisplayName = u.UserProfile.DisplayName
		UsersList.Data[u.UserProfile.UserID].Money = 10
		UsersList.Data[u.UserProfile.UserID].SwallowReturn = 0
		u.SaveUserData()
	}
}

// GetSenderInfo ...
func (u *CurrentUserProfile) GetSenderInfo(event *linebot.Event, bot *linebot.Client) {
	u.Event = event
	switch u.Event.Source.Type {
	case linebot.EventSourceTypeGroup:
		if senderProfile, err := bot.GetGroupMemberProfile(u.Event.Source.GroupID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			u.checkUserExist()
		}
	case linebot.EventSourceTypeRoom:
		if senderProfile, err := bot.GetRoomMemberProfile(u.Event.Source.RoomID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			u.checkUserExist()
		}
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			u.checkUserExist()
		}
		//return event.Source.UserID
	}
}

// GetSenderID - Get event sender's id
func GetSenderID(event *linebot.Event) string {
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
