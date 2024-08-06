package tron

import (
	"chaintx/reporter"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type TransferItem struct {
	TransactionHash, TransferFromAddress, TransferToAddress string
	Confirmed                                               bool
	Amount                                                  int // In SUN
	Timestamp                                               int
}

type TransferResult struct {
	Total int
	Data  []TransferItem
}

func TrxTransfers(tokenAddress string) ([]reporter.Transfer, error) {
	// The limit has a cap of 50 which is not documented
	resultLimit := 50
	start := 0
	client := NewThrottleClient(500)

	url := fmt.Sprintf(
		"https://apilist.tronscanapi.com/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
		resultLimit,
		start,
		tokenAddress,
	)

	get, err := client.Get(url)
	if err != nil {
		return []reporter.Transfer{}, errors.New("failed to fetch data")
	}

	body, err := io.ReadAll(get.Body)
	defer get.Body.Close()

	var data TransferResult
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
		return []reporter.Transfer{}, errors.New("failed to unmarshal data")
	}

	var result []reporter.Transfer

	for len(data.Data) > 0 {
		for _, transfer := range data.Data {
			transferType := "SEND"
			if transfer.TransferToAddress == tokenAddress {
				transferType = "RECEIVE"
			}

			// Use append since we have no way of knowing the exact number of rows,
			// since the API returns the wrong amount of rows
			result = append(result, reporter.Transfer{
				Txid:         transfer.TransactionHash,
				From:         transfer.TransferFromAddress,
				To:           transfer.TransferToAddress,
				Timestamp:    strconv.Itoa(transfer.Timestamp),
				Amount:       strconv.Itoa(transfer.Amount),
				TransferType: transferType,
			})

		}

		start += len(data.Data)
		url = fmt.Sprintf(
			"https://apilist.tronscanapi.com/api/trx/transfer?sort=-timestamp&count=true&limit=%d&start=%d&address=%s&filterTokenValue=0",
			resultLimit,
			start,
			tokenAddress,
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
