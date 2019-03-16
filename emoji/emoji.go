package emoji

import (
	"html"
	"strconv"
)

func EmojiCode(input string) string{
	x , _ := strconv.ParseInt(input, 16, 64)
	return html.UnescapeString("&#" + strconv.Itoa(int(x)) + ";")
}

func Emoji(input string) string {
	if _, err := emojiCodeMap[input]; err == nil {
		return emojiCodeMap[input]
	}
	return nil
}