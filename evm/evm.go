package evm

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"time"
)

func getErc20Abi() abi.ABI {
	contractAbi, err := abi.JSON(strings.NewReader(string(Erc20MetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}
	return contractAbi
}

func ChaseEvmTransfers(
	client *ethclient.Client,
	fromBlock *big.Int,
	toBlock *big.Int,
	interval int,
) {

	for {
		// fromBlock >= toBlock
		if fromBlock.Cmp(toBlock) >= 0 {
			log.Println("Waiting for new blocks...")
			time.Sleep(1 * time.Second)

			blockNumber, err := client.BlockNumber(context.Background())
			if err != nil {
			}

			// What if blockNumber is 2^64?!
			// Well, then the chain has another problem on their hands lmao xD
			toBlock.Set(big.NewInt(int64(blockNumber)))

			// TODO: blockInterval = min(BigInt(endBlock - fromBlock), 100000n);

			continue
		}

		// TODO: Implement chasing
	}
}

func EvmTransfers(
	client *ethclient.Client,
	address string,
	fromBlock *big.Int,
	toBlock *big.Int,
) []Transfer {
	var transfers []Transfer

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	fromAddressHash := common.HexToHash(address)

	filter := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
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

			transfer := Transfer{
				txid:         vLog.TxHash.Hex(),
				timestamp:    vLog.BlockHash.Hex(),
				transferType: getTransferType(address, transferEvent.From.Hex()),
				from:         transferEvent.From.Hex(),
				to:           transferEvent.To.Hex(),
				amount:       transferEvent.Value.Text(10), // TODO: Format to canonical amount
			}

			transfers = append(transfers, transfer)
		}
	}

	return transfers
}

func getTransferType(address string, from string) string {
	if strings.ToLower(address) == strings.ToLower(from) {
		return "SEND"
	}
	return "RECEIVE"
}

// 8===========D----

// TODO: DRY this shit away
type Transfer struct {
	txid         string
	timestamp    string
	transferType string
	from         string
	to           string
	amount       string
}

//// Erc20Transfer The struct returned from the eth node's `getLogs` rpc call
//type Erc20Transfer struct {
//	From   common.Address
//	To     common.Address
//	Tokens *big.Int
//}
