package update

import (
	"encoding/json"
	"time"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	redis "vube/core.go/connect/sexyredis"
	"vube/core.go/connect/sexyredis/supersexyclient"
	"vube/practice/points/constants"
	"vube/practice/points/reporter"
	"math/rand"
)

type Result struct {
	Win  Winner `json:"winner"`
	Loss Loser  `json:"loser"`
}

type Winner struct {
	Player string `json:"player"`
	Xter   string `json:"xter"`
}

type Loser struct {
	Player string `json:"player"`
	Xter   string `json:"xter"`
}

var BjUpcoming = make([]string,0)
var GUpcoming = make([]string,0)

func UnmarshalResult(r *http.Request) (result Result, err error) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error-reading-game-results. err: %s\n", err)
		return
	}
	fmt.Println("unmarsh", string(b))
	if err := json.Unmarshal(b, &result); err != nil {
		log.Printf("Unmarshaling-error. err: %s\n", err)
	}
	legit := verifyResult(result)
	if !legit {
		err = fmt.Errorf("Typo-in-result. Result-not-legit")
		log.Printf("Error-reading-game-results. err: %s\n", err)
		return
	}
	return
}

func UpdateAndRecordResults(result Result) string {
	success := `{"result": "ok"}`
	failure := `{"result": "fail"}`

	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("redis.Get failed to return a connection")
		return failure
	}
	defer rConn.Quit()

	winnerName := result.Win.Player
	winnerXterName := result.Win.Xter
	loserName := result.Loss.Player
	loserXterName := result.Loss.Xter

	var winner reporter.Player
	var wxster reporter.Character
	var woppo reporter.Opponent
	var loser reporter.Player
	var lxster reporter.Character
	var loppo reporter.Opponent

	err = winner.GetPlayerStats(winnerName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s err: %s\n", winnerName, err)
		return failure
	}
	err = wxster.GetCharacterStats(winnerName, winnerXterName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s's %s err: %s\n", winnerName, winnerXterName, err)
		return failure
	}
	err = woppo.GetOpponentStats(winnerName, winnerXterName, loserXterName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s's %s's opponent %s err: %s\n", winnerName, winnerXterName, loserXterName, err)
		return failure
	}

	err = loser.GetPlayerStats(loserName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s err: %s\n", loserName, err)
		return failure
	}
	err = lxster.GetCharacterStats(loserName, loserXterName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s's %s err: %s\n", loserName, loserXterName, err)
		return failure
	}
	err = loppo.GetOpponentStats(loserName, loserXterName, winnerXterName, rConn)
	if err != nil {
		log.Printf("Failed to get stats for %s's %s's opponent %s err: %s\n", winnerName, winnerXterName, loserXterName, err)
		return failure
	}
	winner.Characters = append(winner.Characters, wxster)
	winner.Characters[0].Opponents = append(winner.Characters[0].Opponents, woppo)
	loser.Characters = append(loser.Characters, lxster)
	loser.Characters[0].Opponents = append(loser.Characters[0].Opponents, loppo)

	winner, loser = ProcessResults(winner, loser)
	err = UpdateRedis(winner, loser, rConn)
	if err != nil {
		log.Printf("failed to update redis with game results. err: %s\n")
		return failure
	}

	NextFight(rConn)

	return success
}

