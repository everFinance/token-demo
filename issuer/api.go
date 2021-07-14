package issuer

import (
	"encoding/json"
	"net/http"

	"github.com/everFinance/token-demo/cache"
	"github.com/everFinance/token-demo/token"
	"github.com/gin-gonic/gin"
)

func (i *Issuer) runAPI(port string) {
	i.engine.Static("/token", "./web/")
	i.engine.Use(i.waitingSync)
	i.engine.GET("/balanceOf/:address", i.balanceOf)
	i.engine.GET("/txs/:address", i.txsByAddress)
	i.engine.GET("/txs", i.txs)
	i.engine.POST("/tx", i.submitTx)

	i.engine.Run(port)
}

func (i *Issuer) waitingSync(c *gin.Context) {
	if !i.tracker.IsSynced {
		c.JSON(http.StatusOK, gin.H{
			"error": "node sync...",
		})
		c.Abort()
	}

	c.Next()
}

func (i *Issuer) balanceOf(c *gin.Context) {
	address := c.Param("address")
	amount := i.token.BalanceOf(address)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": amount.String(),
	})
}

func (i *Issuer) txsByAddress(c *gin.Context) {
	address := c.Param("address")
	txs := i.cache.GetTxsByAddress(address)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"txs":     txs,
	})
}

func (i *Issuer) txs(c *gin.Context) {
	txs := i.cache.GetTxs()

	c.JSON(http.StatusOK, gin.H{
		"txs": txs,
	})
}

func (i *Issuer) submitTx(c *gin.Context) {
	tx := token.Tx{
		Owner: i.token.Owner,
	}
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := i.token.ExecuteTx(tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txData, err := json.Marshal(tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i.rollup.AddTx() <- txData

	i.cache.AddPendingTx(tx)
	i.cache.AddTx(cache.TxResponse{
		Nonce:  tx.Nonce,
		Type:   tx.Type,
		From:   tx.From,
		To:     tx.To,
		Amount: tx.Amount,
		Sign:   tx.Sign,
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
