package title

var titleCodeMap = map[string]string{
	"three consecutive": "三連冠",
	"five consecutive":  "五連冠",
	"seven consecutive": "七連冠",
	"ten consecutive":   "十連冠",
	"one shot":          "一拳超人",
	"sogeking":          "狙擊之王",
	"one punch":         "一拳",
	"sec kill":          "秒殺",
	"shunsatsu":         "瞬殺",
}

// Title ...
func Title(input string) string {
	if _, exist := titleCodeMap[input]; exist {
		return titleCodeMap[input]
	}
	return ""
}
