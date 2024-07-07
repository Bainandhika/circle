package main

import (
	"circle/app/config" // use package initialization method to call init function
	"circle/app/pkg/delivery/router"
	"circle/lib/connection/database"
	"circle/lib/connection/nosql"
	"circle/lib/logger" // func init called too
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

	if err := r.Run(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)); err != nil {
		logger.Error.Println("Error start circle, err: " + err.Error())
		panic(err)
	}
}
