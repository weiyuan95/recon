package main

import (
	"chaintx/chains"
	"chaintx/evm"
	"chaintx/reporter"
	"chaintx/store"
	"chaintx/tron"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func main() {
	r := gin.Default()

	// Kickoff the tron.Watch job when initialising the application
	if err := tron.Watch(chains.TronShasta); err != nil {
		log.Fatal("Failed to start tron.Watch for TronShasta:", err)
		return
	}
	if err := tron.Watch(chains.Tron); err != nil {
		log.Fatal("Failed to start tron.Watch for Tron:", err)
		return
	}

	// POST /watch
	// Example request:
	// ```json
	// {
	//   accounts: [
	//     { chain: 'ethMainnet', address: '0xabc' },
	//   ]
	// }
	// ```
	// Example responses:
	//   - HTTP 201 { status: 'processing' }
	//   - HTTP 200 { status: 'done' }
	//   - HTTP 400 { status: 'error', reason: 'Bad input.' }
	//   - HTTP 500 { status: 'error', reason: 'Maximum address limit reached.' }
	type Account struct {
		Chain     chains.ChainName `json:"chain"`
		Address   string           `json:"address"`
		FromBlock uint64           `json:"fromBlock"`
	}

	type WatchRequest struct {
		Accounts []Account `json:"accounts"`
	}

	r.POST("/watch", func(c *gin.Context) {

		// Bind request body to struct
		var request WatchRequest
		err := c.Bind(&request)
		if err != nil {
			log.Println("Encountered json unmarshal error:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"reason": "Bad input.",
			})
		}

		// Validate input first
		for _, account := range request.Accounts {
			if !chains.IsValidChain(account.Chain) {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"reason": fmt.Sprintf("Invalid chain: %s", account.Chain),
				})
				return
			}

			// TODO: if already watching the chain-address, return 200
		}

		// Fire off goroutines to watch accounts
		for _, account := range request.Accounts {
			log.Println("Watching", account.Chain, account.Address)

			if chains.IsEVM(account.Chain) {
				if err := evm.Watch(account.Address, account.Chain, account.FromBlock); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status": "error",
						"reason": err.Error(),
					})
					return
				}
			} else if chains.IsTvm(account.Chain) {
				// Add the address to be watched into the store, it will be picked up by the
				// tron.Watch scheduled job
				store.LocalWatchedAddressStore.Add(account.Chain, reporter.WatchedAddressInfo{
					Address: account.Address,
				})
			} else {
				// Not a valid chain, shouldn't hit this case since we already validate upfront
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"reason": fmt.Sprintf("Invalid chain: %s", account.Chain),
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "processing",
		})
	})

	// - GET /transfers?addresses=a,b,c&as=json/csv
	//   - 200 { status: 'processing' }
	//   - 200 { status: 'done', transfers: [...], next: 'cursor' }
	//   - 404 { status: 'error', reason: 'Not a registered address.' }
	// TODO: Handle `as`
	r.GET("/transfers", func(c *gin.Context) {

		addresses := c.Query("addresses")
		if addresses == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"reason": "No addresses provided.",
			})
			return
		}

		// Validate address list
		// TODO: validate for Tron addresses as well, for now, do not validate
		//addressesList := strings.Split(addresses, ",")
		//for _, address := range addressesList {
		//	if !common.IsHexAddress(address) {
		//		c.JSON(http.StatusBadRequest, gin.H{
		//			"status": "error",
		//			"reason": "Invalid address provided: " + address,
		//		})
		//		return
		//	}
		//}

		// Concatenate all transfers and return
		allTransfers := make([]reporter.Transfer, 0)
		for _, address := range strings.Split(addresses, ",") {
			allTransfers = append(allTransfers, store.LocalTransferStore.ListByAddress(address)...)
		}
		c.JSON(http.StatusOK, allTransfers)
	})

	// Health
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
