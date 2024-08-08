package tron

import (
	"chaintx/chains"
	"chaintx/reporter"
	"chaintx/scheduler"
	"chaintx/store"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"strconv"
)

type TrxTransferItem struct {
	TransactionHash, TransferFromAddress, TransferToAddress string
	Confirmed                                               bool
	Amount                                                  int // In SUN
	Timestamp                                               int
}

type TrxTransferResult struct {
	Total int
	Data  []TrxTransferItem
}

func TrxTransfers(walletAddress string, chainName chains.ChainName) ([]reporter.Transfer, error) {
	// The limit has a cap of 50 which is not documented
	resultLimit := 50
	start := 0
	client := NewThrottleClient(500)
	var baseUrl string

	// TODO: better way to handle this?
	if chainName == chains.Tron {
		baseUrl = "https://apilist.tronscanapi.com"
	} else if chainName == chains.TronShasta {
		baseUrl = "https://shastapi.tronscan.org"
	} else {
		return []reporter.Transfer{}, errors.New("invalid chain name")
	}

	url := fmt.Sprintf(
		"%s/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
		baseUrl,
		resultLimit,
		start,
		walletAddress,
	)

	get, err := client.Get(url)
	if err != nil {
		return []reporter.Transfer{}, errors.New("failed to fetch data")
	}

	body, err := io.ReadAll(get.Body)
	defer get.Body.Close()

	var data TrxTransferResult
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		return []reporter.Transfer{}, errors.New("failed to unmarshal data")
	}

	var result []reporter.Transfer

	for len(data.Data) > 0 {
		for _, transfer := range data.Data {
			transferType := "SEND"
			if transfer.TransferToAddress == walletAddress {
				transferType = "RECEIVE"
			}
			canonicalAmount, err := toCanonical(strconv.Itoa(transfer.Amount), 6)

			if err != nil {
				return []reporter.Transfer{}, errors.New("failed to convert amount")
			}

			// Use append since we have no way of knowing the exact number of rows,
			// since the API returns the wrong amount of rows
			result = append(result, reporter.Transfer{
				Chain:        chainName,
				Txid:         transfer.TransactionHash,
				From:         transfer.TransferFromAddress,
				To:           transfer.TransferToAddress,
				Timestamp:    strconv.Itoa(transfer.Timestamp),
				Amount:       canonicalAmount,
				TransferType: transferType,
				TokenType:    "TRX",
			})

		}

		start += len(data.Data)
		url = fmt.Sprintf(
			"https://apilist.tronscanapi.com/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
			resultLimit,
			start,
			walletAddress,
		)

		get, err = client.Get(url)
		if err != nil {
			return []reporter.Transfer{}, errors.New("failed to fetch data")
		}

		body, err = io.ReadAll(get.Body)
		// Do not defer since we are in a loop - close it immediately after reading
		// this is because the code that is deferred will only run _after_ the loop completes
		get.Body.Close()

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println(err)
			return []reporter.Transfer{}, errors.New("failed to unmarshal data")
		}

	}

	return result, nil
}

type Trc20TokenInfo struct {
	TokenAbbr    string
	TokenDecimal int
}

type Trc20TransferItem struct {
	TransactionId  string `json:"transaction_id"`
	BlockTimestamp int    `json:"block_ts"`
	From           string `json:"from_address"`
	To             string `json:"to_address"`
	Block          int
	Amount         string `json:"quant"`
	Confirmed      bool
	TokenInfo      Trc20TokenInfo
}

type Trc20TransferResult struct {
	TokenTransfers []Trc20TransferItem `json:"token_transfers"`
}

