package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

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

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, blobPath string, hash string, source string) error {
	logId := ctx.GetStub().GetTxID()
	exists, err := s.AssetExists(ctx, logId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", logId)
	}

	asset := Asset{
		LogID:     logId,
		BlobPath:  blobPath,
		Hash:      hash,
		Source:    source,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(logId, assetJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
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
