package main

import (
	"fmt"

	"circle-2.0/app/config" // use package initialization method to call init function
	"circle-2.0/app/pkg/delivery/router"
	"circle-2.0/lib/connection/database"
	"circle-2.0/lib/connection/nosql"
	"circle-2.0/lib/logger" // func init called too

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
