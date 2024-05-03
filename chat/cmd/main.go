package main

import (
	"chat/internal/models"
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

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

func sendUserPrivateMessage(c echo.Context) error {
	var request models.SendPrivateMessageRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func getUserPrivateMessages(c echo.Context) error {
	response := models.PrivateMessagesResponse{
		&models.Message{
			ID:        1,
			Text:      "text message",
			Timestamp: time.Now().Unix(),
		},
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.PUT("/api/v1/chat/private/:user_id", sendUserPrivateMessage)
	e.GET("/api/v1/chat/private/:user_id", getUserPrivateMessages)

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
