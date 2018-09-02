package main

import (
	"math/big"
	"time"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
	"github.com/bocheninc/L0/core/coordinate"
	"github.com/bocheninc/L0/core/types"
	"github.com/bocheninc/base/log"
)

// TransferTx transfer transaction
func TransferTx(privateKey *crypto.PrivateKey, receiver accounts.Address, assetID uint32, value *big.Int) (*types.Transaction, error) {
	sender := accounts.PublicKeyToAddress(*privateKey.Public())
	tx := types.NewTransaction(
		coordinate.HexToChainCoordinate("00"),
		coordinate.HexToChainCoordinate("00"),
		types.TypeAtomic,
		0,
		sender,
		receiver,
		assetID,
		value,
		big.NewInt(0),
		uint32(time.Now().Nanosecond()),
	)
	sig, err := privateKey.Sign(tx.SignHash().Bytes())
	if err != nil {
		log.Errorf("> TransferTx %v(sender: %v, receiver %v, assetID: %v, value: %v) --- %v", tx.Hash(), sender, receiver, assetID, value, err)
		return nil, err
	}
	tx.WithSignature(sig)

	log.Infof("> TransferTx %v(sender: %v, receiver %v, assetID: %v, value: %v)", tx.Hash(), sender, receiver, assetID, value)
	return tx, nil
}
