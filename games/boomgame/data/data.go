package data

// CardData ...
var CardData map[string]CardOption = make(map[string]CardOption)

// GachaCardData ...
var GachaCardData map[string]CardOption = make(map[string]CardOption)

func init() {
	CardData = make(map[string]CardOption)
	merge(CardData, RCard, SRCard, SSRCard, LimitedCard, MythosCard)
	GachaCardData = make(map[string]CardOption)
	merge(GachaCardData, RCard, SRCard, SSRCard)
}

func merge(thisData map[string]CardOption, args ...map[string]CardOption) {
	for _, data := range args {
		for key, value := range data {
			thisData[key] = value
		}
	}
}
