package detector

import (
	"github.com/everFinance/token-demo/cache"
	"github.com/everFinance/token-demo/token"
	"github.com/everFinance/token-demo/tracker"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logger = log15.New(log15.Ctx{"module": "detector"})

// Detector for everyone to verify the calculation results
type Detector struct {
	engine  *gin.Engine
	token   *token.Token
	tracker *tracker.Tracker
	cache   *cache.Cache
}

func New(symbol, owner string) *Detector {
	return &Detector{
		engine:  gin.Default(),
		token:   token.New(symbol, owner),
		tracker: tracker.New(symbol, owner),
		cache:   cache.New(),
	}
}

func (d *Detector) Run(port string) {
	go d.runProcess()
	go d.runAPI(port)
}
