package evm

import (
	"chaintx/chains"
	"chaintx/reporter"
	"chaintx/store"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

func sleepWithJitter(seconds int, reason string) {
	//log.Println("Sleeping", seconds, "seconds -", reason)

	jitter := 1 + rand.Intn(3)
	duration := time.Duration(seconds + jitter)
	time.Sleep(duration * time.Second)
}

// ChaseTransfers
// TODO: Move out SlidingWindow logic
func ChaseTransfers(
	chain chains.ChainName,
	address string,
	fromBlock uint64,
	maxBlocks uint64,
	transfersChan chan reporter.Transfer,
) {
	defer close(transfersChan)

	client, err := GetClient(chain)
	if err != nil {
		log.Println("Error encountered fetching client:", err)
		return
	}

	endBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Println("Error encountered fetching block number:", err)
		return
	}

	toBlock := fromBlock + min(endBlock-fromBlock, maxBlocks)

	for {

		// We've reached the tip of the chain. Sleep until another block is produced.
		if fromBlock >= endBlock {
			sleepWithJitter(1, "Reached tip of chain.")

			newTipHeight, err := client.BlockNumber(context.Background())
			if err != nil {
				// Might be an intermittent network issue - try again next iter
				log.Println("Error encountered fetching block number:", err)
				continue
			}

			if newTipHeight >= endBlock {
				endBlock = newTipHeight
				toBlock = endBlock
			}

			continue
		}

		// Chase the chain
		for toBlock <= endBlock {
			transfers := transfers(
				client,
				chain,
				address,
				fromBlock,
				toBlock,
			)

			log.Println("Block", fromBlock, "to", toBlock, ":", len(transfers), "transfers")

			for _, transfer := range transfers {
				transfersChan <- transfer
			}

			// Throttle
			sleepWithJitter(1, "Throttling")

			// Check if we have a new tip
			newTipHeight, err := client.BlockNumber(context.Background())
			if err != nil {
				// We won't break here, since it could be an intermittent network issue.
				// We'll just leave toBlock as-is.
				fmt.Println("Error encountered fetching block number:", err)
			} else {
				endBlock = newTipHeight
			}

			// Move the window forward, limited by `maxBlocks` to reduce load on the node provider
			// E.g.:
			//   endBlock  = 100
			//   fromBlock = 50 --- we're 50 blocks away from the chain tip
			//   maxBlocks = 10 --- we only move the window 10 blocks forward, instead of 50
			fromBlock = toBlock + 1
			toBlock = fromBlock + min(endBlock-fromBlock, maxBlocks)
		}
	}
}

func transfers(
	client *ethclient.Client,
	chain chains.ChainName,
	address string,
	fromBlock uint64,
	toBlock uint64,
) []reporter.Transfer {
	//log.Println("\tLooking from", fromBlock, "to", toBlock)

	var transfers []reporter.Transfer

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	fromAddressHash := common.HexToHash(address)

	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Topics: [][]common.Hash{
			{logTransferSigHash},
			{fromAddressHash},
		},
	}

	logs, err := client.FilterLogs(context.Background(), filter)
	if err != nil {
		// TODO: Add error logging
		log.Println(err)
		return nil
	}

	erc20Abi := getErc20Abi()
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, vLog := range logs {

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			var transferEvent Erc20Transfer

			err := erc20Abi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Println(err)
				return nil
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			transfer := reporter.Transfer{
				Chain:        chain,
				Txid:         vLog.TxHash.Hex(),
				Timestamp:    vLog.BlockHash.Hex(),
				TransferType: getTransferType(address, transferEvent.From.Hex()),
				From:         transferEvent.From.Hex(),
				To:           transferEvent.To.Hex(),
				Amount:       transferEvent.Value.Text(10), // TODO: Format to canonical amount
			}

			transfers = append(transfers, transfer)
		}
	}

	return transfers
}

func Watch(address string, chainName chains.ChainName, fromBlock uint64) error {
	// Set up channel
	transfers := make(chan reporter.Transfer)

	// Fire off goroutine to chase transfers
	// TODO: maxBlocks should be chain specific, an internal impl detail of ChaseTransfers
	go ChaseTransfers(chainName, address, fromBlock, 100000, transfers)

	// Fire off goroutine to process transfers
	go func() {
		for transfer := range transfers {
			store.LocalTransferStore.Add(address, transfer)
		}
	}()

	// No errors with watching addresses
	return nil
}

func getErc20Abi() abi.ABI {
	contractAbi, err := abi.JSON(strings.NewReader(Erc20MetaData.ABI))
	if err != nil {
		log.Fatal(err)
	}
	return contractAbi
}

func getTransferType(address string, from string) string {
	if strings.ToLower(address) == strings.ToLower(from) {
		return "SEND"
	}
	return "RECEIVE"
}
