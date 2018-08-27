package user

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/JacobDenver007/UserSystemL0/log"
	"github.com/Jeffail/gabs"
	"github.com/bochentec/zipperone/zservices/coin"
)

var RPCHOST string

const (
	methodCall                   = "eth_call"
	methodGetBlockNumber         = "eth_blockNumber"
	methodGetBlockByNumber       = "eth_getBlockByNumber"
	methodGasPrice               = "eth_gasPrice"
	methodGetTransactionCount    = "eth_getTransactionCount"
	methodTxPool                 = "txpool_content"
	methodSendRawTransaction     = "eth_sendRawTransaction"
	methodGetTransactionReceipt  = "eth_getTransactionReceipt"
	methodGetBlockTracerByNumber = "debug_traceBlockByNumber"
	methodGetBalance             = "Ledger.GetBalance"
	methodGetTxsByBlockNumber    = "Ledger.GetTxsByBlockNumber"
)

// RPCClient rpc
type RPCClient struct {
	rpchost     string
	rpcuser     string
	rpcpassword string

	// 	GetRawMemPool() ([]ITransaction, error)
	// 	GetBlockNumber() (*big.Int, error)
	// 	GetBlockByNumber(number *big.Int) (IBlock, error)
	//	GetBlockByNumberJSON(number *big.Int) (string, error)
	// 	GetGasPrice() (*big.Int, error)
	// 	SendRawTransaction(signed string) (string, error)
}

// GetBlockNumber 获取最新高度
func (client *RPCClient) GetBlockNumber() (*big.Int, error) {
	request := common.NewRPCRequest("2.0", methodGetBlockNumber)
	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		return nil, fmt.Errorf("GetBlockNumber SendRPCRequst error --- %s", err)
	}

	if /*value*/ _, ok := jsonParsed.Path("error.code").Data().(float64); ok /*&& value > 0*/ {
		msg, _ := jsonParsed.Path("error.message").Data().(string)
		return nil, fmt.Errorf("GetBlockNumber error --- %s", msg)
	}

	r, ok := jsonParsed.Path("result").Data().(string)
	if !ok {
		return nil, fmt.Errorf("GetBlockNumber Path('result') interface error --- %s", err)
	}

	ret := new(big.Int)
	ret.UnmarshalJSON([]byte(r))
	return ret, nil
}

func (client *RPCClient) GetBlockByNumber(number *big.Int) (coin.IBlock, error) {
	t := time.Now()
	cnt := int64(0)
	defer func() {
		log.Debugf("GetBlockByNumber %s elpase: %s, txs: %d\n", number, time.Now().Sub(t), cnt)
	}()

	request := common.NewRPCRequest("2.0", methodGetTxsByBlockNumber, number.Int64())

	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		log.Errorf("getBlockNumber SendRPCRequst error --- %s --- %d", err, number.Int64())
		return nil, fmt.Errorf("getBlockNumber SendRPCRequst error --- %s", err)
	}

	if _, ok := jsonParsed.Path("error").Data().(string); ok {
		msg, _ := jsonParsed.Path("error.message").Data().(string)
		log.Errorf("getBlockNumber rpc error --- %s --- %d", err, number.Int64())
		return nil, fmt.Errorf("getBlockByNumber rpc error --- %s", msg)
	}

	return client.decodeBlockJSON(jsonParsed.Path("result"))
}

// SendRawTransaction 发送交易
func (client *RPCClient) SendRawTransaction(signed string) (string, error) {
	request := common.NewRPCRequest("2.0", methodSendRawTransaction, signed)
	jsonParsed, err := common.SendRPCRequst(client.rpchost, request)
	if err != nil {
		return "", fmt.Errorf("SendRawTransaction SendRPCRequst error --- %s", err)
	}

	if /*value*/ _, ok := jsonParsed.Path("error.code").Data().(float64); ok /*&& value > 0*/ {
		msg, _ := jsonParsed.Path("error.message").Data().(string)
		return "", fmt.Errorf("SendRawTransaction rpc error --- %s", msg)
	}

	r, ok := jsonParsed.Path("result").Data().(string)
	if !ok {
		return "", fmt.Errorf("SendRawTransaction result error")
	}

	return r, nil
}

