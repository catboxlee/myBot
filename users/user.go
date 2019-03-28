package users

import (
	"database/sql"
	"log"
	"myBot/mydb"

	"github.com/line/line-bot-sdk-go/linebot"
)

// UserDataType ...
type UserDataType struct {
	UserID      string
	DisplayName string
	Money       int
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
	rows, err := mydb.Db.Query("SELECT userid, displayname, money FROM users")
	checkError(err)
	defer rows.Close()
	var data UserDataType
	for rows.Next() {
		switch err := rows.Scan(&data.UserID, &data.DisplayName, &data.Money); err {
		case sql.ErrNoRows:
			//fmt.Println("No rows were returned")
		case nil:
			u.Data[data.UserID] = &data
		default:
			checkError(err)
		}
	}
}

func (u *CurrentUserProfile) sveUserData() {
	query := `insert into users(userid, displayname, money) values($1, $2, 0)
					on conflict(userid)
					do update set displayname = $2`
	mydb.Db.QueryRow(query, u.UserProfile.UserID, u.UserProfile.DisplayName)
}

func (u *CurrentUserProfile) checkUserExist() {
	if _, exist := UsersList.Data[u.UserProfile.UserID]; !exist {
		UsersList.Data[u.UserProfile.UserID] = &UserDataType{}
		UsersList.Data[u.UserProfile.UserID].UserID = u.UserProfile.UserID
		UsersList.Data[u.UserProfile.UserID].DisplayName = u.UserProfile.DisplayName
		UsersList.Data[u.UserProfile.UserID].Money = 10
		u.sveUserData()
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
