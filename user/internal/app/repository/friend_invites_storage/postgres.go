package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Nixonxp/discord/user/internal/app/models"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const friendInvitesTable = "friend_invites"

type PGFriendInvitesRepository struct {
	queryBuilder sq.StatementBuilderType
	driver       repository.QueryEngineProvider
	log          *log.Logger
}

func NewFriendInvitesPostgresqlRepository(driver repository.QueryEngineProvider, log *log.Logger) *PGFriendInvitesRepository {
	return &PGFriendInvitesRepository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		driver:       driver,
		log:          log,
	}
}

func (r *PGFriendInvitesRepository) CreateInvite(ctx context.Context, invite *models.FriendInvite) error {
	row := newFriendInviteRowFromModelsFriendInvite(invite)
	query := sq.Insert(friendInvitesTable).
		SetMap(row.ValuesMap()).
		PlaceholderFormat(sq.Dollar)

	if _, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			r.log.WithContext(ctx).WithError(err).WithField("owner_id", invite.OwnerId.String()).Error("create friend invite exec unique error repo")
			return models.ErrAlreadyExists
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

	rows, err := r.driver.GetQueryEngine(ctx).Query(ctx, query)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("get user friend invites query error repo")
		return nil, err
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
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("get user friend invites query rows error repo")
		return nil, err
	}

	return &models.UserInvitesInfo{
		Invites: invites,
	}, nil
}

func (r *PGFriendInvitesRepository) GetInviteById(ctx context.Context, inviteId models.InviteId) (*models.FriendInvite, error) {
	query := sq.Select(columns()...).
		From(friendInvitesTable).
		Suffix(fmt.Sprintf("WHERE id::text = '%s'", inviteId.String())).
		PlaceholderFormat(sq.Dollar)

	inviteRow := &friendInviteRow{}

	if err := r.driver.GetQueryEngine(ctx).Getx(ctx, inviteRow, query); err != nil {
		if err.Error() == repository.ErrNoRows.Error() {
			r.log.WithContext(ctx).WithError(err).WithField("inviteId", inviteId.String()).Error("invite not found in repository")
			return nil, pkgerrors.Wrap("invite not found", models.ErrNotFound)
		}

		r.log.WithContext(ctx).WithError(err).WithField("inviteId", inviteId.String()).Error("get invite exec error repo")
		return nil, err
	}

	resultInvite := newFriendInviteModelsFromFriendInviteRow(inviteRow)
	return resultInvite, nil
}

func (r *PGFriendInvitesRepository) DeclineInvite(ctx context.Context, inviteId models.InviteId) error {
	query := sq.Update(friendInvitesTable).
		SetMap(map[string]any{
			"status": models.DeclineStatus,
		}).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("WHERE id::text = '%s' AND status = '%s'", inviteId.String(), models.PendingStatus))

	result, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("inviteId", inviteId.String()).Error("decline invite friend exec error repo")
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *PGFriendInvitesRepository) AcceptInvite(ctx context.Context, inviteId models.InviteId) error {
	query := sq.Update(friendInvitesTable).
		SetMap(map[string]any{
			"status": models.AcceptedStatus,
		}).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("WHERE id::text = '%s' AND status = '%s'", inviteId.String(), models.PendingStatus))

	result, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("inviteId", inviteId.String()).Error("accept invite friend exec error repo")
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
