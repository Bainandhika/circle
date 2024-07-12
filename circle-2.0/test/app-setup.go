package test

import (
	"circle-2.0/app/pkg/delivery/router"
	"circle-2.0/lib/connection/database"
	"circle-2.0/lib/connection/nosql"
	"circle-2.0/lib/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jasonlvhit/gocron"
)

func SetUpTestApp() *fiber.App {
	go func() {
		sch := gocron.NewScheduler()
		logger.InitializeLoggerScheduler(sch)
		<-sch.Start()
	}()

	database.InitMySQLConnection()
	mysql := database.MySQLConnect
	mysqlDB, _ := mysql.DB()

	redis := nosql.RedisConnect()

	defer func() {
		mysqlDB.Close()
		redis.Close()
	}()

	return router.SetupRouter(mysql, redis)
}
