package racegame

import (
	"fmt"
	"log"
	"math/rand"
	"myBot/dice"
	"myBot/emoji"
	"myBot/games/racegame/players"
	"myBot/games/racegame/scheduler"
	"myBot/helper"
	"strconv"
	"strings"
)

func (g *GameType) resetInfo() {
	g.Info = &InfoType{
		Turn:      0,
		Phase:     false,
		Meter:     g.defaultMeter,
		GameOver:  false,
		Queue:     []string{},
		PlayQueue: []string{},
	}
	log.Println("Info", g.Info)
}

func (g *GameType) showGameInfo() string {
	var strs []string

	if g.Info.Phase == false {
		strs = append(strs, fmt.Sprintf("[[%s賽羚娘 %s%d]]", emoji.Emoji(":game_die:"), emoji.Emoji(":chequered_flag:"), g.Info.Meter))
		var skills []string
		for i := 1; i <= 5; i++ {
			id := strconv.Itoa(i)
			skills = append(skills, fmt.Sprintf("%s.%s ", id, g.mythosCards.Cards[id].GetDisplayName()))
		}
		if len(skills) > 0 {
			strs = append(strs, strings.Join(skills, ","))
		}
		horses := g.Info.Queue
		for _, userID := range horses {
			strs = append(strs, fmt.Sprintf("%s%s 等待中...", g.Player(userID).GetDisplayName(), func() string {
				if len(g.Player(userID).Buff) > 0 {
					return fmt.Sprintf("(%s)", g.mythosCards.Cards[g.Player(userID).Buff[0]].GetDisplayName())
				}
				return ""
			}()))
		}
		strs = append(strs, fmt.Sprintf("\n[/m] 設定跑道距離"))
		strs = append(strs, fmt.Sprintf("[+] 進入賽場"))
		if len(g.Info.Queue) > 0 {
			strs = append(strs, fmt.Sprintf("[-] 開始比賽"))
		}

		return strings.Join(strs, "\n")
	}

	// 賽場狀態
	gameTurn := (g.Info.Turn)/len(g.Info.Queue) + 1
	strs = append(strs, fmt.Sprintf("[[%s賽羚娘 %s%d 回合%d]]", emoji.Emoji(":game_die:"), emoji.Emoji(":chequered_flag:"), g.Info.Meter, gameTurn))
	g.Info.GameOver = true
	horses := g.Info.Queue
	for i, userID := range horses {
		if g.Player(userID).TotalMove >= g.Info.Meter {
			strs = append(strs, fmt.Sprintf("%s 抵達終點.", g.Player(userID).GetDisplayName()))
		} else {
			g.Info.GameOver = false
			rk := g.GetRanking(userID)
			runLine := strings.Repeat(">", len(g.Info.Queue)-rk)
			strs = append(strs, fmt.Sprintf("%d.%s %s%s%d(%+d)", i+1, g.Player(userID).GetDisplayName(), emoji.Emoji(":footprints:"), runLine, g.Player(userID).GetTotalMove(), g.Player(userID).GetMove()))
		}

	}
	if g.Info.GameOver {
		strs = append(strs, fmt.Sprintf("比賽結束."))
		return strings.Join(strs, "\n")

	}

	// 當前馬匹
	strs = append(strs, fmt.Sprintf(""))
	g.makeGameTurn()
	nextPlayer := g.Player(g.Info.Queue[(g.Info.Turn)%len(g.Info.Queue)])
	if nextPlayer.GetTurn() < gameTurn {
		if len(nextPlayer.CardPile) >= 5 {
			rnd := rand.Intn(len(nextPlayer.CardPile))
			nextPlayer.CardPile = append(nextPlayer.CardPile[:rnd], nextPlayer.CardPile[rnd+1:]...)
		}
		g.PopCard(nextPlayer)
		nextPlayer.MakeTurn()
	}
	strs = append(strs, fmt.Sprintf("%s 行動:", nextPlayer.GetDisplayName()))
	var tmp []string
	if len(nextPlayer.Buff) > 0 {
		for _, cid := range nextPlayer.Buff {
			if _, exist := g.mythosCards.Cards[cid]; exist {
				tmp = append(tmp, g.mythosCards.Cards[cid].GetDisplayName())
				g.mythosCards.Cards[cid].OnEffectFunc(nextPlayer)
			}
		}
	}
	if len(nextPlayer.DeBuff) > 0 {
		for _, cid := range nextPlayer.DeBuff {
			if _, exist := g.mythosCards.Cards[cid]; exist {
				tmp = append(tmp, g.mythosCards.Cards[cid].GetDisplayName())
			}
		}
	}
	if len(tmp) > 0 {
		strs = append(strs, strings.Join(tmp, ""))
	}
	if nextPlayer.GetTurn() > 1 {
		strs = append(strs, fmt.Sprintf("%s", nextPlayer.ViewCardsInfo()))
	}
	a, b, c := nextPlayer.GetDice()
	strs = append(strs, fmt.Sprintf("0.前進(%s%d~%d%+d)", emoji.Emoji(":game_die:"), 1+a, 6+b, c))
	nextPlayer.ResetDice()
	log.Println(g.mythosCards.Deck)
	return strings.Join(strs, "\n")
}

