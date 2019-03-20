package player

import "fmt"

type playersType struct {
	data map[string]*playerType
}

type playerType struct {
	id    string
	name  string
	money int
}

// Players Data
var Players playersType

// Test ...
func Test() {
	fmt.Println("player loading...")
	Players.loadData()
	Players.get("catbox")
	Players.get("catbox2")
	Players.new("catbox3")
	Players.data["catbox2"].Set()
	Players.get("catbox2")

}

func (p *playersType) loadData() {
	Players.data = map[string]*playerType{
		"catbox":  {"1", "catbox", 100},
		"catbox2": {"2", "catbox2", 200},
	}
	fmt.Println(Players)
}

func (p *playerType) Set() {
	p.money = 500
}

func (p *playersType) get(username string) {

	if key, exist := p.data[username]; exist {
		fmt.Println("Found key")
		fmt.Println(key)
	} else {
		fmt.Println("Not Found")
	}
}

func (p *playersType) new(username string) {

	user := &playerType{"3", username, 300}
	p.data[username] = user
}