func Trc20Transfers(walletAddress string, chainName chains.ChainName) ([]reporter.Transfer, error) {
	// The limit has a cap of 50 which is not documented
	resultLimit := 50
	start := 0
	client := NewThrottleClient(500)

	var baseUrl string
	var usdtContractAddress string

	// TODO: better way to handle this?
	if chainName == chains.Tron {
		baseUrl = "https://apilist.tronscanapi.com"
		usdtContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	} else if chainName == chains.TronShasta {
		baseUrl = "https://shastapi.tronscan.org"
		usdtContractAddress = "TG3XXyExBkPp9nzdajDZsozEu4BkaSJozs"
	} else {
		return []reporter.Transfer{}, errors.New("invalid chain name")
	}

	url := fmt.Sprintf(
		"%s/api/filter/trc20/transfers?sort=-timestamp&limit=%d&start=%d&trc20Id=%s&relatedAddress=%s&filterTokenValue=0",
		baseUrl,
		resultLimit,
		start,
		usdtContractAddress,
		walletAddress,
	)

	get, err := client.Get(url)
	if err != nil {
		return []reporter.Transfer{}, errors.New("failed to fetch data")
	}

	body, err := io.ReadAll(get.Body)
	defer get.Body.Close()

	var result []reporter.Transfer

	var data Trc20TransferResult
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		return []reporter.Transfer{}, errors.New("failed to unmarshal data")
	}

	for len(data.TokenTransfers) > 0 {
		for _, transfer := range data.TokenTransfers {
			transferType := "SEND"
			if transfer.To == walletAddress {
				transferType = "RECEIVE"
			}

			canonicalAmount, err := toCanonical(transfer.Amount, 6)

			if err != nil {
				return []reporter.Transfer{}, errors.New("failed to convert amount")
			}

			// Use append since we have no way of knowing the exact number of rows,
			// since the API returns the wrong amount of rows
			result = append(result, reporter.Transfer{
				Chain:        chainName,
				Txid:         transfer.TransactionId,
				From:         transfer.From,
				To:           transfer.To,
				Timestamp:    strconv.Itoa(transfer.BlockTimestamp),
				Amount:       canonicalAmount,
				TransferType: transferType,
				TokenType:    transfer.TokenInfo.TokenAbbr,
			})

		}

		start += len(data.TokenTransfers)
		url := fmt.Sprintf(
			"%s/api/filter/trc20/transfers?sort=-timestamp&limit=%d&start=%d&trc20Id=%s&relatedAddress=%s&filterTokenValue=0",
			baseUrl,
			resultLimit,
			start,
			usdtContractAddress,
			walletAddress,
		)

		get, err = client.Get(url)
		if err != nil {
			return []reporter.Transfer{}, errors.New("failed to fetch data")
		}

		body, err = io.ReadAll(get.Body)
		// Do not defer since we are in a loop - close it immediately after reading
		// this is because the code that is deferred will only run _after_ the loop completes
		get.Body.Close()

		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println(err)
			return []reporter.Transfer{}, errors.New("failed to unmarshal data")
		}

	}

	return result, nil
}

func toCanonical(amount string, decimals int) (string, error) {
	bigFloat, _, _ := new(big.Float).Parse(amount, 10)
	final := bigFloat.Mul(bigFloat, big.NewFloat(math.Pow(10, float64(-decimals))))
	return final.Text('f', decimals), nil
}

func Watch(chainName chains.ChainName) error {
	// Schedule a job to pull all the Tron addresses we are watching, then pull all the TRX and TRC20 data for
	// those addresses
	scheduler.Schedule(func() {
		infos, present := store.LocalWatchedAddressStore.Get(chainName)
		// TODO: Goroutine?
		if present {
			for _, info := range infos {
				trxTransfers, err := TrxTransfers(info.Address, chainName)
				if err != nil {
					fmt.Println("Failed to fetch TRX trxTransfers for", info.Address)
					return
				}

				trc20Transfers, err := Trc20Transfers(info.Address, chainName)
				if err != nil {
					fmt.Println("Failed to fetch TRC20 trxTransfers for", info.Address)
					return
				}

				for _, transfer := range trxTransfers {
					store.LocalTransferStore.Add(info.Address, transfer)
				}

				for _, transfer := range trc20Transfers {
					store.LocalTransferStore.Add(info.Address, transfer)
				}
			}
		}

	}, 20_000) // Every 10 seconds TODO: change to a longer time for 'prod' use

	return nil
}
