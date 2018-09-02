package user

// import (
// 	"context"
// 	"math/big"
// 	"reflect"
// 	"runtime"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/bochentec/zipperone/zservices/coin"
// 	"github.com/bochentec/zipperone/zservices/log"
// )

// // Scanning sync new blocks and new pending txs from main blockchain.
// func Scanning(startHeight *big.Int, endHeight *big.Int) {
// 	//最新区块
// 	var curBlock coin.IBlock
// 	//初始化 回滚
// 	for {
// 		//获取本地最新区块
// 		block, err := DBClient.GetBestBlock()
// 		if err != nil {
// 			log.Errorf("[Scanning] GetBestBlock --- %s", err)
// 			log.Panic(err)
// 		}
// 		curBlock = block

// 		//本地没有区块可回滚
// 		if reflect.ValueOf(curBlock).IsNil() {
// 			break
// 		}

// 		//获取节点最新区块
// 		block, err = rpcClient.GetBlockByNumber(curBlock.Number())
// 		if err != nil {
// 			log.Errorf("[Scanning] GetBlockByNumber %s --- %s", curBlock.Number(), err)
// 			log.Panic(err)
// 		}

// 		//区块相同，无需回滚
// 		if strings.Compare(block.Hash(), curBlock.Hash()) == 0 {
// 			break
// 		}

// 		//区块回滚
// 		log.Infof("[Scanning] RollBack Block: height: %s hash: %s", curBlock.Number(), curBlock.Hash())
// 		if err := coinInst.DeleteBlock(curBlock); err != nil {
// 			log.Errorf("[Scanning] DeleteBlock %s --- %s", curBlock.Number(), err)
// 			log.Panic(err)
// 		}
// 	}

// 	//开始高度 取最大
// 	fromNumber := big.NewInt(0)
// 	if !reflect.ValueOf(curBlock).IsNil() {
// 		fromNumber = new(big.Int).Add(curBlock.Number(), big.NewInt(1))
// 	}
// 	if startHeight != nil && startHeight.Cmp(fromNumber) > 0 {
// 		fromNumber = startHeight
// 	}

// 	//结束高度 取最小
// 	toNumber, err := coinInst.GetBlockNumber()
// 	if err != nil {
// 		log.Errorf("[Scanning] GetBlockNumber--- %s", err)
// 		log.Panic(err)
// 	}
// 	if endHeight != nil && toNumber.Cmp(endHeight) > 0 {
// 		toNumber = endHeight
// 	}

// 	//cpu个数
// 	cpus := runtime.NumCPU()
// 	if cpus > 10 {
// 		cpus = 10
// 	}
// 	if !bcthread {
// 		cpus = 1
// 	}
// 	log.Infof("[Scanning] FromNumber:%d ===> ToNumber:%d cpus %d", fromNumber, toNumber, cpus)
// 	//间距
// 	DValue := new(big.Int).Sub(toNumber, fromNumber)
// 	if DValue.Sign() < 0 {
// 		log.Errorf("[Scanning] FromNumber:%d ===> ToNumber:%d cpus %d", fromNumber, toNumber, cpus)
// 		return
// 	}

// 	DValue = new(big.Int).Sub(toNumber, fromNumber)
// 	safetyConfirmations := big.NewInt(10)
// 	if DValue.Cmp(safetyConfirmations) > 0 {
// 		//差距大于10个块，启用流水线
// 		startNumber := big.NewInt(fromNumber.Int64())
// 		chParse := make(chan *big.Int, 10)
// 		chSave := make(chan coin.IBlock, 100)

// 		wg := &sync.WaitGroup{}
// 		cctx, cancel := context.WithCancel(ctx)
// 		//rpc goroutine
// 		wg.Add(1)
// 		go func(ctx context.Context, chParse chan *big.Int) {
// 			defer wg.Done()
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				default:
// 					if new(big.Int).Add(fromNumber, safetyConfirmations).Cmp(toNumber) <= 0 {
// 						chParse <- fromNumber
// 						fromNumber = new(big.Int).Add(fromNumber, big.NewInt(1))
// 					}
// 				}
// 			}
// 		}(cctx, chParse)

// 		//rpc goroutine
// 		wg.Add(1)
// 		go func(ctx context.Context, chParse chan *big.Int, chSave chan coin.IBlock) {
// 			defer wg.Done()
// 			defer close(chParse)
// 			if cpus == 1 {
// 				for {
// 					select {
// 					case <-ctx.Done():
// 						return
// 					case number := <-chParse:
// 						for {
// 							block, err := coinInst.GetBlockByNumber(number)
// 							if err != nil {
// 								log.Errorf("[Scanning] GetBlockByNumber %s --- %s", number, err)
// 								continue
// 							}
// 							chSave <- block
// 							break
// 						}
// 					}
// 				}
// 			} else {
// 				workers := cpus
// 				chQueue := make([]chan coin.IBlock, workers)
// 				for i := 0; i < workers; i++ {
// 					chQueue[i] = make(chan coin.IBlock)
// 				}
// 				defer func() {
// 					for _, ch := range chQueue {
// 						close(ch)
// 					}
// 				}()
// 				twg := &sync.WaitGroup{}
// 				twg.Add(1)
// 				go func(chQueue []chan coin.IBlock) {
// 					defer twg.Done()
// 					i := 0
// 					ch := chQueue[i]
// 					for {
// 						select {
// 						case <-ctx.Done():
// 							return
// 						case ret := <-ch:
// 							chSave <- ret
// 							i++
// 							if i == len(chQueue) {
// 								i = 0
// 							}
// 							ch = chQueue[i]
// 						}
// 					}
// 				}(chQueue)

