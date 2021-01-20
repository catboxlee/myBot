package players

import (
	"log"
	"myBot/games/zombile/cards"
	"myBot/games/zombile/power"
	"strconv"
	"strings"
)

// Players ...
type Players struct {
	Player map[string]*PlayerOption
}

const (
	defaultActionTimes int = 2
	defaultAttackTimes int = 1
	defaultHealth      int = 6
	defaultHorror      int = 9
	defaultCombat      int = 1
)

// PlayerJoin ...
func (p *Players) PlayerJoin(playerID string, playerName string) (bool, string) {
	p.Player[playerID] = p.newPlayer(playerID, playerName)
	return true, playerName
}

// newPlayer ...
func (p *Players) newPlayer(userid string, displayname string) *PlayerOption {
	return &PlayerOption{userid, displayname, defaultActionTimes, defaultAttackTimes, defaultHealth, defaultHealth, 0, defaultHorror, defaultCombat, 1, cards.Cards{}}
}

// ClearPlayers ...
func (p *Players) ClearPlayers() {
	p.Player = make(map[string]*PlayerOption)
}

// GetPlayer ...
func (p *Players) GetPlayer(playerID string) power.PlayerIF {
	if p.Player[playerID] != nil {
		return p.Player[playerID]
	}
	return nil
}

// ActivateCard ...
func (p *Players) ActivateCard(currentPlayerID string, cardIndex int, command string) string {
	log.Println("Players.activateCard()")

	currentPlayer := p.GetPlayer(currentPlayerID).(*PlayerOption)

	if currentCard := currentPlayer.GetCard(cardIndex); currentCard == nil {
		return ""
	}

	cmds := strings.Split(command, " ")
	var targetPlayer power.PlayerIF = nil
	var target power.FightIF = nil

	if len(cmds) >= 1 {
		if cmds[0][0:1] == "@" {
			cmds[0] = cmds[0][1:]
		}
		if v := p.GetPlayer(cmds[0]); v != nil {
			targetPlayer = v
			target = v
			if len(cmds) > 1 {
				if x, err := strconv.Atoi(cmds[1]); err == nil {
					if c := targetPlayer.GetCard(x); c != nil {
						target = c
					}
				}
			}
		} else {
			if x, err := strconv.Atoi(cmds[0]); err == nil {
				if c := currentPlayer.GetCard(x); c != nil {
					target = c
				}
			}
		}
	}
	return currentPlayer.Card[cardIndex].OnPlay(currentPlayer, targetPlayer, target)
}
