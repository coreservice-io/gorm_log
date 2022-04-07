# GormULog

logger for gorm using ULog

example
```go
package main

import (
	"time"

	"github.com/coreservice-io/GormULog"
	"github.com/coreservice-io/LogrusULog"
	"github.com/coreservice-io/ULog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Person struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"index"`
	Age  int64  `gorm:"index"`
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
			//Level: Silent Error Warn Info. 
			//Info logs all record.
			//Silent turns off log.
			LogLevel:                  GormULog.Info,
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

```