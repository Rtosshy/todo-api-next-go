package dao

func NewDomains() []any {
	return []any{&Status{}, &Task{}, &User{}}
}
