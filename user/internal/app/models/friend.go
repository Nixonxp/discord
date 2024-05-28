package models

import "github.com/google/uuid"

type InviteId uuid.UUID

const (
	PendingStatus  string = "pending"
	AcceptedStatus string = "accepted"
	DeclineStatus  string = "decline"
)

func (v InviteId) String() string {
	return uuid.UUID(v).String()
}

type Friend struct {
	UserId         UserID
	Login          string
	Name           string
	Email          string
	Password       string
	AvatarPhotoUrl string
}

type FriendInvite struct {
	InviteId InviteId
	OwnerId  UserID
	UserId   UserID
	Status   string
}

type UserInvitesInfo struct {
	Invites []*FriendInvite
}

type UserFriendsInfo struct {
	Friends []*User
}
