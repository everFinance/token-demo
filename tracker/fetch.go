package tracker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/everFinance/goar/client"
	"github.com/everFinance/goar/utils"
	"github.com/everFinance/token-demo/token"
)

func (t *Tracker) fetchTokenTxsByID(id string) (txs []token.Tx, err error) {
	data := MustFetchTxData(id, t.arClient)
	tmp := []token.Tx{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		err = fmt.Errorf("data: %v, err: %v", data, err)
		return
	}

	ownerAddr := MustFetchTxOwnerAddress(id, t.arClient)
	for _, tx := range tmp {
		tx.Owner = ownerAddr
		tx.ArTxID = id
		txs = append(txs, tx)
	}

	return
}

func MustFetchTxOwnerAddress(id string, c *client.Client) (address string) {
	if c == nil {
		c = client.New("https://arweave.net")
	}

	for {
		owner, err := c.GetTransactionField(id, "owner")
		if err == nil {
			address = utils.OwnerToAddress(owner)
			break
		}

		logger.Warn("fetch tx, retry 10 secs", "id", id, "err", err)
		time.Sleep(10 * time.Second)
	}

	return
}

func MustFetchTxData(id string, c *client.Client) (res []byte) {
	if c == nil {
		c = client.New("https://arweave.net")
	}

	var err error
	for {
		res, err = c.GetTransactionData(id, "json")
		if err == nil {
			break
		}

		logger.Warn("fetch tx data failed, retry 10 secs", "id", id, "body", string(res), "err", err)
		time.Sleep(10 * time.Second)
	}

	return
}

func MustFetchIds(arql string, c *client.Client) (ids []string) {
	if c == nil {
		c = client.New("https://arweave.net")
	}

	var err error

	for {
		ids, err = c.Arql(arql)
		if err == nil {
			break
		}

		logger.Warn("fetch ids failed, retry 3 secs", "err", err)
		time.Sleep(3 * time.Second)
	}

	return
}

func MustFetchIdsASC(arql string, c *client.Client) (rIds []string) {
	ids := MustFetchIds(arql, c)

	for i := len(ids) - 1; i > -1; i-- {
		rIds = append(rIds, ids[i])
	}
	return
}
