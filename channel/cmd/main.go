package main

import (
	"channel/internal/models"
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"math/rand/v2"
	"net/http"
	"os"

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

func addChannel(c echo.Context) error {
	var request models.AddChannelRequest
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

func deleteChannel(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func joinChannel(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func leaveChannel(c echo.Context) error {
	response := models.SuccessResponse{
		Success: true,
	}
	return c.JSON(http.StatusOK, response)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/channel", addChannel)
	e.DELETE("/api/v1/channel/:channel_idß", deleteChannel)
	e.POST("/api/v1/channel/join/:channel_id", joinChannel)
	e.POST("/api/v1/channel/leave/:channel_id", leaveChannel)

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
