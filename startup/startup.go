package startup

import (
	"fmt"
	"log"
	redis "vube/core.go/connect/sexyredis"
	"vube/practice/points/constants"
)

type InitVal int

const (
	G_FULGORE InitVal = 8
	G_SPINAL          = 1
	G_WULF            = 7
	G_ORCHID          = 3
	G_GLACIUS         = 2
	G_JAGO            = 10
	G_THUNDER         = 9
	G_SADIRA          = 5
	G_TJ 				= 12
	G_MAYA 				= 4
	G_KANRA				= 6
	G_RIPTOR			= 11

	BJ_FULGORE InitVal = 10
	BJ_SPINAL          = 12
	BJ_WULF            = 1
	BJ_ORCHID          = 2
	BJ_GLACIUS         = 7
	BJ_JAGO            = 3
	BJ_THUNDER         = 4
	BJ_SADIRA          = 5
	BJ_TJ 				= 6
	BJ_MAYA 			= 9
	BJ_KANRA			= 11
	BJ_RIPTOR			= 8
)

func LoadXterVals() {
	rConn, err := redis.Get("rw")
	if err != nil {
		log.Printf("redis.Get failed to return a connection")
		panic(err)
	}
	defer rConn.Quit()

	bjVals := map[string]InitVal{"fulgore": BJ_FULGORE, "spinal": BJ_SPINAL, "wulf": BJ_WULF, "orchid": BJ_ORCHID, "glacius": BJ_GLACIUS, "jago": BJ_JAGO, "thunder": BJ_THUNDER, "sadira": BJ_SADIRA, "tj": BJ_TJ, "maya": BJ_MAYA, "kanra": BJ_KANRA, "riptor": BJ_RIPTOR}
	gVals := map[string]InitVal{"fulgore": G_FULGORE, "spinal": G_SPINAL, "wulf": G_WULF, "orchid": G_ORCHID, "glacius": G_GLACIUS, "jago": G_JAGO, "thunder": G_THUNDER, "sadira": G_SADIRA, "tj": G_TJ, "maya": G_MAYA, "kanra": G_KANRA, "riptor": G_RIPTOR}

	for bjkey, bjval := range bjVals {
		_, err := rConn.Set("bj:"+bjkey+":value", bjval)
		if err != nil {
			panic(err)
		}
	}
	for gkey, gval := range gVals {
		_, err := rConn.Set("g:"+gkey+":value", gval)
		if err != nil {
			panic(err)
		}
	}

	res := ""
	res, err = rConn.Get("bj:points")
	if err != nil {
		panic(err)
	}
	log.Printf("bj's previous game score: %s\n", res)

	res, err = rConn.Get("g:points")
	if err != nil {
		panic(err)
	}
	log.Printf("g's previous game score: %s\n", res)

	// res, err = rConn.Get("mike:points")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("mike's previous game score: %s\n", res)

	// bj
	_, err = rConn.Set("bj:points", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("bj:current:wins", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("bj:current:losses", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("bj:current:streak", 0)
	if err != nil {
		panic(err)
	}
	for _, char := range constants.XterList {
		_, err = rConn.Set(fmt.Sprintf("bj:%s:current:wins", char), 0)
		if err != nil {
			panic(err)
		}
		_, err = rConn.Set(fmt.Sprintf("bj:%s:current:losses", char), 0)
		if err != nil {
			panic(err)
		}
		_, err = rConn.Set(fmt.Sprintf("bj:%s:current:streak", char), 0)
		if err != nil {
			panic(err)
		}

		// // OVERALL CHARACTER RESETS, COMMENT OUT FOR REPEATS
		// _, err = rConn.Set(fmt.Sprintf("bj:%s:overall:wins", char), 0)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = rConn.Set(fmt.Sprintf("bj:%s:overall:losses", char), 0)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = rConn.Set(fmt.Sprintf("bj:%s:overall:streak", char), 0)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = rConn.Set(fmt.Sprintf("g:%s:overall:wins", char), 0)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = rConn.Set(fmt.Sprintf("g:%s:overall:losses", char), 0)
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = rConn.Set(fmt.Sprintf("g:%s:overall:streak", char), 0)
		// if err != nil {
		// 	panic(err)
		// }


		for _, oppo := range constants.XterList {
			_, err = rConn.Set(fmt.Sprintf("bj:%s:%s:current:wins", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			_, err = rConn.Set(fmt.Sprintf("bj:%s:%s:current:losses", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			_, err = rConn.Set(fmt.Sprintf("bj:%s:%s:current:streak", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			isRivalInt := 0
			if constants.BJRivals[char] == oppo {
				isRivalInt = 1
			}
			_, err = rConn.Set(fmt.Sprintf("bj:%s:%s:rival", char, oppo), isRivalInt)
			if err != nil {
				panic(err)
			}

			//COMMENT OUT FOR REPEATS
			// _, err = rConn.Set(fmt.Sprintf("bj:%s:%s:overall:wins", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = rConn.Set(fmt.Sprintf("bj:%s:%s:overall:losses", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = rConn.Set(fmt.Sprintf("bj:%s:%s:overall:streak", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
		}
	}

	// g
	_, err = rConn.Set("g:points", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("g:current:wins", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("g:current:losses", 0)
	if err != nil {
		panic(err)
	}
	_, err = rConn.Set("g:current:streak", 0)
	if err != nil {
		panic(err)
	}
	for _, char := range constants.XterList {
		_, err = rConn.Set(fmt.Sprintf("g:%s:current:wins", char), 0)
		if err != nil {
			panic(err)
		}
		_, err = rConn.Set(fmt.Sprintf("g:%s:current:losses", char), 0)
		if err != nil {
			panic(err)
		}
		_, err = rConn.Set(fmt.Sprintf("g:%s:current:streak", char), 0)
		if err != nil {
			panic(err)
		}
		for _, oppo := range constants.XterList {
			_, err = rConn.Set(fmt.Sprintf("g:%s:%s:current:wins", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			_, err = rConn.Set(fmt.Sprintf("g:%s:%s:current:losses", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			_, err = rConn.Set(fmt.Sprintf("g:%s:%s:current:streak", char, oppo), 0)
			if err != nil {
				panic(err)
			}
			isRivalInt := 0
			if constants.GRivals[char] == oppo {
				isRivalInt = 1
			}
			_, err = rConn.Set(fmt.Sprintf("g:%s:%s:rival", char, oppo), isRivalInt)
			if err != nil {
				panic(err)
			}

			//COMMENT OUT FOR REPEATS
			// _, err = rConn.Set(fmt.Sprintf("g:%s:%s:overall:wins", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = rConn.Set(fmt.Sprintf("g:%s:%s:overall:losses", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = rConn.Set(fmt.Sprintf("g:%s:%s:overall:streak", char, oppo), 0)
			// if err != nil {
			// 	panic(err)
			// }
		}
	}

}
