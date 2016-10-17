package retroengine

import "fmt"

//go:generate stringer -type=Role
const (
	Admin Role = iota
	Normal
)

type Role int

func (r *Role) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, r.String())), nil
}
