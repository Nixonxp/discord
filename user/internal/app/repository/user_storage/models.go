package repository

import (
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/google/uuid"
)

type userRow struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Name     string    `db:"name"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

func (r *userRow) ValuesMap() map[string]any {
	return map[string]any{
		"id":       r.ID,
		"login":    r.Login,
		"name":     r.Name,
		"email":    r.Email,
		"password": r.Password,
	}
}

func columns() []string {
	return []string{
		"id",
		"login",
		"name",
		"email",
		"password",
	}
}

func (r *userRow) Values(columns ...string) []any {
	values := make([]any, 0, len(columns))
	m := r.ValuesMap()

	for i := range columns {
		values = append(values, m[columns[i]])
	}

	return values
}

func newUserRowFromModelsUser(user *models.User) (*userRow, error) {
	return &userRow{
		ID:       uuid.UUID(user.Id),
		Login:    user.Login,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func newUserModelsFromUserRow(userRow *userRow) (*models.User, error) {
	return &models.User{
		Id:       models.UserID(userRow.ID),
		Login:    userRow.Login,
		Name:     userRow.Name,
		Email:    userRow.Email,
		Password: userRow.Password,
	}, nil
}
