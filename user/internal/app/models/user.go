package models

import "github.com/google/uuid"

type UserID uuid.UUID

// String - represent UserID as string
func (v UserID) String() string {
	return uuid.UUID(v).String()
}

type User struct {
	Id             UserID
	Login          string
	Name           string
	Email          string
	Password       string
	AvatarPhotoUrl string
}

type UserFriendsInfo struct {
	Friends []*User
}

type ActionInfo struct {
	Success bool
}
