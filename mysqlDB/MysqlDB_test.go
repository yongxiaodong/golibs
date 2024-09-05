package mysqlDB

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMysqlGORM(t *testing.T) {
	assert := assert2.New(t)
	db := NewMysqlGORM(ConnParams{
		Addr:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "1317665590",
		DBName:   "sys",
	}, WithMaxOpenConn(100))
	r := db.Exec("show tables")
	assert.Equal(int64(0), r.RowsAffected, "aaa")

}
