package main

import (
	"fmt"
	"os"

	"log-client/internal"
)

func main() {
	clientFilter := ""
	if len(os.Args) >= 2 {
		clientFilter = os.Args[1]
	}

	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	logs, hashes, err := internal.ReadLogs(contract, clientFilter)

	if err != nil {
		panic(fmt.Errorf("failed to read logs: %w", err))
	}

	for i, logEntry := range logs {
		valid, _ := logEntry.ValidateHash(hashes[i])

		fmt.Printf("LogID: %d\n", logEntry.ID)
		fmt.Printf("Timestamp: %s\n", logEntry.Timestamp.Format("2006-01-02T15:04:05.000000000Z07:00"))
		fmt.Printf("Content: %s\n", logEntry.Content)
		fmt.Printf("Hash: %s\n", hashes[i])
		fmt.Printf("Content Hash Valid: %t\n", valid)
		fmt.Println("-----------------------------------------------------")
	}
}
