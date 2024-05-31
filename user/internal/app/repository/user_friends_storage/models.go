package repository

import (
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/google/uuid"
)

type userFriendRow struct {
	UserId   uuid.UUID `db:"user_id"`
	FriendId uuid.UUID `db:"friend_id"`
}

func (r *userFriendRow) ValuesMap() map[string]any {
	return map[string]any{
		"user_id":   r.UserId,
		"friend_id": r.FriendId,
	}
}

func columns() []string {
	return []string{
		"user_id",
		"friend_id",
	}
}

func (r *userFriendRow) Values(columns ...string) []any {
	values := make([]any, 0, len(columns))
	m := r.ValuesMap()

	for i := range columns {
		values = append(values, m[columns[i]])
	}

	return values
}

func newUserFriendsRowFromModelsUserFriends(friend *models.UserFriends) *userFriendRow {
	return &userFriendRow{
		UserId:   uuid.UUID(friend.UserId),
		FriendId: uuid.UUID(friend.FriendId),
	}
}

func newUserFriendsModelsFromUserFriendsRow(userFriendRow *userFriendRow) *models.UserFriends {
	return &models.UserFriends{
		UserId:   models.UserID(userFriendRow.UserId),
		FriendId: models.UserID(userFriendRow.FriendId),
	}
}

type friendRow struct {
	ID             uuid.UUID `db:"id"`
	Login          string    `db:"login"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	AvatarPhotoUrl *string   `db:"avatar_photo_url"`
}

func (r *friendRow) ValuesMap() map[string]any {
	return map[string]any{
		"id":               r.ID,
		"login":            r.Login,
		"name":             r.Name,
		"email":            r.Email,
		"avatar_photo_url": r.AvatarPhotoUrl,
	}
}

func friendColumns() []string {
	return []string{
		"id",
		"login",
		"name",
		"email",
		"avatar_photo_url",
	}
}

func (r *friendRow) Values(columns ...string) []any {
	values := make([]any, 0, len(columns))
	m := r.ValuesMap()

	for i := range columns {
		values = append(values, m[columns[i]])
	}

	return values
}

func newFriendRowFromModelsFriend(friend *models.Friend) (*friendRow, error) {
	return &friendRow{
		ID:             uuid.UUID(friend.UserId),
		Login:          friend.Login,
		Name:           friend.Name,
		Email:          friend.Email,
		AvatarPhotoUrl: &friend.AvatarPhotoUrl,
	}, nil
}

func newFriendModelsFromFriendRow(friendRow *friendRow) *models.Friend {
	return &models.Friend{
		UserId:         models.UserID(friendRow.ID),
		Login:          friendRow.Login,
		Name:           friendRow.Name,
		Email:          friendRow.Email,
		AvatarPhotoUrl: *friendRow.AvatarPhotoUrl,
	}
}
