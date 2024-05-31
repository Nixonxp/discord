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

const friendInvitesTable = "user_friends"
const userTable = "users"

type PGUserFriendsRepository struct {
	queryBuilder sq.StatementBuilderType
	driver       repository.QueryEngineProvider
	log          *log.Logger
}

func NewUserFriendsPostgresqlRepository(driver repository.QueryEngineProvider, log *log.Logger) *PGUserFriendsRepository {
	return &PGUserFriendsRepository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		driver:       driver,
		log:          log,
	}
}

func (r *PGUserFriendsRepository) AddToFriend(ctx context.Context, friendsInfo []*models.UserFriends) error {
	query := sq.Insert(friendInvitesTable).Columns(columns()...).PlaceholderFormat(sq.Dollar)
	for _, friendsInfoItem := range friendsInfo {
		rowItem := newUserFriendsRowFromModelsUserFriends(friendsInfoItem)
		query = query.Values(rowItem.UserId, rowItem.FriendId)
	}

	if _, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			r.log.WithContext(ctx).WithError(err).Error("add to friend friend exec unique error repo")
			return models.ErrAlreadyExists
		}
		return pkgerrors.Wrap("add to friend exec error repo", err)
	}

	return nil
}

func (r *PGUserFriendsRepository) GetUserFriendsByUserId(ctx context.Context, userId models.UserID) ([]*models.Friend, error) {
	query, _, err := sq.Select("id",
		"login",
		"name",
		"email",
		"avatar_photo_url").
		From(friendInvitesTable).
		Join(fmt.Sprintf("%s ON %s.user_id = %s.id", userTable, friendInvitesTable, userTable)).
		Suffix(fmt.Sprintf("WHERE friend_id::text = '%s'", userId.String())).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.driver.GetQueryEngine(ctx).Query(ctx, query)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("get user friend invites query error repo")
		return nil, err
	}

	friends := make([]*models.Friend, 0)

	defer rows.Close()
	for rows.Next() {
		var friendRow friendRow
		if err := rows.Scan(&friendRow.ID,
			&friendRow.Login,
			&friendRow.Name,
			&friendRow.Email,
			&friendRow.AvatarPhotoUrl,
		); err != nil {
			return nil, pkgerrors.Wrap("row scan user friend invites error repo", err)
		}

		friends = append(friends, newFriendModelsFromFriendRow(&friendRow))
	}
	if err := rows.Err(); err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("get user friend invites query rows error repo")
		return nil, err
	}

	return friends, nil
}

func (r *PGUserFriendsRepository) DeleteFriend(ctx context.Context, userId models.UserID, friendId models.UserID) error {
	query := sq.Delete(friendInvitesTable).
		Suffix(fmt.Sprintf("WHERE (friend_id::text = '%s' AND user_id::text = '%s') OR (user_id::text = '%s' AND friend_id::text = '%s')",
			friendId.String(), userId.String(), friendId.String(), userId.String()),
		).
		PlaceholderFormat(sq.Dollar)

	result, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query)
	if err != nil {
		r.log.WithContext(ctx).WithError(err).WithField("userId", userId.String()).Error("delete friend exec error repo")
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
