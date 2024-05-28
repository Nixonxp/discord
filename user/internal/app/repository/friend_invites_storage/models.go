package repository

import (
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/google/uuid"
)

type friendInviteRow struct {
	InviteId uuid.UUID `db:"id"`
	OwnerId  uuid.UUID `db:"owner_id"`
	UserId   uuid.UUID `db:"user_id"`
	Status   string    `db:"status"`
}

func (r *friendInviteRow) ValuesMap() map[string]any {
	return map[string]any{
		"id":       r.InviteId,
		"owner_id": r.OwnerId,
		"user_id":  r.UserId,
		"status":   r.Status,
	}
}

func columns() []string {
	return []string{
		"id",
		"owner_id",
		"user_id",
		"status",
	}
}

func (r *friendInviteRow) Values(columns ...string) []any {
	values := make([]any, 0, len(columns))
	m := r.ValuesMap()

	for i := range columns {
		values = append(values, m[columns[i]])
	}

	return values
}

func newFriendInviteRowFromModelsFriendInvite(invite *models.FriendInvite) *friendInviteRow {
	return &friendInviteRow{
		InviteId: uuid.UUID(invite.InviteId),
		OwnerId:  uuid.UUID(invite.OwnerId),
		UserId:   uuid.UUID(invite.UserId),
		Status:   invite.Status,
	}
}

func newFriendInviteModelsFromFriendInviteRow(friendInviteRow *friendInviteRow) *models.FriendInvite {
	return &models.FriendInvite{
		InviteId: models.InviteId(friendInviteRow.InviteId),
		OwnerId:  models.UserID(friendInviteRow.OwnerId),
		UserId:   models.UserID(friendInviteRow.UserId),
		Status:   friendInviteRow.Status,
	}
}
