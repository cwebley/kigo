package load

import (
	"encoding/json"
	"io/ioutil"
	"log"
	redis "vube/core.go/connect/sexyredis"
	"vube/practice/points/reporter"
)

func LoadFile(file string) error {
	var report reporter.Report
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error loading file %s. err: %s\n", file, err)
		return err
	}
	if err = json.Unmarshal(f, &report); err != nil {
		log.Printf("Error unmarshaling file: %s, err: %s\n", f, err)
		return err
	}

	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("redis.Get failed to return a connection. err: %s", err)
		return err
	}
	defer rConn.Quit()

	for _, player := range report.Players {
		err = player.SetPlayerStats(rConn)
		if err != nil {
			log.Printf("SetPlayerStats for player %s, err: %s\n", player.Name, err)
		}
	}
	return nil

}