// Only function with bj/g right now, comment out if necessary
// rivalry recalculations not done here.
func NextFight(rConn supersexyclient.SuperSexyClient){
	log.Printf("NEXT FIGHT STUFF")
	if len(GUpcoming) == 3{
		GUpcoming = GUpcoming[1:3]
	}
	if len(BjUpcoming) == 3{
		BjUpcoming = BjUpcoming[1:3]
	}
	// log.Printf("PREUPCOMING: %v, %v", GUpcoming, BjUpcoming)


	//Build upcoming fighter array
	for len(GUpcoming)<3 {
		rand.Seed(time.Now().UnixNano())
		GUpcoming = append(GUpcoming, constants.XterList[rand.Intn(len(constants.XterList))])
		BjUpcoming = append(BjUpcoming, constants.XterList[rand.Intn(len(constants.XterList))])
	}

	log.Printf("UPCOMING MATCH: G: %v, BJ: %v", GUpcoming[0], BjUpcoming[0])

	//get matchup value
	var gPlayer reporter.Player
	var g1X reporter.Character
	var g1O reporter.Opponent
	var g2X reporter.Character
	var g2O reporter.Opponent
	var g3X reporter.Character
	var g3O reporter.Opponent

	var bjPlayer reporter.Player
	var bj1X reporter.Character
	var bj1O reporter.Opponent
	var bj2X reporter.Character
	var bj2O reporter.Opponent
	var bj3X reporter.Character
	var bj3O reporter.Opponent

	gPlayer.GetPlayerStats("g", rConn)
	g1X.GetCharacterStats("g", GUpcoming[0], rConn)
	g1O.GetOpponentStats("g", GUpcoming[0], BjUpcoming[0], rConn)
	g2X.GetCharacterStats("g", GUpcoming[1], rConn)
	g2O.GetOpponentStats("g", GUpcoming[1], BjUpcoming[1], rConn)
	g3X.GetCharacterStats("g", GUpcoming[2], rConn)
	g3O.GetOpponentStats("g", GUpcoming[2], BjUpcoming[2], rConn)

	bjPlayer.GetPlayerStats("bj", rConn)
	bj1X.GetCharacterStats("bj", GUpcoming[0], rConn)
	bj1O.GetOpponentStats("bj", GUpcoming[0], BjUpcoming[0], rConn)
	bj2X.GetCharacterStats("bj", GUpcoming[1], rConn)
	bj2O.GetOpponentStats("bj", GUpcoming[1], BjUpcoming[1], rConn)
	bj3X.GetCharacterStats("bj", GUpcoming[2], rConn)
	bj3O.GetOpponentStats("bj", GUpcoming[2], BjUpcoming[2], rConn)

	var g1Round int
	var g2Round int
	var g3Round int
	var bj1Round int
	var bj2Round int
	var bj3Round int


	// repeats
	if g1X.Name == g2X.Name{
		g2X.Value--
	}
	if g1X.Name == g3X.Name{
		g3X.Value--
	}
	if g2X.Name == g3X.Name{
		g3X.Value--
	}

	//fire
	if g1X.Streak == 2{
		if g1X.Name != g2X.Name{
			g2X.Value++
		}
		if g1X.Name != g3X.Name{
			g3X.Value++
		}
	}
	if g2X.Streak == 2{
		if g2X.Name != g3X.Name{
			g3X.Value++
		}
	}
	if g1X.Streak == 1 && g1X.Name == g2X.Name{
		g3X.Value++
	}

	//rivals
	if g1O.IsRival{
		g1X.Value += 2
	}
	if g2O.IsRival{
		g2X.Value += 2
	}
	if g3O.IsRival{
		g3X.Value += 2
	}

	//repeat
	if bj1X.Name == bj2X.Name{
		bj2X.Value--
	}
	if bj1X.Name == bj3X.Name{
		bj3X.Value--
	}
	if bj2X.Name == bj3X.Name{
		bj3X.Value--
	}

	//fire
	if bj1X.Streak == 2{
		if bj1X.Name != bj2X.Name{
			bj2X.Value++
		}
		if bj1X.Name != bj3X.Name{
			bj3X.Value++
		}
	}
	if bj2X.Streak == 2{
		if bj2X.Name != bj3X.Name{
			bj3X.Value++
		}
	}
	if bj1X.Streak == 1 && bj1X.Name == bj2X.Name{
		bj3X.Value++
	}

	g1Round = gPlayer.Score + g1X.Value
	g2Round = gPlayer.Score + g1X.Value + g2X.Value
	g3Round = gPlayer.Score + g1X.Value + g2X.Value + g3X.Value

	bj1Round = bjPlayer.Score + bj1X.Value
	bj2Round = bjPlayer.Score + bj1X.Value + bj2X.Value
	bj3Round = bjPlayer.Score + bj1X.Value + bj2X.Value + bj3X.Value

	switch {
	case gPlayer.Score >= constants.WinningScore:
		log.Printf("GAME OVER!!! g WINS!\n")
	case g1Round >= constants.WinningScore:
		log.Printf("G CAN WIN NEXT GAME!\n")
	case g2Round >= constants.WinningScore:
		log.Printf("G CAN WIN IN 2 GAMES!\n") 
	case g3Round >= constants.WinningScore:
		log.Printf("G CAN WIN IN 3 GAMES!\n") 
	}

	switch {
	case bjPlayer.Score >= constants.WinningScore:
		log.Printf("GAME OVER!!! bj WINS!\n")
	case bj1Round >= constants.WinningScore:
		log.Printf("BJ CAN WIN NEXT GAME!\n")
	case bj2Round >= constants.WinningScore:
		log.Printf("BJ CAN WIN IN 2 GAMES!\n") 
	case bj3Round >= constants.WinningScore:
		log.Printf("BJ CAN WIN IN 3 GAMES!\n") 
	}

	switch{
	case gPlayer.Score >= constants.WinningScore2:
		log.Printf("GAME OVER BIG!!! g WINSBIG!\n")
	case g1Round >= constants.WinningScore2:
		log.Printf("G CAN WIN BIG NEXT GAME!\n")
	case g2Round >= constants.WinningScore2:
		log.Printf("G CAN WIN BIG IN 2 GAMES!\n") 
	case g3Round >= constants.WinningScore2:
		log.Printf("G CAN WIN BIG IN 3 GAMES!\n") 
	}

	switch{
	case bjPlayer.Score >= constants.WinningScore2:
		log.Printf("GAME OVER BIG!!! g WINS!\n")
	case bj1Round >= constants.WinningScore2:
		log.Printf("BJ CAN WIN BIG NEXT GAME!\n")
	case bj2Round >= constants.WinningScore2:
		log.Printf("BJ CAN WIN BIG IN 2 GAMES!\n") 
	case bj3Round >= constants.WinningScore2:
		log.Printf("BJ CAN WIN BIG IN 3 GAMES!\n") 
	}

}

