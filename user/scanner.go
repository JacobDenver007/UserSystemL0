package user

import (
	"context"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/JacobDenver007/UserSystemL0/log"
)

// Scanning sync new blocks and new pending txs from main blockchain.
func Scanning(startHeight *big.Int, endHeight *big.Int) {
	ctx := context.Background()

	//最新区块
	curBlock := &Block{}

	//获取本地最新区块
	block, err := DBClient.GetBestBlock()
	if err != nil {
		log.Errorf("[Scanning] GetBestBlock --- %s", err)
		log.Panic(err)
	}
	curBlock = block

	//开始高度 取最大
	fromNumber := uint32(0)
	if !reflect.ValueOf(curBlock).IsNil() {
		fromNumber = curBlock.Header.Height + 1
	}
	if startHeight != nil {
		fromNumber = uint32(startHeight.Uint64())
	}

	//结束高度 取最小
	toNumber, err := RPCClient.GetBlockNumber()
	if err != nil {
		log.Errorf("[Scanning] GetBlockNumber--- %s", err)
		log.Panic(err)
	}
	if endHeight != nil {
		toNumber = uint32(startHeight.Uint64())
	}

	//cpu个数
	cpus := 1
	log.Infof("[Scanning] FromNumber:%d ===> ToNumber:%d cpus %d", fromNumber, toNumber, cpus)
	//间距
	DValue := toNumber - fromNumber
	if DValue < 0 {
		log.Errorf("[Scanning] FromNumber:%d ===> ToNumber:%d cpus %d", fromNumber, toNumber, cpus)
		return
	}

	if DValue > 10 {
		//差距大于10个块，启用流水线
		chParse := make(chan uint32, 10)
		chSave := make(chan *Block, 100)

		wg := &sync.WaitGroup{}
		cctx, cancel := context.WithCancel(ctx)
		//rpc goroutine
		wg.Add(1)
		go func(ctx context.Context, chParse chan uint32) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if fromNumber < toNumber {
						chParse <- fromNumber
						fromNumber++
					}
				}
			}
		}(cctx, chParse)

		//rpc goroutine
		wg.Add(1)
		go func(ctx context.Context, chParse chan uint32, chSave chan *Block) {
			defer wg.Done()
			defer close(chParse)

			for {
				select {
				case <-ctx.Done():
					return
				case number := <-chParse:
					for {
						block, err := RPCClient.GetBlockByNumber(number)
						if err != nil {
							log.Errorf("[Scanning] GetBlockByNumber %d --- %s", number, err)
							continue
						}
						chSave <- block
						break
					}
				}
			}

		}(cctx, chParse, chSave)

		//db goroutine
		wg.Add(1)
		go func(ctx context.Context, chSave chan *Block, cancel context.CancelFunc) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case block := <-chSave:
					curBlock = block
					if err := DBClient.InsertBlock(curBlock); err != nil {
						log.Errorf("[Scanning] InsertBlock >= --- %s", err)
						continue
					}
					if fromNumber == toNumber {
						log.Infof("[Scanning] safetyConfirmations Cancel")
						cancel()
					}
				}
			}
		}(cctx, chSave, cancel)

		wg.Wait()
		close(chSave)
		fromNumber++
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		block, err := RPCClient.GetBlockByNumber(fromNumber)
		if err != nil {
			log.Errorf("[Scanning] GetBlockByNumber %d--- %s", fromNumber, err)
			time.Sleep(time.Duration(1))
			continue
		}
		if block == nil || reflect.ValueOf(block).IsNil() {
			time.Sleep(time.Duration(1))
			continue
		}

		curBlock = block
		if err := DBClient.InsertBlock(curBlock); err != nil {
			log.Errorf("[Scanning] InsertBlock %d --- %s", curBlock.Header.Height, err)
			continue
		}
		if !reflect.ValueOf(curBlock).IsNil() {
			fromNumber++
		} else {
			fromNumber = 0
		}
	}
}
