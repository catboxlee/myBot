package emoji

import (
	"html"
	"strconv"
)

// EmojiCode ...
func EmojiCode(input string) string {
	x, _ := strconv.ParseInt(input, 16, 64)
	return html.UnescapeString("&#" + strconv.Itoa(int(x)) + ";")
}

// Emoji ...
func Emoji(input string) string {
	if _, exist := emojiCodeMap[input]; exist {
		return emojiCodeMap[input]
	}
	return ""
}