func (g *GameType) makeGameTurn() {
	nextPlayer := g.Player(g.Info.Queue[(g.Info.Turn)%len(g.Info.Queue)])
	if nextPlayer.TotalMove >= g.Info.Meter {
		g.Info.Turn++
		g.makeGameTurn()
	}
}
func (g *GameType) runPhase(input string) {
	var strs []string
	input = strings.ToLower(input)
	switch input[0:1] {
	case "r":
		// Run
		g.onRun()
	default:
		// Use Card.
		matches := strings.Fields(input)
		if len(matches) > 0 {
			if x, err := strconv.Atoi(matches[0]); err == nil {
				x--
				// Run
				switch x {
				case -1:
					thisPlayer := g.Player(g.Info.currentUserID)
					thisPlayer.Move = 0
					if s := g.onRun(); len(s) > 0 {
						strs = append(strs, s)
					}
					texts = append(texts, strings.Join(strs, "\n"))
					g.Show()
					if g.Info.GameOver {
						g.reset()
						g.Show()
					}
					return
				default:
					// Use Card
					thisPlayer := g.Player(g.Info.currentUserID)
					if len(thisPlayer.CardPile) > x {
						cid := thisPlayer.CardPile[x]
						thisCard := g.mythosCards.Cards[cid]
						if thisCard.OnPlayFunc != nil {
							g.OnPlay()
							var target scheduler.Player
							if len(matches) > 1 {
								if poid, err := strconv.Atoi(matches[1]); err == nil {
									if poid > 0 {
										target = g.Player(g.Info.Queue[poid-1])
									}
								}
							} else {
								target = nil
							}
							if _, s := thisCard.OnPlayFunc(thisPlayer, target); len(s) > 0 {
								strs = append(strs, s)
								thisPlayer.CardPile = append(thisPlayer.CardPile[:x], thisPlayer.CardPile[x+1:]...)
							}
						}
					}
					if s := g.onRun(); len(s) > 0 {
						strs = append(strs, s)
					}
					texts = append(texts, strings.Join(strs, "\n"))
					g.Show()
					if g.Info.GameOver {
						g.reset()
						g.Show()
					}
					log.Println("call updateData()")
					//g.updateInfo()
				}
			}
		}
	}
}

func (g *GameType) onRun() string {
	var strs []string
	_, msg := g.running()
	if len(msg) > 0 {
		strs = append(strs, msg)
	}
	g.Info.Turn++
	//g.Player(g.Info.currentUserID).MakeTurn()
	return strings.Join(strs, "\n")
}

