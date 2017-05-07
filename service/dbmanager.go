// dbmanager
package service

import (
	"database/sql"
	"fmt"
	//	"os"
	//	"os/signal"
	//	"syscall"
	//	"time"

	//	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	//	"littleform/model"
)

const (
	redisAddr = "127.0.0.1:6379"
)

var (
	MaxPoolSize = 20
	redisPool   chan redis.Conn
	sqlPool     *sql.DB
)

func init() {
	//	sqlPool, err := sql.Open("mysql", mysqlAddr)
	//	if err != nil {
	//		fmt.Printf("connect to mysql error: %s", err.Error())
	//	}
	//	sqlPool.SetMaxOpenConns(1000)
	//	sqlPool.SetMaxIdleConns(500)
	//	err = sqlPool.Ping()
	//	if err != nil {
	//		panic(err)
	//	}

	//	err := orm.RegisterDataBase("default", "mysql", mysqlLocal, 30)
	//	if err != nil {
	//		panic(err)
	//	}

	//	orm.RegisterModel(new(model.FormList), new(model.ConfigList), new(model.SubmitIdList), new(model.SubmitList))
	//	orm.RunSyncdb("default", false, true)

}

// redis pool related
func putRedis(conn redis.Conn) {
	if redisPool == nil {
		redisPool = make(chan redis.Conn, MaxPoolSize)
	}
	if len(redisPool) >= MaxPoolSize {
		conn.Close()
		return
	}
	fmt.Printf("redis conn pool size: %d\n", len(redisPool))
	redisPool <- conn
}

func InitRedis(address string) redis.Conn {
	if len(redisPool) == 0 {
		redisPool = make(chan redis.Conn, MaxPoolSize)
		go func() {
			for i := 0; i < MaxPoolSize/2; i++ {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					panic(err)
				}
				putRedis(c)
			}
		}()
	}
	return <-redisPool
}

func getRedisConn() redis.Conn {
	return InitRedis(redisAddr)
}
