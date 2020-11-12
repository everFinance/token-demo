package issuer

import (
	"github.com/everFinance/token-demo/cache"
)

func (i *Issuer) runProcess() {
	i.tracker.Run()

	for {
		tx := <-i.tracker.TransactionStream()

		if !i.tracker.IsSynced {
			if err := i.token.ExecuteTx(tx); err != nil {
				logger.Warn("invalid tx to execute", "err", err)
				continue
			}

			i.cache.AddTx(cache.TxResponse{
				ID:     tx.ArTxID,
				Nonce:  tx.Nonce,
				Type:   tx.Type,
				From:   tx.From,
				To:     tx.To,
				Amount: tx.Amount,
				Sign:   tx.Sign,
			})
		}

		i.cache.AddPackagedID(tx.Sign, tx.ArTxID)

	}
}