func UpdateRedis(winner, loser reporter.Player, rConn supersexyclient.SuperSexyClient) error {
	err := winner.SetPlayerStats(rConn)
	if err != nil {
		log.Printf("failed to set player stats for %s. err: %s\n", winner, err)
		return err
	}
	err = loser.SetPlayerStats(rConn)
	if err != nil {
		log.Printf("failed to set player stats for %s. err: %s\n", loser, err)
	}
	log.Printf("Previous Winner: %s. Total %d. Previous Loser: %s. Total %d", winner.Name, winner.Score, loser.Name, loser.Score)
	return err
}

func ProcessResults(winner, loser reporter.Player) (reporter.Player, reporter.Player) {
	// winning rival match worth 2. basically you both go up a point after the decrement.
	if winner.Characters[0].Opponents[0].IsRival == true {
		log.Printf("RIVLARY MATCH! %s's %s is going up 2 points before stat logging", winner.Name, winner.Characters[0].Name)
		winner.Characters[0].Value++
	}

	log.Printf("%s's %s was worth %d points", winner.Name, winner.Characters[0].Name, winner.Characters[0].Value)
	winner.Record.Wins++
	winner.OverallRecord.Wins++
	winner.Streak++
	if winner.Streak > winner.OverallStreak {
		log.Printf("NEW HIGH STREAK OF %d FOR %s", winner.Streak, winner.Name)
		winner.OverallStreak = winner.Streak
	}
	winner.Score += winner.Characters[0].Value
	if winner.Characters[0].Value != 1 {
		winner.Characters[0].Value--
	}
	if winner.Characters[0].Opponents[0].IsRival == true {
		winner.Characters[0].Value++
	}

	winner.Characters[0].Record.Wins++
	winner.Characters[0].OverallRecord.Wins++
	winner.Characters[0].Streak++
	if winner.Characters[0].Streak > winner.Characters[0].OverallStreak {
		log.Printf("NEW HIGH STREAK OF %d FOR %s", winner.Characters[0].Streak, winner.Characters[0].Name)
		winner.Characters[0].OverallStreak = winner.Characters[0].Streak
	}
	if winner.Characters[0].Streak == 2 {
		log.Printf("%s's %s IS HEATING UP.", winner.Name, winner.Characters[0].Name)
	}

	if winner.Characters[0].Streak == 3 {
		log.Printf("%s's %s IS ON FIRE. +1 value to all characters", winner.Name, winner.Characters[0].Name)
		err := OnFire(winner.Name)
		if err != nil {
			log.Printf("OnFire failed for %s", winner.Name)
		}
	}
	winner.Characters[0].Opponents[0].Record.Wins++
	winner.Characters[0].Opponents[0].OverallRecord.Wins++
	winner.Characters[0].Opponents[0].Streak++
	if winner.Characters[0].Opponents[0].Streak > winner.Characters[0].Opponents[0].OverallStreak {
		winner.Characters[0].OverallStreak = winner.Characters[0].Streak
	}
	// CREATE/ DESTROY RIVALRIES
	if winner.Characters[0].Opponents[0].Record.Wins == winner.Characters[0].Opponents[0].Record.Losses {
		// set new rivalry for next match
		log.Printf("NEW RIVALRY: %s's %s and %s' %s\n", winner.Name, winner.Characters[0].Name, loser.Name, loser.Characters[0].Name)
		winner.Characters[0].Opponents[0].IsRival = true
		loser.Characters[0].Opponents[0].IsRival = true
	}

	if winner.Characters[0].Opponents[0].IsRival == true {
		if winner.Characters[0].Opponents[0].Record.Wins-winner.Characters[0].Opponents[0].Record.Losses == 2 {
			// set new rivalry for next match
			log.Printf("RIVALRY ICED!: %s's %s and %s' %s. nothing happens though.\n", winner.Name, winner.Characters[0].Name, loser.Name, loser.Characters[0].Name)
			winner.Characters[0].Opponents[0].IsRival = false
			loser.Characters[0].Opponents[0].IsRival = false
		}
	}

	loser.Record.Losses++
	loser.OverallRecord.Losses++
	loser.Streak = 0

	loser.Characters[0].Value++
	loser.Characters[0].Record.Losses++
	loser.Characters[0].OverallRecord.Losses++
	if loser.Characters[0].Streak >= 3 {
		log.Printf("%s's %s has been ICED! -1 to all characters", loser.Name, loser.Characters[0].Name)
		err := IceDown(loser.Name)
		if err != nil {
			log.Printf("IceDown failed for %s", winner.Name)
		}
	}
	loser.Characters[0].Streak = 0

	loser.Characters[0].Opponents[0].Record.Losses++
	loser.Characters[0].Opponents[0].OverallRecord.Losses++
	loser.Characters[0].Opponents[0].Streak = 0

	return winner, loser
}

