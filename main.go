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
	"net/http"
	"os"
	"strings"
	"time"

	"myBot/boomgame1"
	"myBot/user"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {

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

	rand.Seed(time.Now().UnixNano())

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			//displayName := GetSenderInfo(event)
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				input := strings.TrimSpace(string(message.Text))
				texts := boomgame1.Boom.Run(input)
				var contents string
				userProfile := GetSenderInfo(event)
				user.LineUser.UserProfile = userProfile
				user.LineUser.Event = event
				contents = user.LineUser.UserProfile.DisplayName
				for _, text := range texts {
					contents += text + "\n"
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(contents+" OK!\n")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

// GetSenderInfo ...
func GetSenderInfo(event *linebot.Event) *linebot.UserProfileResponse {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		if senderProfile, err := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do(); err == nil {
			return senderProfile
		}
	case linebot.EventSourceTypeRoom:
		if senderProfile, err := bot.GetRoomMemberProfile(event.Source.RoomID, event.Source.UserID).Do(); err == nil {
			return senderProfile
		}
	case linebot.EventSourceTypeUser:
		if senderProfile, err := bot.GetProfile(event.Source.UserID).Do(); err == nil {
			return senderProfile
		}
		//return event.Source.UserID
	}
	var senderProfile *linebot.UserProfileResponse
	return senderProfile
}
