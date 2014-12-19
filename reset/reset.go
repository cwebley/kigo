package reset

import (
	"fmt"
	"log"
	redis "vube/core.go/connect/sexyredis"
	"vube/practice/points/constants"
)

func ResetOverallData() error {
	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("Error getting connections. panic time.")
		return err
	}
	defer rConn.Quit()

	for _, player := range constants.PlayerList {
		key := fmt.Sprintf("%s:overall:wins", player)
		_, err = rConn.Set(key, 0)
		if err != nil {
			log.Printf("ERROR resetting key: %s. err: %s", key, err)
			return err
		}
		_, err = rConn.Set(fmt.Sprintf("%s:overall:losses", player), 0)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return err
		}
		_, err = rConn.Set(fmt.Sprintf("%s:overall:streak", player), 0)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return err
		}
		for _, character := range constants.XterList {
			_, err = rConn.Set(fmt.Sprintf("%s:%s:overall:wins", player, character), 0)
			if err != nil {
				log.Printf("ERROR: %s", err)
				return err
			}
			_, err = rConn.Set(fmt.Sprintf("%s:%s:overall:losses", player, character), 0)
			if err != nil {
				log.Printf("ERROR: %s", err)
				return err
			}
			_, err = rConn.Set(fmt.Sprintf("%s:%s:overall:streak", player, character), 0)
			if err != nil {
				log.Printf("ERROR: %s", err)
				return err
			}
			for _, opponent := range constants.XterList {
				_, err = rConn.Set(fmt.Sprintf("%s:%s:%s:overall:wins", player, character, opponent), 0)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return err
				}
				_, err = rConn.Set(fmt.Sprintf("%s:%s:%s:overall:losses", player, character, opponent), 0)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return err
				}
				_, err = rConn.Set(fmt.Sprintf("%s:%s:%s:overall:streak", player, character, opponent), 0)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return err
				}
			}
		}
	}
	return nil
}
