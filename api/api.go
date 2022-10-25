package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"podloxx-collector/structs"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/mogenius/mo-go/logger"

	"github.com/gin-gonic/gin"
)

const MAXPOLL int64 = 100

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const REDISCONSTR = "127.0.0.1:6379"

var redisClient *redis.Client

func TestRedis() {
	for {
		if initRedis() {
			break
		}
		logger.Log.Infof("Waiting for redis to come alive: %s", REDISCONSTR)
		time.Sleep(1 * time.Second)
	}
	data := getRedisData()
	logger.Log.Infof("Received: %d", len(data))
}

func InitApi() {
	for {
		if initRedis() {
			break
		}
		logger.Log.Infof("Waiting for redis to come alive: %s", REDISCONSTR)
		time.Sleep(1 * time.Second)
	}
	go initGin()
}

func InitApiCluster() {
	initGin()
}

func initGin() {
	router := gin.Default()
	router.GET("/traffic", getTraffic)
	router.StaticFile("/traffic-test", "./ui/test-ws.html")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info("listen: %s\n", err)
	}
}

func initRedis() bool {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     REDISCONSTR,
		Password: "",
		DB:       0,
	})

	// CHECK CONNECTION
	_, err := redisClient.Ping().Result()
	if err != nil {
		return false
	}
	logger.Log.Info("REDIS: Connected successfully.")
	return true
}

func getTraffic(c *gin.Context) {
	data := getRedisData()
	logger.Log.Infof("Received: %d", len(data))
	c.IndentedJSON(http.StatusOK, &data)
	getRedisData()
}

func getRedisData() []structs.InterfaceStats {
	var result []structs.InterfaceStats = make([]structs.InterfaceStats, 0)

	var cursor uint64
	var keys []string
	var err error
	keys, cursor, err = redisClient.Scan(cursor, "*", MAXPOLL).Result()
	if err != nil {
		panic(err)
	}

	// get data
	values, errGet := redisClient.MGet(keys...).Result()
	if errGet != nil {
		logger.Log.Infof("Error receiving key %s: %s", keys, errGet)
	}

	// serialize data
	for _, value := range values {
		data := structs.InterfaceStats{}
		errUnm := json.Unmarshal([]byte(value.(string)), &data)
		if errUnm != nil {
			logger.Log.Error(errUnm)
		}
		result = append(result, data)
	}

	// delete processed data
	redisClient.Del(keys...)

	return result
}

func printPrettyPost(c *gin.Context) {
	var out bytes.Buffer
	body, _ := io.ReadAll(c.Request.Body)
	err := json.Indent(&out, []byte(body), "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out.Bytes()))
}
