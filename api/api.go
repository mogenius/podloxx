package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"podloxx-collector/network"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/mogenius/mo-go/logger"

	"github.com/gin-gonic/gin"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const REDISCONSTR = "127.0.0.1:6379"

var redisClient *redis.Client

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
	c.IndentedJSON(http.StatusOK, &network.TrafficData)

	iter := redisClient.Scan(0, "prefix:*", 0).Iterator()
	for iter.Next() {
		fmt.Println("keys", iter.Val())
		// data := structs.InterfaceStats{}
		// json.Unmarshal([]byte(iter.Val()), &data)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	fmt.Println("asdasd")
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
