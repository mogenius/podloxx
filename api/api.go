package api

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"podloxx-collector/structs"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/mogenius/mo-go/logger"

	"github.com/gin-gonic/gin"
)

const MAXPOLL int64 = 100

var HtmlDirFs embed.FS

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
	router.StaticFS("/podloxx", embedFs())
	//router.Static("/podloxx", os.Getenv("PODLOXX_DIST"))
	router.GET("/traffic", getTraffic)

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
	c.IndentedJSON(http.StatusOK, &data)
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

	// return if no new data can be gathered
	if len(keys) == 0 {
		return result
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

func embedFs() http.FileSystem {
	sub, err := fs.Sub(HtmlDirFs, "ui/dist/podloxx")

	dirContent, err := getAllFilenames(&HtmlDirFs, "")
	if err != nil {
		panic(err)
	}

	if len(dirContent) <= 0 {
		panic("dist folder empty. Cannnot serve site. FATAL.")
	} else {
		logger.Log.Noticef("Loaded %d static files from embed.", len(dirContent))
	}
	return http.FS(sub)
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

func getAllFilenames(fs *embed.FS, dir string) (out []string, err error) {
	if len(dir) == 0 {
		dir = "."
	}

	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fp := path.Join(dir, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}

			out = append(out, res...)

			continue
		}

		out = append(out, fp)
	}

	return
}
