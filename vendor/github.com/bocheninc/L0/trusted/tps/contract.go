package main

import (
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/components/utils"
	"github.com/bocheninc/L0/core/accounts"
	"github.com/bocheninc/L0/core/coordinate"
	"github.com/bocheninc/L0/core/types"
	"github.com/bocheninc/base/log"
)

// DeployTx deploy contract
func DeployTx(privkey *crypto.PrivateKey, assetID uint32, value *big.Int, path string, args []string) (accounts.Address, *types.Transaction, error) {
	sender := accounts.PublicKeyToAddress(*privkey.Public())
	contractSpec := new(types.ContractSpec)
	f, _ := os.Open(path)
	buf, _ := ioutil.ReadAll(f)
	contractSpec.ContractCode = buf
	var a accounts.Address
	pubBytes := []byte(sender.String() + string(buf))
	a.SetBytes(crypto.Keccak256(pubBytes[1:])[12:])
	contractSpec.ContractAddr = a.Bytes()
	contractSpec.ContractParams = args

	tx := types.NewTransaction(
		coordinate.HexToChainCoordinate("00"),
		coordinate.HexToChainCoordinate("00"),
		types.TypeLuaContractInit,
		uint32(0),
		sender,
		accounts.NewAddress(contractSpec.ContractAddr),
		assetID,
		value,
		big.NewInt(0),
		uint32(time.Now().Unix()),
	)
	tx.Payload = utils.Serialize(contractSpec)
	sig, err := privkey.Sign(tx.SignHash().Bytes())
	if err != nil {
		log.Errorf("> DeployTx %v(cAddress %v, assetID: %v, value: %v) --- %v", tx.Hash(), accounts.NewAddress(contractSpec.ContractAddr), assetID, value, err)
		return accounts.Address{}, nil, err
	}
	tx.WithSignature(sig)
	log.Infof("> DeployTx %v(cAddress %v, assetID: %v, value: %v)", tx.Hash(), accounts.NewAddress(contractSpec.ContractAddr), assetID, value)
	return accounts.NewAddress(contractSpec.ContractAddr), tx, nil
}

// InvokeTx invoke transaction
func InvokeTx(privkey *crypto.PrivateKey, contractAddr accounts.Address, assetID uint32, value *big.Int, args []string) (*types.Transaction, error) {
	sender := accounts.PublicKeyToAddress(*privkey.Public())
	contractSpec := new(types.ContractSpec)
	contractSpec.ContractAddr = contractAddr.Bytes()
	contractSpec.ContractParams = args
	tx := types.NewTransaction(
		coordinate.HexToChainCoordinate("00"),
		coordinate.HexToChainCoordinate("00"),
		types.TypeContractInvoke,
		uint32(0),
		sender,
		accounts.NewAddress(contractSpec.ContractAddr),
		assetID,
		value,
		big.NewInt(0),
		uint32(time.Now().Unix()),
	)
	tx.Payload = utils.Serialize(contractSpec)
	sig, err := privkey.Sign(tx.SignHash().Bytes())
	if err != nil {
		log.Errorf("> InvokeTx %v(cAddress %v, assetID: %v, value: %v) --- %v", tx.Hash(), accounts.NewAddress(contractSpec.ContractAddr), assetID, value, err)
		return nil, err
	}
	tx.WithSignature(sig)
	log.Infof("> InvokeTx %v(cAddress %v, assetID: %v, value: %v)", tx.Hash(), accounts.NewAddress(contractSpec.ContractAddr), assetID, value)
	return tx, nil
}
