package test

import (
	"circle/app/config"
	"circle/app/pkg/delivery/router"
	"circle/lib/connection/database"
	"circle/lib/connection/nosql"
	"circle/lib/logger"

	"github.com/gofiber/fiber/v2"
)

func SetUpTestApp() *fiber.App {
	config.InitConfig()

	logPath := "c:/go/src/circle/circle/test/log/"

	logConfig := logger.LoggerConfig{LogPath: logPath}
	logConfig.InitLogger()

	db := database.SetMySQL{DBConfig: config.DB, LogPath: logPath}
	db.InitMySQLConnection()

	nosql.InitRedis()

	return router.SetupRouter(database.MySQLConnect, nosql.RedisConnect)
}
