package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Nixonxp/discord/user/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) transaction_manager.QueryEngine
}

const userTable = "users"

type PGUserRepository struct {
	queryBuilder sq.StatementBuilderType
	driver       QueryEngineProvider
	log          *log.Logger
}

func NewUserPostgresqlRepository(driver QueryEngineProvider, log *log.Logger) *PGUserRepository {
	return &PGUserRepository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		driver:       driver,
		log:          log,
	}
}

func (r *PGUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	row, err := newUserRowFromModelsUser(user)
	if err != nil {
		return pkgerrors.Wrap("create error repo", err)
	}

	query := sq.Insert(userTable).
		SetMap(row.ValuesMap()).
		PlaceholderFormat(sq.Dollar)

	if _, err = r.driver.GetQueryEngine(ctx).Execx(ctx, query); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap("create user exec unique error repo", models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap("create user exec error repo", err)
	}

	return nil
}

func (r *PGUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := sq.Update(userTable).
		SetMap(map[string]any{
			"login":            user.Login,
			"name":             user.Name,
			"email":            user.Email,
			"avatar_photo_url": user.AvatarPhotoUrl,
		}).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": user.Id.String()})

	result, err := r.driver.GetQueryEngine(ctx).Execx(ctx, query)
	if err != nil {
		r.ctxLog(ctx).WithError(err).WithField("user", user).Error("update user exec error repo")
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

var ErrNoRows = fmt.Errorf("scanning one: no rows in result set")

func (r *PGUserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := sq.Select(columns()...).
		From(userTable).
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar)

	user := &userRow{}
	if err := r.driver.GetQueryEngine(ctx).Getx(ctx, user, query); err != nil {
		if err.Error() == ErrNoRows.Error() {
			r.ctxLog(ctx).WithError(err).WithField("login", login).Error("user not found in repository")
			return nil, pkgerrors.Wrap("user not found", models.ErrNotFound)
		}

		r.ctxLog(ctx).WithError(err).WithField("login", login).Error("get user exec error repo")
		return nil, err
	}

	resultUser, err := newUserModelsFromUserRow(user)
	if err != nil {
		r.ctxLog(ctx).WithError(err).WithField("login", login).Error("error map row to user model")
		return nil, err
	}

	return resultUser, nil
}

func (r *PGUserRepository) GetUserByOauthId(ctx context.Context, oauthId string) (*models.User, error) {
	query := sq.Select(columns()...).
		From(userTable).
		Where(sq.Eq{"oauth_id": oauthId}).
		PlaceholderFormat(sq.Dollar)

	user := &userRow{}
	if err := r.driver.GetQueryEngine(ctx).Getx(ctx, user, query); err != nil {
		if err.Error() == ErrNoRows.Error() {
			r.ctxLog(ctx).WithError(err).WithField("oauth_id", oauthId).Error("user not found in repository")
			return nil, pkgerrors.Wrap("user not found", models.ErrNotFound)
		}

		r.ctxLog(ctx).WithError(err).WithField("oauth_id", oauthId).Error("get user exec error repo")
		return nil, err
	}

	resultUser, err := newUserModelsFromUserRow(user)
	if err != nil {
		r.ctxLog(ctx).WithError(err).WithField("oauth_id", oauthId).Error("error map row to user model")
		return nil, err
	}

	return resultUser, nil
}

func (r *PGUserRepository) GetUserById(ctx context.Context, userId models.UserID) (*models.User, error) {
	query := sq.Select(columns()...).
		From(userTable).
		Where(sq.Eq{"id": userId.String()}).
		PlaceholderFormat(sq.Dollar)

	user := &userRow{}
	if err := r.driver.GetQueryEngine(ctx).Getx(ctx, user, query); err != nil {
		if err.Error() == ErrNoRows.Error() {
			r.ctxLog(ctx).WithError(err).WithField("userId", userId).Error("get user not found error repo")
			return nil, models.ErrNotFound
		}

		r.ctxLog(ctx).WithError(err).WithField("userId", userId).Error("get user exec error repo")
		return nil, err
	}

	resultUser, err := newUserModelsFromUserRow(user)
	if err != nil {
		return nil, pkgerrors.Wrap("error map row to user model", err)
	}

	return resultUser, nil
}

func (r *PGUserRepository) ctxLog(ctx context.Context) *log.Logger {
	return r.log.WithContext(ctx)
}
