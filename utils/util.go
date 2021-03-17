package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/garyburd/redigo/redis"

	_ "github.com/beego/beego/v2/server/web/session/redis"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"html/template"
	"net/http"
	"strings"
	"time"
)

//字符串进行MD5编码
func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}

//对字符串进行编译，防止注入
func FilterStr(str string) string {
	return strings.TrimSpace(template.HTMLEscapeString(html.UnescapeString(str)))
}

//过滤管理员登录
var CmsLoginFilter = func(ctx *context.Context) {
	user := ctx.Input.Session("cmc_user")

	if user == nil {
		errMsg := map[string]interface{}{
			"code":     http.StatusMethodNotAllowed,
			"message":  "请您先登录系统",
			"Jump_url": "/cms/login",
		}

		bytes, err := json.Marshal(errMsg)
		if err != nil {
			panic("系统错误")
		}
		_, _ = ctx.ResponseWriter.Write(bytes)
		return
	}
}

//redis 缓存驱动
var RedisPool *redis.Pool

//获取缓存驱动
func GetRedis() *redis.Pool {
	return RedisPool
}

//初始化 redis
func Redis() {
	fmt.Println("redis init start ...")

	pool := &redis.Pool{
		MaxActive:   30,
		MaxIdle:     5,
		IdleTimeout: time.Second * time.Duration(60*time.Millisecond),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379",
				redis.DialDatabase(0),
				redis.DialPassword(""),
			)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}

	fmt.Println("redis 初始化完成...")
	RedisPool = pool
}
