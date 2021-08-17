package user

import (
	"ova-exam-api/internal/domain/entity"
)

type User struct {
	entity.Entity
	UserId    uint64
	Email     string
	Password  string
}

func (u *User) String() string {
	return "this is user!"
}
