package mysqlDB

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type ConnParams struct {
	Addr     string
	Port     int
	User     string
	Password string
	DBName   string
}

type GORMOpts struct {
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime time.Duration
	MAxLifeTime time.Duration
}

type Option func(db *gorm.DB)

func WithMaxOpenConn(param int) Option {
	return func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxOpenConns(param)
	}
}

func WithMaxIdleConn(param int) Option {
	return func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetMaxIdleConns(param)
	}
}

func WithMaxLifeTime(param time.Duration) Option {
	return func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDB.SetConnMaxLifetime(param)
	}
}

func WithMaxIdleTime(param time.Duration) Option {
	return func(db *gorm.DB) {
		sqlDb, err := db.DB()
		if err != nil {
			panic(err)
		}
		sqlDb.SetConnMaxIdleTime(param)

	}

}

func applyDefaultOption(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(100)                // 最大链接
	sqlDb.SetMaxIdleConns(50)                 // 空闲链接
	sqlDb.SetConnMaxLifetime(time.Minute * 5) // 最大生存时长
	sqlDb.SetConnMaxIdleTime(time.Minute * 3)
}

// 检查和修正DSN
func checkAndFixDSN(dsn string) string {
	sl := len(dsn)
	defaultLoc := "Local"

	// 如果没有包含 parseTime=true 参数，添加它
	if !strings.Contains(dsn, "parseTime=true") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}

	// 如果没有包含 loc 参数，添加它
	if !strings.Contains(dsn, "loc=") {
		if strings.Contains(dsn, "?") {
			dsn += "&loc=" + defaultLoc
		} else {
			dsn += "?loc=" + defaultLoc
		}
	}
	if sl != len(dsn) {
		log.Println("DSN fixed")
	}
	return dsn
}

func NewMysqlGORM(myp ConnParams, opts ...Option) *gorm.DB {
	log.Println("Conn database: ", myp.DBName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", myp.User, myp.Password, myp.Addr, myp.Port, myp.DBName)
	dsn = checkAndFixDSN(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("%s database conn fail", myp.DBName)
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	if len(opts) == 0 {
		applyDefaultOption(db)

	} else {
		for _, opt := range opts {
			opt(db)
		}
	}
	log.Println("Conn database successes: ", myp.DBName)
	return db
}

func NewMysqlSqlc(myp ConnParams) *sql.DB {
	log.Println("Conn database: ", myp.DBName)
	connP := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", myp.User, myp.Password, myp.Addr, myp.Port, myp.DBName)
	masterDb, err := sql.Open("mysqlDB", connP)
	if err != nil {
		log.Printf("%s database conn fail", myp.DBName)
		panic(err)
	}
	masterDb.SetMaxOpenConns(30)
	masterDb.SetMaxOpenConns(100)
	masterDb.SetConnMaxLifetime(time.Hour)
	if err := masterDb.Ping(); err != nil {
		panic(err)
	}
	log.Println("Conn database succes: ", myp.DBName)
	return masterDb
}

func NewEsConn() {

}

func NewClickhouseConn() {

}

func NewRedisConn() {

}
