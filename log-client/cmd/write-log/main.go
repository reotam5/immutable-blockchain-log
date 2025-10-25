package main

import (
	"fmt"
	"os"

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
		err := internal.WriteLog(contract, line, clientName)
		if err != nil {
			panic(fmt.Errorf("failed to write log: %w", err))
		} else {
			fmt.Println("Wrote log entry to ledger for line: ", line)
		}
	}, nil)
}
