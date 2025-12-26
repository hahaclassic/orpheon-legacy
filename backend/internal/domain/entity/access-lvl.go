package entity

type AccessLevel int

const (
	Unauthorized AccessLevel = iota
	User
	Admin
)

func (a AccessLevel) String() string {
	return []string{"Unauthorized", "User", "Admin"}[a]
}
