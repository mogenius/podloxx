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
	"github.com/mogenius/mo-go/utils"

	"github.com/gin-contrib/cors"
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
var uptime = time.Now()
var lastOverview structs.Overview = structs.Overview{}

func TestRedis() {
	for {
		if initRedis() {
			break
		}
		logger.Log.Infof("Waiting for redis to come alive: %s", REDISCONSTR)
		time.Sleep(1 * time.Second)
	}
	data := getRedisDataTotal()
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
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.StaticFS("/podloxx", embedFs())
	router.GET("/traffic/total", getTrafficTotal)
	router.GET("/traffic/flow", getTrafficFlow)
	router.GET("/traffic/overview", getTrafficOverview)

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

func getTrafficFlow(c *gin.Context) {
	data := getRedisDataFlow()
	c.IndentedJSON(http.StatusOK, &data)
}

func getTrafficTotal(c *gin.Context) {
	data := getRedisDataTotal()
	c.IndentedJSON(http.StatusOK, &data)
}

func getTrafficOverview(c *gin.Context) {
	data := getOverview()
	lastOverview = data
	c.IndentedJSON(http.StatusOK, &data)
}

func getRedisDataFlow() map[string]structs.InterfaceStatsNumbers {
	var result map[string]structs.InterfaceStatsNumbers = make(map[string]structs.InterfaceStatsNumbers)

	var cursor uint64
	var keys []string
	var err error
	keys, _, err = redisClient.Scan(cursor, "traffic_*", MAXPOLL).Result()
	if err != nil {
		logger.Log.Error(err)
		return result
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
		result[data.PodName] = structs.Minify(data)
	}

	// delete processed data
	redisClient.Del(keys...)

	return result
}

func getRedisDataTotal() map[string]structs.InterfaceStats {
	var result map[string]structs.InterfaceStats = make(map[string]structs.InterfaceStats)

	var cursor uint64
	var keys []string
	var err error
	keys, _, err = redisClient.Scan(cursor, "pod_*", MAXPOLL).Result()
	if err != nil {
		logger.Log.Error(err)
		return result
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
		result[data.PodName] = data
	}

	return result
}

func getOverview() structs.Overview {
	result := structs.Overview{ExternalBandwidthPerSec: "0 B", InternalBandwidthPerSec: "0 B"}
	data := getRedisDataTotal()
	for _, entry := range data {
		mini := structs.Minify(entry)
		result.PacketsSum += mini.PacketsSum
		result.TransmitBytes += mini.TransmitBytes
		result.ReceivedBytes += mini.ReceivedBytes
		result.UnknownBytes += mini.UnknownBytes
		result.LocalReceivedBytes += mini.LocalReceivedBytes
		result.LocalTransmitBytes += mini.LocalTransmitBytes
		result.TotalPods += 1
	}

	// nodes count
	nodes := make(map[string]bool)
	for _, value := range data {
		nodes[value.Node] = true
	}
	result.TotalNodes = len(nodes)

	// namespaces count
	namespaces := make(map[string]bool)
	for _, value := range data {
		namespaces[value.Namespace] = true
	}
	result.TotalNamespaces = len(namespaces)

	result.Uptime = uptime.Format(time.RFC3339)

	result.LastUpdate = time.Now().Unix()

	// seconds since lastOverview.LastUpdate
	seconds := result.LastUpdate - lastOverview.LastUpdate
	if seconds > 0 {
		externalTraffic := int((result.ReceivedBytes + result.TransmitBytes) - (result.LocalReceivedBytes + result.LocalTransmitBytes))
		externalTrafficLast := int((lastOverview.ReceivedBytes + lastOverview.TransmitBytes) - (lastOverview.LocalReceivedBytes + lastOverview.LocalTransmitBytes))
		result.ExternalBandwidthPerSec = utils.BytesToHumanReadable(uint64(externalTraffic - externalTrafficLast/int(seconds)))
		result.InternalBandwidthPerSec = utils.BytesToHumanReadable(uint64(int(result.LocalReceivedBytes+result.LocalTransmitBytes) - (int(lastOverview.LocalReceivedBytes+lastOverview.LocalTransmitBytes))/int(seconds)))
		result.PacketsPerSec = (int(result.PacketsSum) - int(lastOverview.PacketsSum)) / int(seconds)
	} else {
		result.InternalBandwidthPerSec = lastOverview.InternalBandwidthPerSec
		result.ExternalBandwidthPerSec = lastOverview.ExternalBandwidthPerSec
	}

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
