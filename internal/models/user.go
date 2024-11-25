package models

import "errors"

const (
	MAX_ROLE_LENGHT  = 10
	MAX_TOKEN_LENGHT = 64
)

type User struct {
	UID   int    `json:"uid"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

var (
	RoleLengthErr  = errors.New("role name too long (max length: 10)")
	TokenLengthErr = errors.New("token too long (max length: 64)")
)

func (u *User) Validate() error {
	if len(u.Token) > MAX_TOKEN_LENGHT {
		return RoleLengthErr
	}
	if len(u.Role) > MAX_ROLE_LENGHT {
		return TokenLengthErr
	}
	return nil
}
