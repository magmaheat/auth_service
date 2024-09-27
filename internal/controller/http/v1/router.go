package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magmaheat/auth_service/internal/service"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}" "uri":${uri} "status":"${status}", "error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))

}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
