package detector

import (
	"github.com/everFinance/token-demo/cache"
)

// tx from Arweave, execute tx in token vm
func (d *Detector) runProcess() {
	d.tracker.Run()

	for {
		tx := <-d.tracker.TransactionStream()

		if err := d.token.ExecuteTx(tx); err != nil {
			logger.Warn("invalid tx to execute", "err", err)
			continue
		}

		d.cache.AddTx(cache.TxResponse{
			ID:     tx.ArTxID,
			Nonce:  tx.Nonce,
			Type:   tx.Type,
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
			Sign:   tx.Sign,
		})
	}
}
