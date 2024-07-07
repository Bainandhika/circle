package main

import (
	"circle-fiber/app/config" // use package initialization method to call init function
	"circle-fiber/app/pkg/delivery/router"
	"circle-fiber/lib/connection/database"
	"circle-fiber/lib/connection/nosql"
	"circle-fiber/lib/logger" // func init called too
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func main() {
	go func() {
		sch := gocron.NewScheduler()
		logger.InitializeLoggerScheduler(sch)
		<-sch.Start()
	}()

	mysql := database.MySQLConnect()
	mysqlDB, _ := mysql.DB()

	redis := nosql.RedisConnect()

	defer func() {
		mysqlDB.Close()
		redis.Close()
	}()

	r := router.SetupRouter(mysql, redis)

	if err := r.Listen(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		logger.Error.Fatalln("Error start circle, err: " + err.Error())
	}
}
