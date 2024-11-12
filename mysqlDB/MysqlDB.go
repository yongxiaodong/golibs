package mysqlDB

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math"
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

func GORMBatchInsert(data []interface{}, destDB *gorm.DB, batchSize int) {
	if batchSize > 10000 || batchSize < 100 {
		log.Println("batchSize cannot be greater than 10000 or less than 100, set default batchSize to 500")
		batchSize = 500
	}
	total := len(data)
	for i := 0; i < total; i += batchSize {
		end := int(math.Min(float64(i+batchSize), float64(total)))
		batchData := data[i:end]
		destDB.Create(&batchData)
	}
}

type BatchOption struct {
	// 单次删除的数据
	BatchSize int
	// 执行一次后休眠时间，毫秒
	SlowTime time.Duration
}

func BatchDelete(db *gorm.DB, sql string, TName string, option ...BatchOption) *gorm.DB {
	var result *gorm.DB
	row := int64(0)
	var size = 50000
	var slowTime time.Duration = 0
	if len(option) > 0 && option[0].BatchSize > 0 {
		size = option[0].BatchSize
	}
	if len(option) > 0 && option[0].SlowTime > 0 {
		slowTime = option[0].SlowTime
	}
	for {
		result = db.Exec(sql, size)
		row = row + result.RowsAffected
		if result.Error != nil {
			break
		}
		log.Printf("分批执行成功: %s - 影响行: %d", TName, result.RowsAffected)
		if result.RowsAffected < int64(size) {
			break
		}
		time.Sleep(time.Millisecond * slowTime)
	}
	result.RowsAffected = result.RowsAffected + row
	return result
}
