package main

import (
	"log"
	"net/http"

	"log-client/internal"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var logPath = ""

func main() {
	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	r := gin.Default()

	stop := make(chan struct{})

	// set logPath from request body
	r.POST("/settings/log", func(c *gin.Context) {
		var json struct {
			Path string `json:"path" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logPath = json.Path

		// stop previous watcher and start new one
		close(stop)
		stop = make(chan struct{})
		go LogWriter(contract, stop)

		c.JSON(http.StatusOK, gin.H{"status": "log path set"})
	})

	// get logPath
	r.GET("/settings/log", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"path": logPath})
	})

	// read logs from chain
	r.GET("/log", func(c *gin.Context) {
		filter := c.Query("filter")

		logs, _, err := internal.ReadLogs(contract, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, logs)
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

// go routine to watch file and write new lines to ledger
func LogWriter(contract *client.Contract, stop chan struct{}) {
	internal.WatchFile(
		logPath,
		func(line string) {
			err := internal.WriteLog(contract, line, "gateway-client")
			if err != nil {
				log.Println("Failed to write log: ", err)
			} else {
				log.Println("Wrote log entry to ledger for line: ", line)
			}
		},
		stop,
	)
}
