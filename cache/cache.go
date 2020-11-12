package cache

import (
	"sync"

	"github.com/everFinance/token-demo/token"
)

type TxResponse struct {
	ID     string `json:"id"` // AR tx id
	Nonce  string `json:"nonce"`
	Type   string `json:"type"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Sign   string `json:"sign"`
}

type Cache struct {
	pendingTxs  []token.Tx
	pendingLock sync.RWMutex

	packagedID   map[string]string // signature -> ar tx id
	packagedLock sync.RWMutex

	txs     []TxResponse
	txsLock sync.RWMutex

	txsByAddress     map[string][]TxResponse // address -> []Tx
	txsByAddressLock sync.RWMutex
}

func New() *Cache {
	return &Cache{
		pendingTxs:   []token.Tx{},
		packagedID:   map[string]string{},
		txsByAddress: map[string][]TxResponse{},
		txs:          []TxResponse{},
	}
}

// PendingTxs
func (c *Cache) GetPendingTxs() (txs []token.Tx) {
	c.pendingLock.RLock()
	defer c.pendingLock.RUnlock()

	txs = make([]token.Tx, len(c.pendingTxs))
	copy(txs, c.pendingTxs)
	return
}

func (c *Cache) AddPendingTx(tx token.Tx) {
	c.pendingLock.Lock()
	defer c.pendingLock.Unlock()

	c.pendingTxs = append(c.pendingTxs, tx)
}

func (c *Cache) EmptyPendingTx() {
	c.pendingLock.Lock()
	defer c.pendingLock.Unlock()

	c.pendingTxs = []token.Tx{}
}

// PackagedID
func (c *Cache) GetPackagedIDBySign(sign string) (id string) {
	c.packagedLock.RLock()
	defer c.packagedLock.RUnlock()

	id = c.packagedID[sign]
	return
}

func (c *Cache) AddPackagedID(sign, id string) {
	c.packagedLock.Lock()
	defer c.packagedLock.Unlock()

	c.packagedID[sign] = id
}

// Txs
func (c *Cache) GetTxs() (txs []TxResponse) {
	c.txsLock.RLock()
	defer c.txsLock.RUnlock()

	txs = make([]TxResponse, len(c.txs))
	copy(txs, c.txs)
	return
}

func (c *Cache) GetTxsByAddress(address string) (txs []TxResponse) {
	c.txsByAddressLock.RLock()
	tmp := c.txsByAddress[address]
	c.txsByAddressLock.RUnlock()

	for _, tx := range tmp {
		if tx.ID == "" {
			tx.ID = c.GetPackagedIDBySign(tx.Sign)
		}

		txs = append(txs, tx)
	}

	return
}

func (c *Cache) AddTx(tx TxResponse) {
	c.addTxs(tx)
	c.addTxsByAddress(tx)
}

func (c *Cache) addTxs(tx TxResponse) {
	c.txsLock.Lock()
	defer c.txsLock.Unlock()

	c.txs = append(c.txs, tx)
}

func (c *Cache) addTxsByAddress(tx TxResponse) {
	c.txsByAddressLock.Lock()
	defer c.txsByAddressLock.Unlock()

	if txs, ok := c.txsByAddress[tx.From]; ok {
		txs = append(txs, tx)
		c.txsByAddress[tx.From] = txs
	} else {
		c.txsByAddress[tx.From] = []TxResponse{tx}
	}

	if tx.To == "" {
		return
	}
	if txs, ok := c.txsByAddress[tx.To]; ok {
		txs = append(txs, tx)
		c.txsByAddress[tx.To] = txs
	} else {
		c.txsByAddress[tx.To] = []TxResponse{tx}
	}
}