// 				twg.Add(workers)
// 				for i := 0; i < workers; i++ {
// 					go func(chParse chan *big.Int, chQueue []chan coin.IBlock) {
// 						defer twg.Done()
// 						for {
// 							select {
// 							case <-ctx.Done():
// 								return
// 							case number := <-chParse:
// 								for {
// 									block, err := coinInst.GetBlockByNumber(number)
// 									if err != nil {
// 										log.Errorf("[Scanning] GetBlockByNumber %s --- %s", number, err)
// 										continue
// 									}
// 									index := new(big.Int).Sub(number, startNumber).Int64() % int64(workers)
// 									select {
// 									case chQueue[index] <- block:
// 										break
// 									case <-ctx.Done():
// 										break
// 									}
// 									break
// 								}

// 							}
// 						}
// 					}(chParse, chQueue)
// 				}

// 				twg.Wait()
// 			}
// 		}(cctx, chParse, chSave)

// 		//db goroutine
// 		wg.Add(1)
// 		go func(ctx context.Context, chSave chan coin.IBlock, cancel context.CancelFunc) {
// 			defer wg.Done()
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case block := <-chSave:
// 					curBlock = block
// 					if err := coinInst.InsertBlock(curBlock); err != nil {
// 						log.Errorf("[Scanning] InsertBlock >= --- %s", err)
// 						continue
// 					}
// 					if new(big.Int).Add(curBlock.Number(), safetyConfirmations).Cmp(toNumber) == 0 {
// 						log.Infof("[Scanning] safetyConfirmations Cancel")
// 						cancel()
// 					}
// 				}
// 			}
// 		}(cctx, chSave, cancel)

// 		wg.Wait()
// 		close(chSave)
// 		fromNumber = new(big.Int).Add(curBlock.Number(), big.NewInt(1))
// 	}

// 	pendingTxs := map[string]coin.ITransaction{}
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			break
// 		default:
// 		}

// 		if endHeight != nil && new(big.Int).Sub(fromNumber, endHeight).Cmp(big.NewInt(10)) > 0 {
// 			return
// 		}

// 		block, err := coinInst.GetBlockByNumber(fromNumber)
// 		if err != nil {
// 			log.Errorf("[Scanning] GetBlockByNumber %s--- %s", fromNumber, err)
// 			continue
// 		}
// 		if block == nil || reflect.ValueOf(block).IsNil() {
// 			txs, err := coinInst.GetRawMemPool(pendingTxs)
// 			if err != nil {
// 				log.Errorf("[Scanning] GetRawMemPool --- %s", err)
// 				continue
// 			}
// 			if err := coinInst.InsertPendingTxs(txs); err != nil {
// 				log.Errorf("[Scanning] InsertPendingTxs --- %s", err)
// 				continue
// 			}
// 			pendingTxs = map[string]coin.ITransaction{}
// 			for _, tx := range txs {
// 				pendingTxs[tx.TxHash()] = tx
// 			}

// 			time.Sleep(time.Duration(pendingduration) * time.Second)
// 			continue
// 		}

// 		if !reflect.ValueOf(curBlock).IsNil() && strings.Compare(curBlock.Hash(), block.ParentHash()) != 0 {
// 			log.Warnf("[Scanning] RollBack Block: height: %s hash: %s", curBlock.Number(), curBlock.Hash())
// 			if err := coinInst.DeleteBlock(curBlock); err != nil {
// 				log.Errorf("[Scanning] DeleteBlock %s--- %s", curBlock.Number(), err)
// 				continue
// 			}

// 			block, err := coinInst.GetBestBlock()
// 			if err != nil {
// 				log.Errorf("[Scanning] GetBestBlock --- %s", err)
// 				continue
// 			}
// 			curBlock = block
// 		} else {
// 			curBlock = block
// 			if err := coinInst.InsertBlock(curBlock); err != nil {
// 				log.Errorf("[Scanning] InsertBlock %s --- %s", curBlock.Number(), err)
// 				continue
// 			}
// 		}
// 		if !reflect.ValueOf(curBlock).IsNil() {
// 			fromNumber = new(big.Int).Add(curBlock.Number(), big.NewInt(1))
// 		} else {
// 			fromNumber = big.NewInt(0)
// 		}
// 	}
// }