func (g *GameType) running() (r bool, s string) {
	var strs []string
	g.OnPlay()
	thisPlayer := g.Player(g.Info.currentUserID)
	lastMyRank := g.GetRanking(thisPlayer.GetUserID())
	thisPlayer.Property.Move = 0
	if s := g.checkDeBuffAttackFunc(thisPlayer); len(s) > 0 {
		strs = append(strs, s)
	}

	if thisPlayer.Property.Stop == false {
		if s := g.checkBuffFunc(thisPlayer); len(s) > 0 {
			strs = append(strs, s)
		}
	}
	if thisPlayer.Property.Stop == false {
		if s := g.checkDeBuffEffectFunc(thisPlayer); len(s) > 0 {
			strs = append(strs, s)
		}

		a, b, c := thisPlayer.GetDice()
		runDice := fmt.Sprintf("%dD%d+%d", 1+a, 6+b, 0)
		dice.Dice.Roll(runDice)
		rolls := dice.Dice.Rolls
		d := dice.Dice.Hit
		thisPlayer.Property.DiceHit += helper.Max(d+c, 0)
		thisPlayer.Property.Move += thisPlayer.Property.DiceHit

		strs = append(strs, fmt.Sprintf("%s %s%+d(%s%s)", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), thisPlayer.Property.Move, func() string {
			var tmpstrs []string
			for _, val := range rolls {
				tmpstrs = append(tmpstrs, fmt.Sprintf("%s%d", emoji.Emoji(":game_die:"), val))
			}
			return strings.Join(tmpstrs, "")
		}(), func() string {
			return fmt.Sprintf("%+d", c)
		}()))
		if thisPlayer.Property.Move <= 0 {
			thisPlayer.Stop = true
			strs = append(strs, fmt.Sprintf("%s停留原地", thisPlayer.GetDisplayName()))
		}

		if thisPlayer.Property.Stop == false {
			if s := g.checkDeBuffSpeedLimitFunc(thisPlayer); len(s) > 0 {
				strs = append(strs, s)
			}
		}

		if thisPlayer.Property.Stop == false {
			thisPlayer.Property.TotalMove += thisPlayer.Property.Move
			if len(thisPlayer.Buff) > 0 {
				if thisPlayer.Buff[0] == "3" && thisPlayer.GetTurn() != 1 {
					nowMyRank := g.GetRanking(thisPlayer.GetUserID())
					if nowMyRank < lastMyRank {
						rk := g.getRaceSort()
						for i := 0; i < lastMyRank-nowMyRank; i++ {
							rk[i+1].AddDeBuff("speed_down1")
							thisPlayer.AddDeBuff("speed_up2")
							strs = append(strs, fmt.Sprintf("%s「攻城車」撞擊%s", thisPlayer.GetDisplayName(), rk[i+1].GetDisplayName()))
						}
					}
				}
			}
			strs = append(strs, fmt.Sprintf("%s %s%d(%+d)", thisPlayer.GetDisplayName(), emoji.Emoji(":footprints:"), thisPlayer.Property.TotalMove, thisPlayer.Property.Move))
		} else {
			thisPlayer.Property.Move = 0
		}
	}
	thisPlayer.ResetDice()
	r = true
	s = strings.Join(strs, "\n")
	return
}

func (g *GameType) checkBuffFunc(thisPlayer *players.PlayerOption) string {
	var strs []string
	for _, coid := range thisPlayer.Property.Buff {
		if g.mythosCards.Cards[coid].OnEffectFunc != nil {
			if r, s := g.mythosCards.Cards[coid].OnEffectFunc(thisPlayer); r {
				if len(s) > 0 {
					strs = append(strs, s)
				}
			}
		}
	}
	return strings.Join(strs, "\n")
}
func (g *GameType) checkDeBuffAttackFunc(thisPlayer *players.PlayerOption) string {
	var strs []string
	var debuffs []string
	for _, coid := range thisPlayer.Property.DeBuff {
		if _, exist := g.mythosCards.Cards[coid]; exist {
			if g.mythosCards.Cards[coid].OnAttackFunc != nil {
				if r, s := g.mythosCards.Cards[coid].OnAttackFunc(thisPlayer); r {
					if len(s) > 0 {
						strs = append(strs, s)
					}
				}
			} else {
				debuffs = append(debuffs, coid)
			}
		}
	}
	thisPlayer.Property.DeBuff = debuffs
	return strings.Join(strs, "\n")
}

func (g *GameType) checkDeBuffEffectFunc(thisPlayer *players.PlayerOption) string {
	var strs []string
	var debuffs []string
	for _, coid := range thisPlayer.Property.DeBuff {
		if _, exist := g.mythosCards.Cards[coid]; exist {
			if g.mythosCards.Cards[coid].OnEffectFunc != nil {
				if r, s := g.mythosCards.Cards[coid].OnEffectFunc(thisPlayer); r {
					if len(s) > 0 {
						strs = append(strs, s)
					}
				}
			} else {
				debuffs = append(debuffs, coid)
			}
		}
	}
	thisPlayer.Property.DeBuff = debuffs
	return strings.Join(strs, "\n")
}

func (g *GameType) checkDeBuffSpeedLimitFunc(thisPlayer *players.PlayerOption) string {
	var strs []string
	var debuffs []string
	for _, coid := range thisPlayer.Property.DeBuff {
		if _, exist := g.mythosCards.Cards[coid]; exist {
			if g.mythosCards.Cards[coid].OnSpeedLimitFunc != nil || thisPlayer.Stop == false {
				if r, s := g.mythosCards.Cards[coid].OnSpeedLimitFunc(thisPlayer); r {
					if len(s) > 0 {
						strs = append(strs, s)
					}
				}
			} else {
				debuffs = append(debuffs, coid)
			}
		}
	}
	thisPlayer.Property.DeBuff = debuffs
	return strings.Join(strs, "\n")
}

func (g *GameType) newSeason() {
	g.ClearPlayers()
	g.resetInfo()
	g.Show()
}

func (g *GameType) newRound() {
	g.resetInfo()
	g.Show()
}
