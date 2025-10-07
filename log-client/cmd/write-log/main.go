package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log-client/internal"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/write-log/main.go <filename> <client-name>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	clientName := os.Args[2]

	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	// everytime a new line is added to the file, create a new asset on the ledger
	internal.WatchFile(filePath, func(line string) {
		createAsset(contract, line, clientName)
	})
}

// this is the core of the application - uploads to off-chain storage and create a new asset on the ledger
func createAsset(contract *client.Contract, content string, clientID string) {
	var logEntry internal.LogEntry
	logEntry.Content = strings.TrimSpace(content)
	logEntry.Timestamp = time.Now()
	err := logEntry.WriteToDB()
	if err != nil {
		panic(fmt.Errorf("failed to write log entry to database: %w", err))
	}

	logHash, err := logEntry.Hash()
	if err != nil {
		panic(fmt.Errorf("failed to hash log entry: %w", err))
	}

	_, commit, err := contract.SubmitAsync("CreateAsset", client.WithArguments(fmt.Sprint(logEntry.ID), logHash, clientID))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Println("Transaction committed successfully: ", logEntry.String())
}
