package tracker

import (
	"github.com/everFinance/goar"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustFetchTx(t *testing.T) {
	addr := MustFetchTxOwnerAddress("xkAigg50YoM6pisCXn0jk_6LDV_N9zHLXHYmgQj12BU", nil)
	assert.Equal(t, "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8", addr)
}

func TestMustFetchTxData(t *testing.T) {
	res := MustFetchTxData("RRdqyGM-VtFJ7hzXrWaZWOzBTHQdrBnpGRBpwQdDaSA", goar.New("https://seed-dev.everpay.io"))
	t.Log(string(res))
}

func TestNew(t *testing.T) {
	c := goar.New("https://seed-dev.everpay.io")
	res, err := c.GetTransactionData("RRdqyGM-VtFJ7hzXrWaZWOzBTHQdrBnpGRBpwQdDaSA", "json")
	t.Log(res)
	t.Log(err)
}

func TestMustFetchIds(t *testing.T) {
	ids := MustFetchIds(`
	{
		transactions(
			first: 10000
			tags: [
					{
							name: "TokenSymbol",
							values: "ROL"
					},
					{
							name: "CreatedBy",
							values: "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"
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
	}`, nil)
	assert.True(t, len(ids) > 0)
}
