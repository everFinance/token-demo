package tracker

import (
	"time"

	"github.com/go-co-op/gocron"
)

func (t *Tracker) runJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(2).Minutes().Do(t.jobTxsPull)

	s.StartImmediately()
	s.StartAsync()
}
func reverseIDs(ids []string) (rIDs []string) {
	for i := len(ids) - 1; i >= 0; i-- {
		rIDs = append(rIDs, ids[i])
	}
	return
}

func (t *Tracker) jobTxsPull() {

	logger.Info("job txs pull running...")
	defer func() {
		logger.Info("job txs pull done")
		t.IsSynced = true
	}()

	// get all Arweave txs
	ids := MustFetchIds(`
	{
		transactions(
			first: 100
			tags: [
					{
							name: "TokenSymbol",
							values: "`+t.symbol+`"
					},
					{
							name: "CreatedBy",
							values: "`+t.owner+`"
					},
			]
			sort: HEIGHT_ASC
		) {
			edges {
				node {
					id
				}
			}
		}
	}`, t.arClient)

	if len(ids) <= len(t.ids) {
		return
	}

	newIds := reverseIDs(ids)

	// process txs
	txsCounts := 0
	for _, id := range newIds {
		txs, err := t.fetchTokenTxsByID(id)
		if err != nil {
			logger.Warn("invalid tx from Arweave", "id", id, "err", err)
			continue
		}

		for _, tx := range txs {
			t.transaction <- tx

			txsCounts++
		}
	}

	// saved newIds
	t.ids = append(t.ids, newIds...)
	logger.Info("job txs processed", "counts", txsCounts)

}
