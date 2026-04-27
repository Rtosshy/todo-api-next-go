package dao

func NewDAOs() []any {
	return []any{&Status{}, &Task{}, &User{}}
}
