package constants

const (
	FULGORE string = "fulgore"
	SPINAL         = "spinal"
	WULF           = "wulf"
	ORCHID         = "orchid"
	GLACIUS        = "glacius"
	JAGO           = "jago"
	THUNDER        = "thunder"
	SADIRA         = "sadira"
	TJ				= "tj"
	MAYA 			= "maya"
	KANRA			= "kanra"
	RIPTOR 			= "riptor"
)

const (
	G    string = "g"
	BJ          = "bj"
	MIKE        = "mike"
)

var XterList = []string{FULGORE, SPINAL, WULF, ORCHID, GLACIUS, JAGO, THUNDER, SADIRA, TJ, MAYA, KANRA, RIPTOR}
var PlayerList = []string{G, BJ}

// map[G_XTER]BJ_XTER
var GRivals = map[string]string{FULGORE: THUNDER, SPINAL: WULF, JAGO: JAGO, THUNDER: FULGORE, SADIRA: SADIRA, TJ: TJ, MAYA: MAYA, KANRA: KANRA, RIPTOR: RIPTOR}
var BJRivals = map[string]string{THUNDER: FULGORE, WULF: SPINAL, JAGO: JAGO, FULGORE: THUNDER, SADIRA: SADIRA, TJ: TJ, MAYA: MAYA, KANRA: KANRA, RIPTOR: RIPTOR}

var WinningScore = 150
var WinningScore2 = 200