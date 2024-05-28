package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Nixonxp/discord/user/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	"github.com/Nixonxp/discord/user/pkg/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const userTable = "users"

type PGUserRepository struct {
	queryBuilder sq.StatementBuilderType
	conn         *postgres.Connection
}

func NewUserPostgresqlRepository(conn *postgres.Connection) *PGUserRepository {
	return &PGUserRepository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		conn:         conn,
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

	if _, err = r.conn.Execx(ctx, query); err != nil {
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

	result, err := r.conn.Execx(ctx, query)
	if err != nil {
		return pkgerrors.Wrap("update user exec error repo", err)
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
	if err := r.conn.Getx(ctx, user, query); err != nil {
		if err.Error() == ErrNoRows.Error() {
			return nil, pkgerrors.Wrap("get user not found error repo", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("get user exec error repo", err)
	}

	resultUser, err := newUserModelsFromUserRow(user)
	if err != nil {
		return nil, pkgerrors.Wrap("error map row to user model", err)
	}

	return resultUser, nil
}

func (r *PGUserRepository) GetUserById(ctx context.Context, userId models.UserID) (*models.User, error) {
	query := sq.Select(columns()...).
		From(userTable).
		Where(sq.Eq{"id": userId.String()}).
		PlaceholderFormat(sq.Dollar)

	user := &userRow{}
	if err := r.conn.Getx(ctx, user, query); err != nil {
		if err.Error() == ErrNoRows.Error() {
			return nil, pkgerrors.Wrap("get user not found error repo", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("get user exec error repo", err)
	}

	resultUser, err := newUserModelsFromUserRow(user)
	if err != nil {
		return nil, pkgerrors.Wrap("error map row to user model", err)
	}

	return resultUser, nil
}
