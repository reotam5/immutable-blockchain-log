package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"log-client/internal"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type RawChain struct {
	BlobPath  string `json:"BlobPath"`
	Hash      string `json:"Hash"`
	LogID     string `json:"LogID"`
	Source    string `json:"Source"`
	Timestamp string `json:"Timestamp"`
}

func main() {
	clientFilter := ""
	if len(os.Args) >= 2 {
		clientFilter = os.Args[1]
	}

	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	readLogs(contract, clientFilter)
}

func readLogs(contract *client.Contract, clientFilter string) {
	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets", clientFilter)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	if len(evaluateResult) == 0 {
		fmt.Println("No logs found")
		return
	}

	var rawChain []RawChain
	if err := json.Unmarshal([]byte(evaluateResult), &rawChain); err != nil {
		panic(err)
	}

	// sort by timestamp
	sort.SliceStable(rawChain, func(i, j int) bool {
		return rawChain[i].Timestamp < rawChain[j].Timestamp
	})

	// print out raw chain in nice format
	for _, entry := range rawChain {
		var logEntry internal.LogEntry
		dbId, _ := strconv.Atoi(entry.BlobPath)
		_ = logEntry.LoadFromDB(uint(dbId))
		valid, _ := logEntry.ValidateHash(entry.Hash)

		fmt.Printf("LogID: %s\n", entry.LogID)
		fmt.Printf("Source: %s\n", entry.Source)
		fmt.Printf("Timestamp: %s\n", entry.Timestamp)
		fmt.Printf("Hash: %s\n", entry.Hash)
		fmt.Printf("BlobPath (DB ID): %s\n", entry.BlobPath)
		fmt.Printf("Content: %s\n", logEntry.Content)
		fmt.Printf("Content Hash Valid: %t\n", valid)
		fmt.Println("-----------------------------------------------------")
	}
}
