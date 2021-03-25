package data

// CardData ...
var CardData map[string]CardOption = make(map[string]CardOption)
var ResidentCardData map[string]CardOption = make(map[string]CardOption)
var LimitedCardData []string

func init() {
	CardData = make(map[string]CardOption)
	merge(CardData, ResidentCard, LimitedCard)
	ResidentCardData = make(map[string]CardOption)
	merge(ResidentCardData, ResidentCard)
	LimitedCardData = nil
	for id, _ := range LimitedCard {
		LimitedCardData = append(LimitedCardData, id)
	}
}

func merge(thisData map[string]CardOption, args ...map[string]CardOption) {
	for _, data := range args {
		for key, value := range data {
			thisData[key] = value
		}
	}
}
