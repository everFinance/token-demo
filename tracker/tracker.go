package tracker

import (
	"github.com/everFinance/goar"
	"github.com/everFinance/token-demo/token"
	"github.com/inconshreveable/log15"
)

var logger = log15.New(log15.Ctx{"module": "tracker"})

// Tracker transactions from Arweave
type Tracker struct {
	arClient    *goar.Client
	transaction chan token.Tx

	symbol string
	owner  string
	ids    []string

	IsSynced bool
}

func New(symbol, owner string) *Tracker {
	return &Tracker{
		arClient:    goar.NewClient("https://arweave.net"),
		transaction: make(chan token.Tx),

		symbol: symbol,
		owner:  owner,
		ids:    make([]string, 0),

		IsSynced: false,
	}
}

// Run Tracker, auto load txs from arweave
func (t *Tracker) Run() {
	t.runJobs()
}

func (t *Tracker) TransactionStream() <-chan token.Tx {
	return t.transaction
}

func (t *Tracker) TransactionsCount() int {
	return len(t.ids)
}
