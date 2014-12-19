package reporter

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	redis "vube/core.go/connect/sexyredis"
	"vube/core.go/connect/sexyredis/supersexyclient"
	"vube/practice/points/constants"
)

type Report struct {
	Players []Player `json:"players"`
}
type Player struct {
	Name          string      `json:"name"`
	Score         int         `json:"score"`
	Characters    []Character `json:"characters"`
	Record        Ratio       `json:"current_record"`
	OverallRecord Ratio       `json:"overall_record"`
	Streak        int         `json:"streak"`
	OverallStreak int         `json:"overall_streak"`
}
type Character struct {
	Name          string     `json:"name"`
	Value         int        `json:"value"`
	Record        Ratio      `json:"current_record"`
	OverallRecord Ratio      `json:"overall_record"`
	Streak        int        `json:"streak"`
	OverallStreak int        `json:"overall_streak"`
	Opponents     []Opponent `json:"opponents"`
}
type Opponent struct {
	Name          string `json:"name"`
	Record        Ratio  `json:"current_record"`
	OverallRecord Ratio  `json:"overall_record"`
	Streak        int    `json:"streak"`
	OverallStreak int    `json:"overall_streak"`
	IsRival       bool   `json:"is_rival"`
}
type Ratio struct {
	Wins   int `json:"wins"`
	Losses int `json:"losses"`
}

func NewReport() (string, error) {
	var report Report
	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("redis.Get failed to return a connection. err: %s", err)
		return "", err
	}
	defer rConn.Quit()

	for _, playerName := range constants.PlayerList {
		// var player Player
		player, err := GetPlayerAllStats(playerName, rConn)
		if err != nil {
			log.Printf("GetPlayerAllStats")
			return "", err
		}
		report.Players = append(report.Players, player)
	}
	reportJson, err := json.Marshal(report)
	if err != nil {
		log.Printf("json Marshal failed for report. err: %s\n", err)
		return "", nil
	}
	return string(reportJson), nil
}
func GetPlayerAllStats(playerName string, rConn supersexyclient.SuperSexyClient) (Player, error) {
	var p Player
	if err := p.GetPlayerStats(playerName, rConn); err != nil {
		log.Printf("Error retrieving stats. err: %s", err)
		return p, err
	}
	for i, characterName := range constants.XterList {
		var c Character
		if err := c.GetCharacterStats(playerName, characterName, rConn); err != nil {
			log.Printf("Error retrieving stats. err: %s", err)
			return p, err
		}
		p.Characters = append(p.Characters, c)
		for _, opponentName := range constants.XterList {
			var o Opponent
			if err := o.GetOpponentStats(playerName, characterName, opponentName, rConn); err != nil {
				log.Printf("Error retrieving stats. err: %s", err)
				return p, err
			}
			p.Characters[i].Opponents = append(p.Characters[i].Opponents, o)
		}
	}
	return p, nil
}