func IceDown(playerName string) error {
	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("Error getting connections. panic time.")
		return err
	}
	defer rConn.Quit()

	for _, character := range constants.XterList {
		valStr, err := rConn.Get(fmt.Sprintf("%s:%s:value", playerName, character))
		if err != nil {
			return fmt.Errorf("failed to IceDown get for %s's %s", playerName, character)
		}
		valInt, err := strconv.Atoi(valStr)
		if err != nil {
			log.Printf("strconv failed %s, err %s", valStr, err)
			return err
		}
		valIntNew := valInt - 1


		_, err = rConn.Set(fmt.Sprintf("%s:%s:value", playerName, character), valIntNew)
		if err != nil {
			return fmt.Errorf("failed to IceDown get for %s's %s", playerName, character)
		}
	}
	return nil
}

func OnFire(playerName string) error {
	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("Error getting connections. panic time.")
		return err
	}
	defer rConn.Quit()

	for _, character := range constants.XterList {
		_, err = rConn.Incr(fmt.Sprintf("%s:%s:value", playerName, character))
		if err != nil {
			return fmt.Errorf("failed to OnFire incr for %s's %s", playerName, character)
		}
	}
	return nil
}

func verifyResult(result Result) (legit bool) {
	var (
		xter, player string
	)
	var xterList = constants.XterList
	var playerList = constants.PlayerList

	for _, xter = range xterList {
		if result.Win.Xter == xter {
			legit = true
			break
		}
	}
	if !legit {
		return
	}
	legit = false
	for _, player = range playerList {
		if result.Win.Player == player {
			legit = true
			break
		}
	}
	if !legit {
		return
	}
	legit = false
	for _, xter = range xterList {
		if result.Loss.Xter == xter {
			legit = true
			break
		}
	}
	if !legit {
		return
	}
	legit = false
	for _, player = range playerList {
		if result.Loss.Player == player {
			legit = true
			break
		}
	}
	if !legit {
		return
	}
	return
}
