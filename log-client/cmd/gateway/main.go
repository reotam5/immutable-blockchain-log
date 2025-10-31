package main

import (
	"log"
	"net/http"
	"os"

	"log-client/internal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var logPath = ""

func main() {
	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

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
		// check if path is valid file
		_, err := os.Open(json.Path)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file path"})
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

		logs, hashes, err := internal.ReadLogs(contract, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		detailedLogs := []internal.DetailedLogEntry{}
		for i := range logs {
			logEntry := &logs[i]
			detaildLogEntry, _ := logEntry.GetDetailedLogEntry(hashes[i])
			detailedLogs = append(detailedLogs, *detaildLogEntry)
		}

		c.JSON(http.StatusOK, detailedLogs)
	})

	log.Println("Server starting on :" + internal.PORT)
	r.Run(":" + internal.PORT)
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
