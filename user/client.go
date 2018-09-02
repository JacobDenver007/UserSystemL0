package user

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
	"github.com/bocheninc/L0/core/coordinate"
	"github.com/bocheninc/L0/core/types"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/JacobDenver007/UserSystemL0/log"
)

type RPC struct {
	rpchost     string
	rpcuser     string
	rpcpassword string
}

// type Block struct {
// 	Timestamp        big.Int             `json:"timestamp"`          // 交易时间
// 	Number           big.Int             `json:"number"`             // 区块号
// 	ParentHash       string              `json:"parent_hash"`        // 区块父哈希
// 	Hash             string              `json:"hash"`               // 区块哈希
// 	TxCnt            int                 `json:"tx_cnt"`             //交易个数
// 	Transactions     []*Transaction      `json:"transactions"`       //交易列表
// }

// type Transaction struct {
// 	Timestamp   big.Int                `json:"timestamp"`    // 交易时间
// 	BlockNumber big.Int                `json:"block_number"` // 区块号
// 	Hash        string                 `json:"hash"`         // tx id
// 	From        string                 `json:"from"`         //  发起者
// 	To          string                 `json:"to"`           // 接受者（合约地址）
// 	Value       big.Int                `json:"value"`        // eth number
// 	Gas         big.Int                `json:"gas"`          // Gas最多消耗
// 	UsedGas     big.Int                `json:"used_gas"`     // Gas消耗
// 	GasPrice    big.Int                `json:"gas_price"`    // Gas单价
// 	Nonce       big.Int                `json:"nonce"`        // Nonce
// 	Error       string                 `json:"error"`        // 错误
// }

var (
	methodGetBlockHeaderByNumber = "Ledger.GetBlockByNumber"
	methodGetBlockTxsByNumber    = "Ledger.GetTxsByBlockNumber"
	methodSendTransaction        = "Transaction.Broadcast"
)

// func (client *RPCClient) GetBlockByNumber(number uint32) (*Block, error) {
// 	t := time.Now()
// 	cnt := int64(0)
// 	defer func() {
// 		log.Debugf("GetBlockByNumber %s elpase: %s, txs: %d\n", number, time.Now().Sub(t), cnt)
// 	}()

// 	request := common.NewRPCRequest("2.0", methodGetBlockHeaderByNumber, number)

// 	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
// 	if err != nil {
// 		log.Errorf("GetBlockByNumber SendRPCRequst error --- %s --- %d", err, number.Int64())
// 		return nil, fmt.Errorf("GetBlockByNumber SendRPCRequst error --- %s", err)
// 	}

// 	if _, ok := jsonParsed.Path("error.code").Data().(float64); ok /*&& value > 0*/ {
// 		if internal == true {
// 			log.Info("getBlockNumber true rpc error --- %s --- %d", err, number.Int64())
// 			return client.getBlockByNumberForZipperone(number, false)
// 		} else {
// 			msg, _ := jsonParsed.Path("error.message").Data().(string)
// 			log.Errorf("getBlockNumber rpc error --- %s --- %d", err, number.Int64())
// 			return nil, fmt.Errorf("getBlockByNumber rpc error --- %s", msg)
// 		}

// 	}

// 	if jsonParsed.Path("result").Data() == nil {
// 		return nil, nil
// 	}

// 	return client.decodeBlockJSON(jsonParsed.Path("result"))
// }

func (client *RPC) SendTransaction(from string, to string, assetId uint32, value *big.Int) error {
	t := time.Now()
	defer func() {
		log.Debugf("SendTransaction elpase: %s\n", time.Now().Sub(t))
	}()

	param, _ := generateSendParam(from, to, assetId, value)

	request := common.NewRPCRequest("2.0", methodSendTransaction, param)

	_, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		log.Errorf("SendTransaction SendRPCRequst error --- %s", err)
		return fmt.Errorf("SendTransaction SendRPCRequst error --- %s", err)
	}

	return nil
}

func generateSendParam(from string, to string, assetID uint32, value *big.Int) (string, error) {
	account, err := DBClient.GetAccountInfo(from)
	if err != nil {
		return "", err
	}

	privateKey, _ := crypto.HexToECDSA(account.PrivateKey)

	sender := accounts.PublicKeyToAddress(*privateKey.Public())
	toByte, _ := hex.DecodeString(to)
	var receiver accounts.Address
	receiver.SetBytes(toByte)

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
		return "", err
	}
	tx.WithSignature(sig)

	log.Infof("> TransferTx %v(sender: %v, receiver %v, assetID: %v, value: %v)", tx.Hash(), sender, receiver, assetID, value)
	return hex.EncodeToString(tx.Serialize()), nil
}
