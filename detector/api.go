package detector

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Detector) runAPI(port string) {
	d.engine.GET("/balanceOf/:address", d.balanceOf)
	d.engine.GET("/txs/:address", d.txsByAddress)
	d.engine.GET("/txs", d.txs)

	d.engine.Run(port)
}

func (d *Detector) balanceOf(c *gin.Context) {
	address := c.Param("address")
	amount := d.token.BalanceOf(address)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": amount.String(),
	})
}

func (d *Detector) txsByAddress(c *gin.Context) {
	address := c.Param("address")
	txs := d.cache.GetTxsByAddress(address)

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"txs":     txs,
	})
}

func (d *Detector) txs(c *gin.Context) {
	txs := d.cache.GetTxs()

	c.JSON(http.StatusOK, gin.H{
		"txs": txs,
	})
}
