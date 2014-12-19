package saves

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
	redis "vube/core.go/connect/sexyredis"
	"vube/practice/points/constants"
	"vube/practice/points/reporter"
)

const hardCodedPath string = "/usr/local/vube/go/src/vube/practice/points/saves/"

func SaveAsFile() error {
	layout := "2006-01-02"
	fileName := hardCodedPath + time.Now().Format(layout) + ".json"

	var report reporter.Report

	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("redis.Get failed to return a connection. err: %s", err)
		return err
	}
	defer rConn.Quit()

	for _, playerName := range constants.PlayerList {
		var player reporter.Player

		err := player.GetPlayerStats(playerName, rConn)
		if err != nil {
			log.Printf("GetPlayerStats failed for %s", playerName)
			return err
		}

		for i, characterName := range constants.XterList {
			var character reporter.Character
			character.GetCharacterStats(playerName, characterName, rConn)
			player.Characters = append(player.Characters, character)
			for _, opponentName := range constants.XterList {
				var opponent reporter.Opponent
				opponent.GetOpponentStats(playerName, characterName, opponentName, rConn)
				player.Characters[i].Opponents = append(player.Characters[i].Opponents, opponent)
			}
		}
		report.Players = append(report.Players, player)
	}
	reportJson, err := json.Marshal(report)
	if err != nil {
		log.Printf("json Marshal failed for report. err: %s\n", err)
		return err
	}

	err = ioutil.WriteFile(fileName, reportJson, 0777)
	if err != nil {
		log.Printf("failed to writefile %s. err : %s", fileName, err)
		return err
	}
	return nil
}
