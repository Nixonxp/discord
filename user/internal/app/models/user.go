package models

type User struct {
	UserID   uint64
	Login    string
	Name     string
	Email    string
	Password string
}

type UserFriendsInfo struct {
	Friends []*User
}

type ActionInfo struct {
	Success bool
}
