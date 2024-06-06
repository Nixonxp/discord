package repository

import (
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/google/uuid"
)

type userRow struct {
	ID             uuid.UUID `db:"id"`
	Login          string    `db:"login"`
	Name           string    `db:"name"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	AvatarPhotoUrl *string   `db:"avatar_photo_url"`
	OauthId        *string   `db:"oauth_id"`
}

func (r *userRow) ValuesMap() map[string]any {
	return map[string]any{
		"id":               r.ID,
		"login":            r.Login,
		"name":             r.Name,
		"email":            r.Email,
		"password":         r.Password,
		"avatar_photo_url": r.AvatarPhotoUrl,
		"oauth_id":         r.OauthId,
	}
}

func columns() []string {
	return []string{
		"id",
		"login",
		"name",
		"email",
		"password",
		"avatar_photo_url",
		"oauth_id",
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
		ID:             uuid.UUID(user.Id),
		Login:          user.Login,
		Name:           user.Name,
		Email:          user.Email,
		Password:       user.Password,
		AvatarPhotoUrl: &user.AvatarPhotoUrl,
		OauthId:        &user.OauthId,
	}, nil
}

func newUserModelsFromUserRow(userRow *userRow) (*models.User, error) {
	var str string
	if userRow.AvatarPhotoUrl == nil {
		userRow.AvatarPhotoUrl = &str
	}

	if userRow.OauthId == nil {
		userRow.OauthId = &str
	}

	return &models.User{
		Id:             models.UserID(userRow.ID),
		Login:          userRow.Login,
		Name:           userRow.Name,
		Email:          userRow.Email,
		Password:       userRow.Password,
		AvatarPhotoUrl: *userRow.AvatarPhotoUrl,
		OauthId:        *userRow.OauthId,
	}, nil
}
