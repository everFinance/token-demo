package issuer

import (
	"encoding/json"
	"time"

	"github.com/everFinance/goar/types"
	"github.com/go-co-op/gocron"
)

// rollup txs to Arweave
func (i *Issuer) runRollup() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Minutes().Do(i.jobRollup)

	s.StartAsync()
}

func (i *Issuer) jobRollup() {

	logger.Info("job rollup running...")
	defer logger.Info("job rollup done")

	txs := i.cache.GetPendingTxs()

	pendingCounts := len(txs)
	if pendingCounts == 0 {
		logger.Info("pendingTxs is empty")
		return
	}

	pending, err := json.Marshal(txs)
	if err != nil {
		logger.Error("pending marshal failed", "txs", txs, "err", err)
		return
	}

	id, status, err := i.wallet.SendData(pending, []types.Tag{
		types.Tag{
			Name:  "TokenSymbol",
			Value: i.Symbol,
		},
		types.Tag{
			Name:  "Version",
			Value: i.Version,
		},
		types.Tag{
			Name:  "CreatedBy",
			Value: i.Owner,
		},
	})

	if status != "OK" {
		logger.Error("submit pendingTxs failed", "id", id, "status", status, "err", err)
		return
	}

	i.cache.EmptyPendingTx()
	logger.Info("submit success", "id", id, "pendingCounts", pendingCounts)

}