func (client *RPCClient) decodeBlockJSON(jsonParsed *gabs.Container) (*Block, error) {
	// {
	//     "difficulty": "0xb130dd5b5eb21",
	//     "extraData": "0x737061726b706f6f6c2d636e2d6e6f64652d32",
	//     "gasLimit": "0x79f39e",
	//     "gasUsed": "0x79bda3",
	//     "hash": "0xfacdc77b1fedd019a871660f6e4bb86199a981faa476ad455b1df009978238f5",
	//     "logsBloom": "0x4000a42a000841002301a20408c0108836101101086124b09202848440a8402c01809000010e02013030110f09021e4240852c102110959282100001082040510700004a18b24c550406415c004011008040000b53408200a071032d0a009c1500200020829064202810d41280200e15882b0201022456e0048163515020c6875608114041b020028a4480206822020a380404922083e00100008104004220000a01c0d202424400000a0464453480c10820c2140801804100004280230a4800240c422210160007010092001400064048822415018e00b420111024a0a0e7876010c11001248001500800849b00ed1441460206c410801386601b80000007d0",
	//     "miner": "0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c",
	//     "mixHash": "0x9462b13cf76c375a656b4c4ba1bd59a21645410bfa79e96ee8ac563752a25d3a",
	//     "nonce": "0xab729340082a6470",
	//     "number": "0x57bf36",
	//     "parentHash": "0x9057876ae042dff46e4f73dfe4863b5a87ecee01e5c984cd71691badd1c5b781",
	//     "receiptsRoot": "0x7baf2205ed8a13aa424603ab6cbe654fe073e694ce6134310360947c3f8f0fff",
	//     "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	//     "size": "0x576e",
	//     "stateRoot": "0xfe2faea0c1558f0bec678071b05ce1c553be474f37873aafe70abbdf3a1059c5",
	//     "timestamp": "0x5b19d431",
	//     "totalDifficulty": "0xfb5354c3e6a2309f64",
	//     "transactions": [
	//     ],
	//     "transactionsRoot": "0x7679a1814d8f9796534c7aa32eac4e937f1d9ecd47ca1a8e9dc9c3fb706975b7",
	//     "uncles": []
	// }

	blk := NewBlock()
	blk.Header.ID = jsonParsed.Path("hash").Data().(string)
	blk.Header.PrevID = jsonParsed.Path("parentHash").Data().(string)
	ret := new(big.Int)
	ret.UnmarshalJSON([]byte(jsonParsed.Path("number").Data().(string)))
	blk.Header.Height = ret.Int64()
	ret.UnmarshalJSON([]byte(jsonParsed.Path("timestamp").Data().(string)))
	blk.Header.Time = ret.Int64()
	ret.UnmarshalJSON([]byte(jsonParsed.Path("gasUsed").Data().(string)))
	blk.Header.Size = ret.Int64()
	blk.Header.Miner = strings.ToLower(jsonParsed.Path("miner").Data().(string))
	uncles, _ := jsonParsed.Path("uncles").ArrayCount()
	if blk.Header.Height <= 4370000 {
		blk.Header.Reward = new(big.Int).Mul(big.NewInt(5), big.NewInt(1e18))
		if uncles > 0 {
			uReward := new(big.Int).Mul(big.NewInt(int64(uncles)), new(big.Int).Div(blk.Header.Reward, big.NewInt(32)))
			blk.Header.Reward = new(big.Int).Add(blk.Header.Reward, uReward)
		}
	} else {
		blk.Header.Reward = new(big.Int).Mul(big.NewInt(3), big.NewInt(1e18))
		if uncles > 0 {
			uReward := new(big.Int).Mul(big.NewInt(int64(uncles)), new(big.Int).Div(blk.Header.Reward, big.NewInt(32)))
			blk.Header.Reward = new(big.Int).Add(blk.Header.Reward, uReward)
		}
	}

	children, _ := jsonParsed.S("transactions").Children()
	for _, child := range children {
		tx, err := client.decodeTransactionJSON(child)
		if err != nil {
			return nil, err
		}
		if tx == nil {
			return nil, nil
		}
		tx.Header.Height = blk.Header.Height
		tx.Header.Time = blk.Header.Time
		blk.TTxs[tx.Header.ID] = tx
		blk.Header.Txs = append(blk.Header.Txs, tx.Header.ID)
	}
	return blk, nil
}

func (client *RPCClient) decodeTransactionJSON(jsonParsed *gabs.Container) (*Transaction, error) {
	tx := NewTransaction()
	tx.Header.ID = jsonParsed.Path("txHash").Data().(string)
	tx.Header.Time = time.Now().Unix()
	tx.Header.From = strings.ToLower(jsonParsed.Path("result.from").Data().(string))
	if jsonParsed.Path("to").Data() == nil {
		tx.Header.To = strings.ToLower(jsonParsed.Path("result.to").Data().(string))
	} else {
		tx.Header.To = "UNKOWN"
	}
	tx.Header.Value = new(big.Int)
	tx.Header.Value.UnmarshalJSON([]byte(jsonParsed.Path("result.value").Data().(string)))
	gasprice := new(big.Int)
	gasprice.UnmarshalJSON([]byte(jsonParsed.Path("gasPrice").Data().(string)))
	gasused := new(big.Int)
	//最新块中可能会出现gasUsed为null的结果，需要等待几秒再执行
	if jsonParsed.Path("gasUsed").Data() == nil {
		return nil, nil
	}
	gasused.UnmarshalJSON([]byte(jsonParsed.Path("gasUsed").Data().(string)))
	tx.Header.Size = gasused.Int64()

	tx.Header.Fee = new(big.Int).Mul(big.NewInt(tx.Header.Size), gasprice)

	logs := jsonParsed.Path("logs").Data()
	if logs != nil {
		client.getTokenTxs(jsonParsed, tx.Header)
	}

	if jsonParsed.Path("result.error").Data() != nil {
		tx.Header.Error = jsonParsed.Path("result.error").Data().(string)
		return tx, nil
	}
	callsChildren, _ := jsonParsed.S("result", "calls").Children()
	for _, callsChild := range callsChildren {
		itx := client.getInternalTxs(callsChild)
		tx.Header.InternalTransactions = append(tx.Header.InternalTransactions, itx)
	}
	return tx, nil
}

func getBalance(address string, number *big.Int) (*big.Int, error) {
	h := "pending"
	if number != nil {
		h = fmt.Sprintf("0x%s", number.Text(16))
	}
	request := common.NewRPCRequest("2.0", methodGetBalance, address, h)
	jsonParsed, err := common.SendRPCRequst(RPCHOST, request)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("getBalance SendRPCRequst error --- %s", err)
	}

	if /*value*/ _, ok := jsonParsed.Path("error.code").Data().(float64); ok /*&& value > 0*/ {
		msg, _ := jsonParsed.Path("error.message").Data().(string)
		return nil, fmt.Errorf("getBalance %s error --- %s", address, msg)
	}

	r, ok := jsonParsed.Path("result").Data().(string)
	if !ok {
		return big.NewInt(0), fmt.Errorf("getBalance Path('result') interface error --- %s", err)
	}

	var ret = big.NewInt(0)
	ret.UnmarshalJSON([]byte(r))
	return ret, nil
}
