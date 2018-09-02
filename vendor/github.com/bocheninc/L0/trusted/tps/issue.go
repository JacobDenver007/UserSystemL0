package main

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
	"github.com/bocheninc/L0/core/coordinate"
	"github.com/bocheninc/L0/core/types"
	"github.com/bocheninc/base/log"
)

// IssueTx issue asset
func IssueTx(assetID uint32, supply *big.Int, owner accounts.Address) (*types.Transaction, error) {
	issuePriKeyHex := "496c663b994c3f6a8e99373c3308ee43031d7ea5120baf044168c95c45fbcf83"
	privateKey, _ := crypto.HexToECDSA(issuePriKeyHex)
	sender := accounts.PublicKeyToAddress(*privateKey.Public())

	tx := types.NewTransaction(
		coordinate.HexToChainCoordinate("00"),
		coordinate.HexToChainCoordinate("00"),
		types.TypeIssue,
		0,
		sender,
		owner,
		assetID,
		supply,
		big.NewInt(0),
		uint32(time.Now().Nanosecond()),
	)
	issueCoin := make(map[string]interface{})
	issueCoin["id"] = assetID
	tx.Payload, _ = json.Marshal(issueCoin)

	sig, err := privateKey.Sign(tx.SignHash().Bytes())
	if err != nil {
		log.Errorf("> IssueTx %v(assetID: %v, supply: %v owner: %v) --- %v", tx.Hash(), assetID, supply, owner, err)
		return nil, err
	}
	tx.WithSignature(sig)

	log.Infof("> IssueTx %v(assetID: %v, supply: %v owner: %v)", tx.Hash(), assetID, supply, owner)
	return tx, err
}
