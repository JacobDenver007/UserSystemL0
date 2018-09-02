package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bocheninc/L0/components/crypto"
	"github.com/bocheninc/L0/core/accounts"
	"github.com/bocheninc/L0/core/types"
	"github.com/bocheninc/base/log"
	"github.com/bocheninc/base/utils"
)

var atmoicPriKeyHex = "396c663b994c3f6a8e99373c3308ee43031d7ea5120baf044168c95c45fbcf83"
var atmoicPriKey, _ = crypto.HexToECDSA("396c663b994c3f6a8e99373c3308ee43031d7ea5120baf044168c95c45fbcf83")
var txChan = make(chan *types.Transaction, 1000)
var caddr accounts.Address
var taddrs = map[accounts.Address]*crypto.PrivateKey{}
var assetID = uint32(1)

func init() {
	privateKey, _ := crypto.GenerateKey()
	addr := accounts.PublicKeyToAddress(*privateKey.Public())
	for i := 0; i < rand.Intn(10)+1; i++ {
		taddrs[addr] = privateKey
	}
}

func main() {
	asset := flag.Uint("asset", 1, "asset id")
	supply := flag.String("supply", "10000000000000", "token total supply")
	transfer := flag.Bool("transfer", true, "transfer")
	contract := flag.Bool("contract", false, "contract invoke")
	tps := flag.Uint64("tps", 0, "max tps")
	limit := flag.Int("limit", 10, "loop number")
	workers := flag.Int("workers", 1, "tps workers")
	flag.Parse()

	TCPSend([]string{"127.0.0.1:20166"})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// broadcast tx
	go func() {
		wg.Done()
		for {
			select {
			case tx, ok := <-txChan:
				if !ok {
					return
				}
				//log.Infof("Hash: %v Sender: %v Asset: %v Type: %v TxChan: %v", tx.Hash(), tx.Sender(), tx.AssetID(), tx.GetType(), len(txChan))
				if len(conn) > 0 {
					Relay(NewMsg(0x14, tx.Serialize()))
				} else {
					BroadcastTx(hex.EncodeToString(tx.Serialize()))
				}
			}
		}
	}()

	assetID = uint32(*asset)
	// Issue
	owner := accounts.PublicKeyToAddress(*atmoicPriKey.Public())
	if !utils.FileExist(fmt.Sprintf("%s-%d", owner, assetID)) {
		amount, ok := new(big.Int).SetString(*supply, 10)
		if !ok {
			panic(fmt.Errorf("failed to convert to bigint --- %s", *supply))
		}
		tx, err := IssueTx(assetID, amount, accounts.PublicKeyToAddress(*atmoicPriKey.Public()))
		if err != nil {
			panic(fmt.Errorf("issue tx failed -- %s", err))
		}
		txChan <- tx

		caddr, tx, err = DeployTx(atmoicPriKey, assetID, big.NewInt(0), "./transfer.lua", []string{})
		if err != nil {
			panic(fmt.Errorf("deploy tx failed -- %s", err))
		}
		txChan <- tx
	}
	for i := 0; i < *limit; i++ {
		TPS(*workers, *tps, *contract, *transfer)
	}

	time.Sleep(5 * time.Second)
	close(txChan)
	wg.Wait()
}

func TPS(workers int, tps uint64, contract bool, transfer bool) {
	n := uint64(0)
	cn := uint64(0)
	t := time.Now()
	defer func() {
		log.Infof("TPS eplase %s (unlimit %v, Txs %d contractTx %d, transferTx %d)", time.Now().Sub(t), tps == 0, n, cn, n-cn)
	}()
	wg := &sync.WaitGroup{}
	ticker := time.NewTicker(time.Second)
	if contract {
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				if tps > 0 {
				climitcontract:
					for {
						select {
						case <-ticker.C:
							break climitcontract
						default:
							for taddr := range taddrs {
								val := big.NewInt(rand.Int63n(1000))
								args := []string{}
								args = append(args, "transfer")
								args = append(args, taddr.String())
								args = append(args, big.NewInt(int64(assetID)).String())
								args = append(args, val.String())
								tx, _ := InvokeTx(atmoicPriKey, caddr, assetID, val, args)
								txChan <- tx
								atomic.AddUint64(&cn, 1)
								tn := atomic.AddUint64(&n, 1)
								if tn >= tps {
									break climitcontract
								}
								break
							}
						}
					}
				} else {
				cunlimitcontract:
					for {
						select {
						case <-ticker.C:
							break cunlimitcontract
						default:
							for taddr := range taddrs {
								val := big.NewInt(rand.Int63n(1000))
								args := []string{}
								args = append(args, "transfer")
								args = append(args, taddr.String())
								args = append(args, big.NewInt(int64(assetID)).String())
								args = append(args, val.String())
								tx, _ := InvokeTx(atmoicPriKey, caddr, assetID, val, args)
								txChan <- tx
								atomic.AddUint64(&cn, 1)
								atomic.AddUint64(&n, 1)
								break
							}
						}
					}
				}
			}()
		}
	}

	if transfer {
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer wg.Done()
				if tps > 0 {
				tlimitcontract:
					for {
						select {
						case <-ticker.C:
							break tlimitcontract
						default:
							for taddr := range taddrs {
								tx, _ := TransferTx(atmoicPriKey, taddr, assetID, big.NewInt(rand.Int63n(100)))
								txChan <- tx
								tn := atomic.AddUint64(&n, 1)
								if tn >= tps {
									break tlimitcontract
								}
								break
							}
						}
					}
				} else {
				tunlimitcontract:
					for {
						select {
						case <-ticker.C:
							break tunlimitcontract
						default:
							for taddr := range taddrs {
								tx, _ := TransferTx(atmoicPriKey, taddr, assetID, big.NewInt(rand.Int63n(100)))
								txChan <- tx
								atomic.AddUint64(&n, 1)
								break
							}
						}
					}
				}
			}()
		}
	}
	wg.Wait()
}
