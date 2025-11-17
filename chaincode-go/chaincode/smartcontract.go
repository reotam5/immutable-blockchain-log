package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	BlobPath  string `json:"BlobPath"`
	Hash      string `json:"Hash"`
	LogID     string `json:"LogID"`
	Source    string `json:"Source"`
	Timestamp string `json:"Timestamp"`
}

type PaginatedQueryResult struct {
	Records             []*Asset `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
	HasNextPage         bool     `json:"hasNextPage"`
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, blobPath string, hash string, source string) error {
	txTime, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return err
	}

	// create unique key based on timestamp and uuid (using timestamp to help with ordering)
	ts := time.Unix(txTime.Seconds, int64(txTime.Nanos)).UTC()
	timestampStr := ts.Format("20060102T150405Z")
	id := uuid.New().String()
	key := fmt.Sprintf("asset:%s:%s", timestampStr, id)

	exists, err := assetExists(ctx, key)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", key)
	}

	asset := Asset{
		LogID:     key,
		BlobPath:  blobPath,
		Hash:      hash,
		Source:    source,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, assetJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func assetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface, sourceFilter string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		if sourceFilter == "" || asset.Source == sourceFilter {
			assets = append(assets, &asset)
		}
	}

	return assets, nil
}

func (s *SmartContract) GetAssetsWithFilter(ctx contractapi.TransactionContextInterface, source string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {
	query := fmt.Sprintf(`{
			"selector": {
					"Source": "%s"
			}
		}`, source)
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(query, int32(pageSize), bookmark)
	if err != nil {
		return nil, nil
	}
	defer resultsIterator.Close()

	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, nil
	}

	hasNextPage := false
	if responseMetadata.Bookmark != "nil" {
		_, nextPageResponseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(query, int32(1), responseMetadata.Bookmark)
		if err != nil {
			return nil, nil
		}
		if nextPageResponseMetadata.FetchedRecordsCount > 0 {
			hasNextPage = true
		}
	}

	if assets == nil {
		return &PaginatedQueryResult{
			Records:             []*Asset{},
			FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
			Bookmark:            responseMetadata.Bookmark,
			HasNextPage:         hasNextPage,
		}, nil
	} else {
		return &PaginatedQueryResult{
			Records:             assets,
			FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
			Bookmark:            responseMetadata.Bookmark,
			HasNextPage:         hasNextPage,
		}, nil
	}
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
