package test

import (
	"circle-2.0/app/config"
	"circle-2.0/app/pkg/delivery/router"
	"circle-2.0/lib/connection/database"
	"circle-2.0/lib/connection/nosql"
	"circle-2.0/lib/logger"
	"github.com/gofiber/fiber/v2"
)

func SetUpTestApp() *fiber.App {
	configPath := "c:/go/src/circle/circle-2.0/circle-config.env"
	config.InitConfig(configPath)

	logPath := "c:/go/src/circle/circle-2.0/test/log/"

	logConfig := logger.LoggerConfig{LogPath: logPath}
	logConfig.InitLogger()

	db := database.SetMySQL{DBConfig: config.DB, LogPath: logPath}
	db.InitMySQLConnection()

	nosql.InitRedis()

	return router.SetupRouter(database.MySQLConnect, nosql.RedisConnect)
}
