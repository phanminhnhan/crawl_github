package model


type Role int

const (

	MEMBER Role = iota
	ADMIN
	SUPPERADMIN

	)

func (r Role)String()string{
	return []string{"MEMBER", "ADMIN", "SUPPERADMIN"}[r]
}