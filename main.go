package main

import (
	"chaintx/chains"
	"chaintx/evm"
	"chaintx/reporter"
	"chaintx/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func main() {
	r := gin.Default()
	localStore := store.NewTransferStore()

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
					"reason": "Invalid chain.",
				})
				return
			}

			// TODO: if already watching the chain-address, return 200
		}

		// Fire off goroutines to watch accounts
		for _, account := range request.Accounts {
			log.Println("Watching", account.Chain, account.Address)

			// Set up client
			client, err := evm.GetClient(account.Chain)
			if err != nil {
				// Since we validated the chain earlier, this should never happen. But let's just be defensive here.
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"reason": err.Error(),
				})
			}

			// Set up channel
			transfers := make(chan reporter.Transfer)

			// Fire off goroutine to chase transfers
			// TODO: maxBlocks should be chain specific, an internal impl detail of ChaseTransfers
			go evm.ChaseTransfers(client, account.Address, account.FromBlock, 100000, transfers)

			// Fire off goroutine to process transfers
			go func() {
				for transfer := range transfers {
					localStore.Add(account.Address, transfer)
				}
			}()
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
		addressesList := strings.Split(addresses, ",")
		for _, address := range addressesList {
			if !common.IsHexAddress(address) {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"reason": "Invalid address provided: " + address,
				})
				return
			}
		}

		// Concatenate all transfers and return
		allTransfers := make([]reporter.Transfer, 0)
		for _, address := range addressesList {
			allTransfers = append(allTransfers, localStore.ListByAddress(address)...)
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
