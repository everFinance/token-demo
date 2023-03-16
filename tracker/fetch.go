package tracker

import (
	"encoding/json"
	"fmt"
	"github.com/everFinance/goar"
	"time"

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

func MustFetchTxOwnerAddress(id string, c *goar.Client) (address string) {
	if c == nil {
		c = goar.NewClient("https://arweave.net")
	}

	for {
		owner, err := c.GetTransactionField(id, "owner")
		if err == nil {
			address, _ = utils.OwnerToAddress(owner)
			break
		}

		logger.Warn("fetch tx, retry 10 secs", "id", id, "err", err)
		time.Sleep(10 * time.Second)
	}

	return
}

func MustFetchTxData(id string, c *goar.Client) (res []byte) {
	if c == nil {
		c = goar.NewClient("https://arweave.net")
	}

	var err error
	for {
		res, err = c.GetTransactionDataByGateway(id)
		if err == nil {
			break
		}
		fmt.Println("jsondata: ", string(res))

		logger.Warn("fetch tx data failed, retry 10 secs", "id", id, "body", string(res), "err", err)
		time.Sleep(10 * time.Second)
	}

	return
}

func MustFetchIds(query string, c *goar.Client) (ids []string) {
	if c == nil {
		c = goar.NewClient("https://arweave.net")
	}

	var err error

	for {
		ids, err = fetchIds(query, c)
		if err == nil {
			break
		}
	}

	return
}

func fetchIds(query string, c *goar.Client) (ids []string, err error) {
	data, err := c.GraphQL(query)
	if err != nil {
		return
	}

	txs := struct {
		Transactions struct {
			Edges []struct{ Node struct{ ID string } }
		}
	}{}

	if err = json.Unmarshal(data, &txs); err != nil {
		return
	}

	for _, ed := range txs.Transactions.Edges {
		ids = append(ids, ed.Node.ID)
	}

	return
}
