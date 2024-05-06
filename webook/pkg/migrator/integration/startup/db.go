package startup

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var localhost = "43.132.133.178"

func InitSrcDB() *gorm.DB {
	return initDB("weboook")
}

func InitDstDB() *gorm.DB {
	return initDB("weboook_intr")
}


func initDB(dbName string) *gorm.DB {
	dsn := fmt.Sprintf("root:root@tcp(%s:13316)/%s", localhost, dbName)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("等待连接 MySQL", err)
	}
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return db
}
