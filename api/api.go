package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"podloxx-collector/network"

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

func InitApi() {
	go initGin()
	go initRedisCon()
}

func initGin() {
	router := gin.Default()
	router.GET("/traffic", getTraffic)
	router.GET("/traffic-ws", getTrafficWs)
	router.StaticFile("/traffic-test", "./ui/test-ws.html")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info("listen: %s\n", err)
	}
}

func initRedisCon() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println(pong)
}

func getTraffic(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, &network.TrafficData)
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

// webSocket returns json format
func getTrafficWs(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()

	for pkt := range network.ReceiverChannel {
		err = ws.WriteJSON(&pkt)
		if err != nil {
			log.Println("error write json: " + err.Error())
			return
		}
	}
}
