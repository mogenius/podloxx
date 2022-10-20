package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"podloxx-collector/network"

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
	router := gin.Default()
	router.GET("/traffic", getTraffic)
	router.GET("/traffic-ws", getTrafficWs)
	router.StaticFile("/traffic-test", "./ui/test-ws.html")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		Handler: router,
	}
	fmt.Println(srv.Addr)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Warning("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Warning("Server forced to shutdown:", err)
	}

	logger.Log.Warning("Server exiting")
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
