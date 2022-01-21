package main

import (
	"time"

	"github.com/universe-30/GormULog"
	"github.com/universe-30/LogrusULog"
	"github.com/universe-30/ULog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Person struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"index"`
	Age  int    `gorm:"index"`
}

func main() {
	//logger instance
	logger, _ := LogrusULog.New("./logs", 2, 20, 30)
	logger.SetLevel(ULog.DebugLevel)

	//new db
	db, err := gorm.Open(sqlite.Open("./sqlite.db"), &gorm.Config{
		//use custom logger
		Logger: GormULog.New_gormLocalLogger(logger, GormULog.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  gormlogger.Info, //Level: Silent Error Warn Info. Info will log all record.
		}),
	})
	if err != nil {
		logger.Fatalln(err)
	}

	err = db.AutoMigrate(
		&Person{},
	)
	if err != nil {
		logger.Fatalln(err)
	}

	p := Person{Name: "Jack", Age: 18}
	db.Create(&p)

	var qp Person
	db.First(&qp)
	logger.Debugln(qp)

	//user Debug() to trace sql
	db.Debug().Last(&qp)
	logger.Debugln(qp)
}
