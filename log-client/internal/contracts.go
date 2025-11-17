package internal

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type rawChain struct {
	BlobPath  string `json:"BlobPath"`
	Hash      string `json:"Hash"`
	LogID     string `json:"LogID"`
	Source    string `json:"Source"`
	Timestamp string `json:"Timestamp"`
}

type rawPaginatedResult struct {
	Records             []*rawChain `json:"records"`
	FetchedRecordsCount int32       `json:"fetchedRecordsCount"`
	Bookmark            string      `json:"bookmark"`
	HasNextPage         bool        `json:"hasNextPage"`
}

func ReadLogs(contract *client.Contract, clientFilter string) ([]LogEntry, []string, error) {
	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets", clientFilter)
	if err != nil {
		return nil, nil, err
	}

	if len(evaluateResult) == 0 {
		return nil, nil, nil
	}

	var rawChain []rawChain
	if err := json.Unmarshal([]byte(evaluateResult), &rawChain); err != nil {
		return nil, nil, err
	}

	// sort by timestamp
	sort.SliceStable(rawChain, func(i, j int) bool {
		return rawChain[i].Timestamp < rawChain[j].Timestamp
	})

	var logEntries []LogEntry
	var hashes []string

	// print out raw chain in nice format
	for _, entry := range rawChain {
		var logEntry LogEntry
		dbId, _ := strconv.Atoi(entry.BlobPath)
		_ = logEntry.LoadFromDB(uint(dbId))
		logEntries = append(logEntries, logEntry)
		hashes = append(hashes, entry.Hash)
	}

	return logEntries, hashes, nil
}

func WriteLog(contract *client.Contract, content string, clientID string) error {
	var logEntry LogEntry
	logEntry.Content = strings.TrimSpace(content)
	logEntry.Timestamp = time.Now()
	logEntry.Source = clientID
	err := logEntry.WriteToDB()
	if err != nil {
		return err
	}

	logHash, err := logEntry.Hash()
	if err != nil {
		return err
	}

	_, commit, err := contract.SubmitAsync("CreateAsset", client.WithArguments(fmt.Sprint(logEntry.ID), logHash, clientID))
	if err != nil {
		return err
	}

	if commitStatus, err := commit.Status(); err != nil {
		return err
	} else if !commitStatus.Successful {
		return fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code))
	}

	return nil
}

func ReadLogsWithPagination(contract *client.Contract, clientFilter string, pageSize int, bookmark string) ([]LogEntry, []string, string, bool, error) {
	evaluateResult, err := contract.EvaluateTransaction("GetAssetsWithFilter", clientFilter, strconv.Itoa(pageSize), bookmark)
	if err != nil {
		return nil, nil, "", false, err
	}

	var rawPaginatedResult rawPaginatedResult
	if err := json.Unmarshal([]byte(evaluateResult), &rawPaginatedResult); err != nil {
		return nil, nil, "", false, err
	}

	var logEntries []LogEntry
	var hashes []string

	// print out raw chain in nice format
	for _, entry := range rawPaginatedResult.Records {
		var logEntry LogEntry
		dbId, _ := strconv.Atoi(entry.BlobPath)
		_ = logEntry.LoadFromDB(uint(dbId))
		logEntries = append(logEntries, logEntry)
		hashes = append(hashes, entry.Hash)
	}

	return logEntries, hashes, rawPaginatedResult.Bookmark, rawPaginatedResult.HasNextPage, nil
}
