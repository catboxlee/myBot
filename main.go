// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"math/rand"
	
	"myBot/mydb"
	"myBot/world"
	"myBot/users"
	
	"net/http"
	"os"
	"strings"
	"time"
	"myBot/games/pokergoal"
	"myBot/boomgame1"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	defer mydb.Db.Close()

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/", hello)
	port := os.Getenv("PORT")
	//port = "8080"
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("!hello world!"))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	doLinebotEvents(events)
}

func doLinebotEvents(events []*linebot.Event) {

	rand.Seed(time.Now().UnixNano())
	log.Println("1111")

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			//displayName := GetSenderInfo(event)
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				users.LineUser.GetSenderInfo(event, bot)
				input := strings.TrimSpace(string(message.Text))
				var texts []string
				log.Println("2222")
				switch world.World.Game {
				case 1:
					texts = boomgame1.Boom.Run(input)
				case 2:
					texts = pokergoal.Pokergoal.Run(input)
				default:
				}
				replyMsg(event, texts)
			}
		}
	}
}

func pushMsg(event *linebot.Event, texts []string) {
	var contents []linebot.SendingMessage
	if len(texts) > 0 {
		for _, text := range texts {
			contents = append(contents, linebot.NewTextMessage(text))
		}
		if _, err := bot.PushMessage(GetSenderID(event), contents...).Do(); err != nil {
			log.Print(err)
		}
	}
}
func replyMsg(event *linebot.Event, texts []string) {
	var contents []linebot.SendingMessage
	if len(texts) > 0 {
		for _, text := range texts {
			contents = append(contents, linebot.NewTextMessage(text))
		}
		if _, err := bot.ReplyMessage(event.ReplyToken, contents...).Do(); err != nil {
			log.Print(err)
		}
	}
}

// GetSenderInfo ...
func GetSenderInfo(event *linebot.Event) {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		if senderProfile, err := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do(); err == nil {
			users.LineUser.UserProfile = senderProfile
			users.LineUser.Event = event
		}
	case linebot.EventSourceTypeRoom:
		if senderProfile, err := bot.GetRoomMemberProfile(event.Source.RoomID, event.Source.UserID).Do(); err == nil {
			users.LineUser.UserProfile = senderProfile
			users.LineUser.Event = event
		}
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(event.Source.UserID).Do(); err == nil {
			users.LineUser.UserProfile = senderProfile
			users.LineUser.Event = event
		}
		//return event.Source.UserID
	}
	//user.LineUser.SaveUserData()
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
