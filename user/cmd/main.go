package main

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"user/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func httpErrorMsg(err error) *models.ErrorMessage {
	if err == nil {
		return nil
	}
	return &models.ErrorMessage{
		Message: err.Error(),
	}
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	var request models.UpdateUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.UpdateUserResponse{
		ID:    int64(userID),
		Login: "user_login",
		Email: "test@mail.ru",
		Name:  "test",
	}
	return c.JSON(http.StatusOK, response)
}

func getUserByLogin(c echo.Context) error {
	login := c.Param("login")

	response := models.GetUserResponse{
		ID:    1,
		Login: login,
		Email: "test@mail.ru",
		Name:  "test",
	}
	return c.JSON(http.StatusOK, response)
}

func addUserToFriend(c echo.Context) error {
	_ = c.Param("user_id")

	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func removeUserFromFriend(c echo.Context) error {
	_ = c.Param("user_id")

	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func acceptUserToFriend(c echo.Context) error {
	_ = c.Param("invite_id")

	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func declineUserToFriend(c echo.Context) error {
	_ = c.Param("invite_id")

	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func getUserFriends(c echo.Context) error {
	response := models.FriendsResponse{
		&models.Friend{
			UserID: 1,
			Login:  "user_login",
			Email:  "test@mail.ru",
			Name:   "test",
		},
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.PUT("/api/v1/users/:id", updateUser)
	e.GET("/api/v1/users/search/login/:login", getUserByLogin)
	e.GET("/api/v1/friends", getUserFriends)
	e.POST("/api/v1/friends/:user_id", addUserToFriend)
	e.DELETE("/api/v1/friends/:user_id", removeUserFromFriend)
	e.POST("/api/v1/friends/invite/accept/:invite_id", acceptUserToFriend)
	e.POST("/api/v1/friends/invite/decline/:invite_id", declineUserToFriend)

	e.GET("/health", func(c echo.Context) error {
		status := http.StatusOK
		statusMessage := "OK"

		if !isServiceOk(10) {
			status = http.StatusInternalServerError
			statusMessage = "Error"
		}

		return c.JSON(status, struct{ Status string }{Status: statusMessage})
	})

	e.GET("/ready", func(c echo.Context) error {
		status := http.StatusOK
		statusMessage := "OK"

		if !isServiceOk(5) {
			status = http.StatusInternalServerError
			statusMessage = "Error"
		}

		return c.JSON(status, struct{ Status string }{Status: statusMessage})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// isServiceOk в зависимости от входящего значения вернет false, например
// передано 5, тогда (100 / 5 = 20) 20% вероятностью вернется false, для теста сервиса
func isServiceOk(probability int) bool {
	randNumber := rand.IntN(probability-1) + 1

	if randNumber == 1 {
		return false
	}

	return true
}
