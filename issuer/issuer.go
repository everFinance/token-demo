package issuer

import (
	"github.com/everFinance/goar/wallet"
	"github.com/everFinance/token-demo/cache"
	"github.com/everFinance/token-demo/token"
	"github.com/everFinance/token-demo/tracker"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logger = log15.New(log15.Ctx{"module": "detector"})

type Issuer struct {
	Symbol  string
	Owner   string
	Version string

	engine  *gin.Engine
	token   *token.Token
	tracker *tracker.Tracker
	cache   *cache.Cache
	wallet  *wallet.Wallet
}

func New(symbol, owner, keyPath string) *Issuer {
	wallet, err := wallet.NewFromPath(keyPath)
	if err != nil {
		panic(err)
	}

	return &Issuer{
		Symbol:  symbol,
		Owner:   owner,
		Version: "0.1",

		engine:  gin.Default(),
		token:   token.New(symbol, owner),
		tracker: tracker.New(symbol, owner),
		cache:   cache.New(),
		wallet:  wallet,
	}
}

func (i *Issuer) Run(port string) {
	go i.runAPI(port)
	go i.runProcess()
	go i.runRollup()
}
