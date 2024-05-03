package main

import (
	"auth/internal/models"
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"math/rand/v2"
	"net/http"
	"os"
)

func httpErrorMsg(err error) *models.ErrorMessage {
	if err == nil {
		return nil
	}
	return &models.ErrorMessage{
		Message: err.Error(),
	}
}

func registerUser(c echo.Context) error {
	var request models.RegisterUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.RegisterUserResponse{
		ID:    1,
		Login: "user_login",
		Email: "test@mail.ru",
		Name:  "test",
	}
	return c.JSON(http.StatusOK, response)
}

func loginUser(c echo.Context) error {
	var request models.LoginUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.LoginUserResponse{
		UserID: 1,
		Login:  "user_login",
		Token:  "9HqLlyZOugoStsXCUfD_0YdwnNnunAJF8V47U3QHXSq",
		Name:   "test",
	}
	return c.JSON(http.StatusOK, response)
}

func callbackOAuthUser(c echo.Context) error {
	var request models.CallbackOAuthUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.LoginUserResponse{
		UserID: 1,
		Login:  "user_login",
		Token:  "9HqLlyZOugoStsXCUfD_0YdwnNnunAJF8V47U3QHXSq",
		Name:   "test",
	}
	return c.JSON(http.StatusOK, response)
}

func loginOAuthUser(c echo.Context) error {

	response := models.LoginOAthUserResponse{
		Code: "9HqLlyZOugoStsXCUfD_0YdwnNnunAJF8V47U3QHXSq",
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/auth/register", registerUser)
	e.POST("/api/v1/auth/login", loginUser)
	e.POST("/api/v1/oauth/login", loginOAuthUser)
	e.POST("/api/v1/oauth/callback", callbackOAuthUser)

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
