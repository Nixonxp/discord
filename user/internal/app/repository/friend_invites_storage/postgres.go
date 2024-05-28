package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Nixonxp/discord/user/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	"github.com/Nixonxp/discord/user/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const friendInvitesTable = "friend_invites"

type PGFriendInvitesRepository struct {
	queryBuilder sq.StatementBuilderType
	conn         *postgres.Connection
}

func NewFriendInvitesPostgresqlRepository(conn *postgres.Connection) *PGFriendInvitesRepository {
	return &PGFriendInvitesRepository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		conn:         conn,
	}
}

func (r *PGFriendInvitesRepository) CreateInvite(ctx context.Context, invite *models.FriendInvite) error {
	row := newFriendInviteRowFromModelsFriendInvite(invite)
	query := sq.Insert(friendInvitesTable).
		SetMap(row.ValuesMap()).
		PlaceholderFormat(sq.Dollar)

	if _, err := r.conn.Execx(ctx, query); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap("create friend invite exec unique error repo", models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap("create friend invite exec error repo", err)
	}

	return nil
}

func (r *PGFriendInvitesRepository) GetInvitesByUserId(ctx context.Context, userId models.UserID) (*models.UserInvitesInfo, error) {
	query, _, err := sq.Select(columns()...).
		From(friendInvitesTable).
		Suffix(fmt.Sprintf("WHERE user_id::text = '%s'", userId.String())).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, pkgerrors.Wrap("get user friend invites query error repo", err)
	}

	invites := make([]*models.FriendInvite, 0)

	defer rows.Close()
	for rows.Next() {
		var friendInviteRowItem friendInviteRow
		if err := rows.Scan(&friendInviteRowItem.InviteId,
			&friendInviteRowItem.OwnerId,
			&friendInviteRowItem.UserId,
			&friendInviteRowItem.Status,
		); err != nil {
			return nil, pkgerrors.Wrap("row scan user friend invites error repo", err)
		}

		invites = append(invites, newFriendInviteModelsFromFriendInviteRow(&friendInviteRowItem))
	}
	if err := rows.Err(); err != nil {
		return nil, pkgerrors.Wrap("get user friend invites query rows error repo", err)
	}

	return &models.UserInvitesInfo{
		Invites: invites,
	}, nil
}

func (r *PGFriendInvitesRepository) DeclineInvite(ctx context.Context, inviteId models.InviteId) error {
	query := sq.Update(friendInvitesTable).
		SetMap(map[string]any{
			"status": models.DeclineStatus,
		}).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("WHERE id::text = '%s' AND status = '%s'", inviteId.String(), models.PendingStatus))

	result, err := r.conn.Execx(ctx, query)
	if err != nil {
		return pkgerrors.Wrap("decline invite friend exec error repo", err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
