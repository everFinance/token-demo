package issuer

import (
	"encoding/json"
	"time"

	"github.com/everFinance/dapp-tools/rollup"
	"github.com/everFinance/dapp-tools/tracker"
	"github.com/everFinance/token-demo/cache"
	"github.com/everFinance/token-demo/token"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logger = log15.New(log15.Ctx{"module": "issuer"})

type Issuer struct {
	engine *gin.Engine

	cache *cache.Cache
	token *token.Token

	tracker *tracker.Tracker
	rollup  *rollup.Rollup
}

func New(symbol, owner, keyPath string) *Issuer {
	tok := token.New(symbol, owner)

	return &Issuer{
		engine: gin.Default(),

		cache: cache.New(),
		token: tok,

		tracker: tracker.New(tok.Tags, "https://arweave.net", owner),
		rollup:  rollup.New("", "https://arweave.net", keyPath, owner, tok.Tags),
	}
}

func (i *Issuer) Run(port string) {
	go i.runAPI(port)
	go i.runTracker()
	go i.runRollup()
}

func (i *Issuer) runTracker() {
	i.tracker.Run()

	for {
		txRaw := <-i.tracker.SubscribeTx()

		tx := token.Tx{}
		if err := json.Unmarshal(txRaw.Data, &tx); err != nil {
			logger.Warn("can not unmarshal tx data from tracker", "txData", (txRaw.Data), "err", err)
			continue
		}
		tx.ArTxID = txRaw.ID
		tx.Owner = txRaw.Owner

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

func (i *Issuer) runRollup() {
	i.rollup.Run(time.Minute*1, 10000)
}
