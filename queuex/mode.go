package queuex

// Mode represents queue behavior mode
type Mode int

const (
	ModeFIFO Mode = iota // First In First Out (default)
	ModeLIFO             // Last In First Out
)
