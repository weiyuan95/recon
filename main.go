package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// REST API
// - POST /watch { addresses: [a,b,c] }
//
// - GET /transfers?addresses=a,b,c&as=json/csv
//   - 200 { status: 'processing' }
//   - 200 { status: 'done', transfers: [...], next: 'cursor' }
//   - 404 { status: 'error', reason: 'Not a registered address.' }
//   - 500 { status: 'error', reason: 'Maximum address limit reached.' }

func main() {
	r := gin.Default()
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
