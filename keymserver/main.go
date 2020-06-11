package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/shellow/keyman"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

var Logger *zap.SugaredLogger

var LISTENADDR string
var RedisAddr string
var RedisPool *redis.Pool
var RedisPass string
var Keym *keyman.Keyman

func main() {
	initApp()
	server()
}

func initarg() {
	flag.StringVar(&RedisPass, "rpass", "passwd", "redis passwd")
	flag.StringVar(&LISTENADDR, "addr", ":8080", "listen address")
	flag.StringVar(&RedisAddr, "raddr", "127.0.0.1:6379", "redis address")
	flag.Parse()
}

func initApp() {
	gin.SetMode(gin.ReleaseMode)
	initarg()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	Logger = logger.Sugar()

	// init redis
	RedisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   100,
		IdleTimeout: 5 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", RedisAddr,
				redis.DialPassword(RedisPass),
				redis.DialDatabase(0),
				redis.DialConnectTimeout(3*time.Second),
				redis.DialReadTimeout(3*time.Second),
				redis.DialWriteTimeout(3*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	Logger.Info("init finish")
}

func server() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "Hello World")
	})
	router.POST("/evimem/enable", Keym.Enable)
	router.POST("/evimem/addkey", Keym.Addkey)
	router.POST("/evimem/delkey", Keym.Delkey)
	router.POST("/evimem/getkey", Keym.Getkey)
	router.GET("/evimem/listkey", Keym.Listkey)
	router.POST("/evimem/diskey", Keym.Diskey)

	s := &http.Server{
		Addr:           LISTENADDR,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}

	Logger.Info("server run")
	err := s.ListenAndServe()
	if err != nil {
		Logger.Error(err)
		os.Exit(-1)
	}
}
