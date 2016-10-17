package retroengine

import "fmt"

//go:generate stringer -type=State
const (
	WaitingForPlayers State = iota
	WaitingForCountdown
	CountingDown
	Started
	Finished
)

type State int

func (s *State) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, s.String())), nil
}
