package main

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"math/rand/v2"
	"net/http"
	"os"
	"server/internal/models"
	"strconv"
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

func addServer(c echo.Context) error {
	var request models.CreateServerRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.CreateServerResponse{
		ID:   1,
		Name: *request.Name,
	}
	return c.JSON(http.StatusOK, response)
}

func searchServer(c echo.Context) error {
	var request models.SearchServerRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.SearchServerResponse{
		&models.Server{
			ID:   1,
			Name: "server name",
		},
	}
	return c.JSON(http.StatusOK, response)
}

func subscribeServer(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}
func unsubscribeServer(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func searchServerByUserId(c echo.Context) error {
	id := c.Param("user_id")
	_, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	response := models.SearchServerResponse{
		&models.Server{
			ID:   1,
			Name: "server name",
		},
	}
	return c.JSON(http.StatusOK, response)
}

func inviteUserToServer(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func publicationMessageOnServer(c echo.Context) error {
	var request models.SendMessageServerRequest
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

func getMessagesFromServer(c echo.Context) error {
	response := models.ServerMessagesResponse{
		&models.Message{
			ID:        1,
			Text:      "text",
			Timestamp: time.Now().Unix(),
		},
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/servers/create", addServer)
	e.POST("/api/v1/servers/search", searchServer)
	e.POST("/api/v1/servers/subscribe/:server_id", subscribeServer)
	e.POST("/api/v1/servers/unsubscribe/:server_id", unsubscribeServer)
	e.POST("/api/v1/servers/search/user/:user_id", searchServerByUserId)
	e.POST("/api/v1/servers/:server_id/invite/:user_id", inviteUserToServer)
	e.POST("/api/v1/servers/:server_id/messages/send", publicationMessageOnServer)
	e.GET("/api/v1/servers/:server_id/messages", getMessagesFromServer)

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
