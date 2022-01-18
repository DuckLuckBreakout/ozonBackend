package models

import (
	uuid "github.com/satori/go.uuid"
)

const (
	ExpireSessionCookie = 90 * 24 * 3600
)

type Session struct {
	Value  string
	UserId uint64
}

type UserId struct {
	Id uint64
}

func NewSession(userId uint64) *Session {
	newValue := uuid.NewV4()
	return &Session{
		Value:  newValue.String(),
		UserId: userId,
	}
}

type DtoSession struct {
	Value  string
	UserId uint64
}

type DtoUserId struct {
	Id uint64
}
