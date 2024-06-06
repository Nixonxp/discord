package models

import "github.com/google/uuid"

type UserID uuid.UUID

// String - represent UserID as string
func (v UserID) String() string {
	return uuid.UUID(v).String()
}

type User struct {
	UserID         UserID
	Login          string
	Name           string
	Email          string
	Password       string
	AvatarPhotoUrl string
}

type Login struct {
	Login    string
	Password string
}
