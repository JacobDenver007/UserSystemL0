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

type Block struct {
	Header *BlockHeader
	Txs    []*Transaction
}

type BlockHeader struct {
	PreviousHash string `json:"previousHash" `
	TimeStamp    uint32 `json:"timeStamp"`
	Nonce        uint32 `json:"nonce" `
	Height       uint32 `json:"height" `
}

type Transaction struct {
	Data txdata `json:"data"`

	Hash string `json:"hash"`
}

type txdata struct {
	Type       uint32   `json:"type"`
	Nonce      uint32   `json:"nonce"`
	Sender     string   `json:"sender"`
	Recipient  string   `json:"recipient"`
	AssetID    uint32   `json:"assetid"`
	Amount     *big.Int `json:"amount"`
	Fee        *big.Int `json:"fee"`
	CreateTime uint32   `json:"createTime"`
}

var (
	methodGetBlockHeaderByNumber = "Ledger.GetBlockByNumber"
	methodGetBlockTxsByNumber    = "Ledger.GetTxsByBlockNumber"
	methodSendTransaction        = "Transaction.Broadcast"
)

func GetBlockByNumber(number uint32) (*Block, error) {
	blockHeader, err := RPCClient.GetBlockHeaderByNumber(number)
	if err != nil {
		return nil, err
	}
	blockTxs, err := RPCClient.GetBlockTxsByNumber(number)
	if err != nil {
		return nil, err
	}

	block := &Block{Header: blockHeader, Txs: blockTxs}
	return block, nil
}

func (client *RPC) GetBlockHeaderByNumber(number uint32) (*BlockHeader, error) {
	t := time.Now()
	defer func() {
		log.Debugf("GetBlockHeaderByNumber %s elpase: %s\n", number, time.Now().Sub(t))
	}()

	request := common.NewRPCRequest("2.0", methodGetBlockHeaderByNumber, number)

	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		log.Errorf("GetBlockHeaderByNumber SendRPCRequst error --- %s --- %d", err, number)
		return nil, fmt.Errorf("GetBlockHeaderByNumber SendRPCRequst error --- %s", err)
	}

	blockHeader := &BlockHeader{}
	blockHeader.PreviousHash = jsonParsed.Path("previousHash").Data().(string)
	blockHeader.TimeStamp = uint32(jsonParsed.Path("timeStamp").Data().(float64))
	blockHeader.Height = uint32(jsonParsed.Path("height").Data().(float64))
	blockHeader.Nonce = uint32(jsonParsed.Path("nonce").Data().(float64))

	return blockHeader, nil
}

func (client *RPC) GetBlockTxsByNumber(number uint32) ([]*Transaction, error) {
	t := time.Now()
	cnt := int64(0)
	defer func() {
		log.Debugf("GetBlockTxsByNumber %s elpase: %s, txs: %d\n", number, time.Now().Sub(t), cnt)
	}()

	request := common.NewRPCRequest("2.0", methodGetBlockTxsByNumber, number)

	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		log.Errorf("GetBlockTxsByNumber SendRPCRequst error --- %s --- %d", err, number)
		return nil, fmt.Errorf("GetBlockTxsByNumber SendRPCRequst error --- %s", err)
	}

	txs := make([]*Transaction, 0)

	children, _ := jsonParsed.S("transactions").Children()
	for _, child := range children {
		tx := &Transaction{}
		tx.Data.Sender = child.Path("data.sender").Data().(string)
		tx.Data.Recipient = child.Path("data.recipient").Data().(string)
		tx.Data.Amount = new(big.Int)
		tx.Data.Amount.UnmarshalJSON([]byte(jsonParsed.Path("data.amount").Data().(string)))
		tx.Data.Fee = new(big.Int)
		tx.Data.Fee.UnmarshalJSON([]byte(jsonParsed.Path("data.fee").Data().(string)))
		tx.Data.AssetID = uint32(child.Path("data.assetid").Data().(float64))
		tx.Data.Type = uint32(child.Path("data.type").Data().(float64))
		tx.Data.CreateTime = uint32(child.Path("data.createTime").Data().(float64))
		tx.Hash = child.Path("hash").Data().(string)

		txs = append(txs, tx)
	}

	return txs, nil
}

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
