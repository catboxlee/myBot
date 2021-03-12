package users

import (
	"database/sql"
	"log"
	"myBot/mydb"

	"github.com/line/line-bot-sdk-go/linebot"
)

// TypeUsers ...
type TypeUsers struct {
	Data map[string]*UserOption
}

// CurrentUserProfile ...
type CurrentUserProfile struct {
	UserProfile *linebot.UserProfileResponse
	Event       *linebot.Event
}

// UsersList ...
var UsersList TypeUsers

// LineUser ...
var LineUser CurrentUserProfile

func init() {
	LineUser.UserProfile = &linebot.UserProfileResponse{}
	UsersList.loadUsersData()
}

// User ...
func (u *TypeUsers) User(thisUserID string) *UserOption {
	if _, exist := u.Data[thisUserID]; exist {
		return u.Data[thisUserID]
	}
	return nil
}

func (u *TypeUsers) loadUsersData() {
	u.Data = make(map[string]*UserOption)
	rows, err := mydb.Db.Query("SELECT userid, displayname, money, gemstone FROM users")
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		var data UserOption
		switch err := rows.Scan(&data.UserID, &data.DisplayName, &data.Money, &data.GemStone); err {
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

// CheckUserExist ...
func (u *TypeUsers) CheckUserExist(userProfile *linebot.UserProfileResponse) {
	u.checkUserExist(userProfile)
}

func (u *TypeUsers) checkUserExist(userProfile *linebot.UserProfileResponse) {
	if _, exist := u.Data[userProfile.UserID]; !exist {
		log.Println("New User:", userProfile.UserID)
		u.Data[LineUser.UserProfile.UserID] = u.newUser(userProfile)
		u.Data[LineUser.UserProfile.UserID].addData()
	} else if u.Data[LineUser.UserProfile.UserID].GetDisplayName() != LineUser.UserProfile.DisplayName {
		u.Data[LineUser.UserProfile.UserID].setDisplayName(LineUser.UserProfile.DisplayName)
	}
}

func (u *TypeUsers) newUser(userProfile *linebot.UserProfileResponse) *UserOption {
	nu := new(UserOption)
	nu.UserID = userProfile.UserID
	nu.DisplayName = userProfile.DisplayName
	nu.Money = 100
	nu.GemStone = 13900
	return nu
}

// SaveData ...
func (u *TypeUsers) SaveData(userid string) {
	u.Data[userid].saveData()
}

// DeleteUserData ...
func (u *TypeUsers) DeleteUserData() {

}

// GetSenderInfo ...
func (u *CurrentUserProfile) GetSenderInfo(event *linebot.Event, bot *linebot.Client) {
	u.Event = event
	switch u.Event.Source.Type {
	case linebot.EventSourceTypeGroup:
		if senderProfile, err := bot.GetGroupMemberProfile(u.Event.Source.GroupID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			UsersList.checkUserExist(u.UserProfile)
		}
	case linebot.EventSourceTypeRoom:
		if senderProfile, err := bot.GetRoomMemberProfile(u.Event.Source.RoomID, u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			UsersList.checkUserExist(u.UserProfile)
		}
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(u.Event.Source.UserID).Do(); err == nil {
			u.UserProfile = senderProfile
			UsersList.checkUserExist(u.UserProfile)
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
