package domain

func NewDomains() []any {
	return []any{&Status{}, &Task{}, &User{}}
}
