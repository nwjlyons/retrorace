package handlers

//go:generate stringer -type=Action
const (
	CloseToNewPlayers Action = iota
	Start
	Increment
	Reset
)

type Action int
