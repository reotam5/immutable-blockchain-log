package main

import (
	"fmt"
	"os"
	"strconv"

	"log-client/internal"
)

func main() {
	clientFilter := ""
	pageSize := 10
	bookmark := ""
	if len(os.Args) >= 2 {
		clientFilter = os.Args[1]
	}

	if len(os.Args) >= 3 {
		num, err := strconv.Atoi(os.Args[2])
		if err == nil {
			pageSize = num
		}
	}

	if len(os.Args) >= 4 {
		bookmark = os.Args[3]
	}

	// get smart contract connection
	_, _, contract := internal.GetConnection()
	defer internal.CloseConnection()

	logs, hashes, new_bookmark, hasNextPage, err := internal.ReadLogsWithPagination(contract, clientFilter, pageSize, bookmark)

	if err != nil {
		panic(fmt.Errorf("failed to read logs: %w", err))
	}

	fmt.Printf("Next Bookmark: %s\n", new_bookmark)
	fmt.Printf("Has Next Page: %t\n", hasNextPage)
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