func (p *Player) GetPlayerStats(playerName string, rConn supersexyclient.SuperSexyClient) error {
	score, err := GetScore(rConn, MakePlayerPtsKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	curWins, err := GetScore(rConn, MakePlayerCurWinsKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	curLosses, err := GetScore(rConn, MakePlayerCurLossesKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	overallWins, err := GetScore(rConn, MakePlayerOverallWinsKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	overallLosses, err := GetScore(rConn, MakePlayerOverallLossesKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	curStreak, err := GetScore(rConn, MakePlayerCurStreakKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	overallStreak, err := GetScore(rConn, MakePlayerOverallStreakKey(playerName))
	if err != nil {
		log.Printf("error getting player score for %s", playerName)
		return err
	}
	p.Name = playerName
	p.Score = score
	p.Record.Wins = curWins
	p.Record.Losses = curLosses
	p.OverallRecord.Wins = overallWins
	p.OverallRecord.Losses = overallLosses
	p.Streak = curStreak
	p.OverallStreak = overallStreak
	return nil
}

func (c *Character) GetCharacterStats(playerName, characterName string, rConn supersexyclient.SuperSexyClient) error {
	value, err := GetScore(rConn, MakeCharacterValueKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's value", playerName, characterName)
		return err
	}
	curWins, err := GetScore(rConn, MakeCharacterCurWinsKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	curLosses, err := GetScore(rConn, MakeCharacterCurLossesKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	overallWins, err := GetScore(rConn, MakeCharacterOverallWinsKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	overallLosses, err := GetScore(rConn, MakeCharacterOverallLossesKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	curStreak, err := GetScore(rConn, MakeCharacterCurStreakKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	overallStreak, err := GetScore(rConn, MakeCharacterOverallStreakKey(playerName, characterName))
	if err != nil {
		log.Printf("error getting %s's %s's score", playerName, characterName)
		return err
	}
	c.Name = characterName
	c.Value = value
	c.Record.Wins = curWins
	c.Record.Losses = curLosses
	c.OverallRecord.Wins = overallWins
	c.OverallRecord.Losses = overallLosses
	c.Streak = curStreak
	c.OverallStreak = overallStreak
	return nil
}

func (o *Opponent) GetOpponentStats(playerName, characterName, opponentName string, rConn supersexyclient.SuperSexyClient) error {
	curWins, err := GetScore(rConn, MakeOpponentCurWinsKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	curLosses, err := GetScore(rConn, MakeOpponentCurLossesKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	overallWins, err := GetScore(rConn, MakeOpponentOverallWinsKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	overallLosses, err := GetScore(rConn, MakeOpponentOverallLossesKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	curStreak, err := GetScore(rConn, MakeOpponentCurStreakKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	overallStreak, err := GetScore(rConn, MakeOpponentOverallStreakKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	isRival := false
	isRivalInt, err := GetScore(rConn, MakeOpponentRivalKey(playerName, characterName, opponentName))
	if err != nil {
		log.Printf("error getting %s's %s's opponent %s score", playerName, characterName, opponentName)
		return err
	}
	if isRivalInt == 1 {
		isRival = true
	}
	o.Name = opponentName
	o.Record.Wins = curWins
	o.Record.Losses = curLosses
	o.OverallRecord.Wins = overallWins
	o.OverallRecord.Losses = overallLosses
	o.Streak = curStreak
	o.OverallStreak = overallStreak
	o.IsRival = isRival
	return nil
}

func GetScore(rConn supersexyclient.SuperSexyClient, key string) (int, error) {
	playerPtsStr, err := rConn.Get(key)
	if err != nil {
		log.Printf("rConn.Get failed %s. err: %s", key, err)
		return 0, err
	}
	if playerPtsStr == "" {
		err := fmt.Errorf("redis value not intialized for key: %s", key)
		log.Printf("%s", err)
		return 0, err
	}
	playerPts, err := strconv.Atoi(playerPtsStr)
	if err != nil {
		log.Printf("strconv failed %s, err %s", playerPtsStr, err)
		return 0, err
	}
	return playerPts, err
}

func (p *Player) SetPlayerStats(rConn supersexyclient.SuperSexyClient) error {
	err := SetScore(rConn, MakePlayerPtsKey(p.Name), p.Score)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerCurWinsKey(p.Name), p.Record.Wins)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerCurLossesKey(p.Name), p.Record.Losses)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerOverallWinsKey(p.Name), p.OverallRecord.Wins)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerOverallLossesKey(p.Name), p.OverallRecord.Losses)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerCurStreakKey(p.Name), p.Streak)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	err = SetScore(rConn, MakePlayerOverallStreakKey(p.Name), p.OverallStreak)
	if err != nil {
		log.Printf("error setting value for %s. err: %s\n", p.Name, err)
		return err
	}
	for _, character := range p.Characters {
		err = SetScore(rConn, MakeCharacterValueKey(p.Name, character.Name), character.Value)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterCurWinsKey(p.Name, character.Name), character.Record.Wins)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterCurLossesKey(p.Name, character.Name), character.Record.Losses)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterOverallWinsKey(p.Name, character.Name), character.OverallRecord.Wins)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterOverallLossesKey(p.Name, character.Name), character.OverallRecord.Losses)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterCurStreakKey(p.Name, character.Name), character.Streak)
		fmt.Println("Character Current STREAK: ", character.Streak, character.Name)

		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		err = SetScore(rConn, MakeCharacterOverallStreakKey(p.Name, character.Name), character.OverallStreak)
		fmt.Println("Character OVERALL STREAK: ", character.OverallStreak, character.Name)
		if err != nil {
			log.Printf("error setting value for %s's %s. err: %s\n", p.Name, character.Name, err)
			return err
		}
		for _, opponent := range character.Opponents {
			err = SetScore(rConn, MakeOpponentCurWinsKey(p.Name, character.Name, opponent.Name), opponent.Record.Wins)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			err = SetScore(rConn, MakeOpponentCurLossesKey(p.Name, character.Name, opponent.Name), opponent.Record.Losses)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			err = SetScore(rConn, MakeOpponentOverallWinsKey(p.Name, character.Name, opponent.Name), opponent.OverallRecord.Wins)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			err = SetScore(rConn, MakeOpponentOverallLossesKey(p.Name, character.Name, opponent.Name), opponent.OverallRecord.Losses)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			err = SetScore(rConn, MakeOpponentCurStreakKey(p.Name, character.Name, opponent.Name), opponent.Streak)

			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			err = SetScore(rConn, MakeOpponentOverallStreakKey(p.Name, character.Name, opponent.Name), opponent.OverallStreak)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
			isRivalInt := 0
			if opponent.IsRival == true {
				isRivalInt = 1
			}
			err = SetScore(rConn, MakeOpponentRivalKey(p.Name, character.Name, opponent.Name), isRivalInt)
			if err != nil {
				log.Printf("error setting value for %s's %s's opponent %s. err: %s\n", p.Name, character.Name, opponent.Name, err)
				return err
			}
		}
	}
	return nil
}

func SetScore(rConn supersexyclient.SuperSexyClient, key string, score int) error {
	_, err := rConn.Set(key, score)
	if err != nil {
		log.Printf("rConn.Get failed %s. err: %s", key, err)
	}

	return err
}

func MakePlayerPtsKey(player string) string {
	return fmt.Sprintf("%s:points", player)
}
func MakeCharacterValueKey(player, character string) string {
	return fmt.Sprintf("%s:%s:value", player, character)
}
func MakePlayerCurWinsKey(player string) string {
	return fmt.Sprintf("%s:current:wins", player)
}
func MakePlayerCurLossesKey(player string) string {
	return fmt.Sprintf("%s:current:losses", player)
}
func MakePlayerOverallWinsKey(player string) string {
	return fmt.Sprintf("%s:overall:wins", player)
}
func MakePlayerOverallLossesKey(player string) string {
	return fmt.Sprintf("%s:overall:losses", player)
}
func MakePlayerCurStreakKey(player string) string {
	return fmt.Sprintf("%s:current:streak", player)
}
func MakePlayerOverallStreakKey(player string) string {
	return fmt.Sprintf("%s:overall:streak", player)
}
func MakeCharacterCurWinsKey(player, character string) string {
	return fmt.Sprintf("%s:%s:current:wins", player, character)
}
func MakeCharacterCurLossesKey(player, character string) string {
	return fmt.Sprintf("%s:%s:current:losses", player, character)
}
func MakeCharacterOverallWinsKey(player, character string) string {
	return fmt.Sprintf("%s:%s:overall:wins", player, character)
}
func MakeCharacterOverallLossesKey(player, character string) string {
	return fmt.Sprintf("%s:%s:overall:losses", player, character)
}
func MakeCharacterCurStreakKey(player, character string) string {
	return fmt.Sprintf("%s:%s:current:streak", player, character)
}
func MakeCharacterOverallStreakKey(player, character string) string {
	return fmt.Sprintf("%s:%s:overall:streak", player, character)
}

//opponents

func MakeOpponentCurWinsKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:current:wins", player, character, opponent)
}
func MakeOpponentCurLossesKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:current:losses", player, character, opponent)
}
func MakeOpponentOverallWinsKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:overall:wins", player, character, opponent)
}
func MakeOpponentOverallLossesKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:overall:losses", player, character, opponent)
}
func MakeOpponentCurStreakKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:current:streak", player, character, opponent)
}
func MakeOpponentOverallStreakKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:overall:streak", player, character, opponent)
}
func MakeOpponentRivalKey(player, character, opponent string) string {
	return fmt.Sprintf("%s:%s:%s:rival", player, character, opponent)
}
