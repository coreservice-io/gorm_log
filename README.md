# gorm_log

logger for gorm implementing log interface 

example
```go

package main

import (
	"time"

	"github.com/coreservice-io/gorm_log"
	"github.com/coreservice-io/log"
	"github.com/coreservice-io/logrus_log"
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
	logger, _ := logrus_log.New("./logs", 2, 20, 30)
	logger.SetLevel(log.DebugLevel)

	//new db
	db, err := gorm.Open(sqlite.Open("./sqlite.db"), &gorm.Config{
		//use custom logger
		Logger: gorm_log.New_gormLocalLogger(logger, gorm_log.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			//Level: Silent Error Warn Info.
			//Info logs all record.
			//Silent turns off log.
			LogLevel: gorm_log.Info,
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