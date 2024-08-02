package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"circle/app/config"
	"circle/lib/logger"
	"circle/lib/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var MySQLConnect *gorm.DB

type SetMySQL struct {
	DBConfig config.DatabaseConfig
	LogPath  string
}

func (m *SetMySQL) setMySQL() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.DBConfig.Username,
		m.DBConfig.Password,
		m.DBConfig.Host,
		m.DBConfig.Port,
		m.DBConfig.Name,
	)

	date := time.Now().Format("20060102")
	fileName := fmt.Sprintf("%scircle-gorm-%s.log", m.LogPath, date)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error.Printf("Failed to open log file! err: %v", err)
		panic(err)
	}
	defer file.Close()

	newLogger := gormLogger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,            // Don't include params in the SQL log
			Colorful:                  false,           // Disable color
		},
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
	})
	if err != nil {
		logger.Error.Printf("Failed to connect to MySQL! err: %v", err)
		panic(err)
	}

	// Automatically migrate the schema to match the struct definitions
	if err = database.AutoMigrate(&model.Users{}); err != nil {
		logger.Error.Printf("Failed to auto migrate user! err: %v", err)
		panic(err)
	}

	if err = database.AutoMigrate(&model.OrderMains{}); err != nil {
		logger.Error.Printf("Failed to auto migrate order_main! err: %v", err)
		panic(err)
	}

	if err = database.AutoMigrate(&model.OrderUsers{}); err != nil {
		logger.Error.Printf("Failed to auto migrate order_person! err: %v", err)
		panic(err)
	}

	if err = database.AutoMigrate(&model.OrderUserItems{}); err != nil {
		logger.Error.Printf("Failed to auto migrate order_person! err: %v", err)
		panic(err)
	}

	if err = database.AutoMigrate(&model.AdditionalCosts{}); err != nil {
		logger.Error.Printf("Failed to auto migrate additional_cost! err: %v", err)
		panic(err)
	}

	if err = database.AutoMigrate(&model.Discounts{}); err != nil {
		logger.Error.Printf("Failed to auto migrate discount! err: %v", err)
		panic(err)
	}

	return database
}

func (m *SetMySQL) InitMySQLConnection() {
	MySQLConnect = m.setMySQL()
}