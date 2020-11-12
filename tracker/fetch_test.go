package tracker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustFetchTx(t *testing.T) {
	addr := MustFetchTxOwnerAddress("xkAigg50YoM6pisCXn0jk_6LDV_N9zHLXHYmgQj12BU", nil)
	assert.Equal(t, "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8", addr)
}

func TestMustFetchTxData(t *testing.T) {
	res := MustFetchTxData("xkAigg50YoM6pisCXn0jk_6LDV_N9zHLXHYmgQj12BU", nil)
	assert.Equal(t, `[{"nonce":"1598169690000","type":"claim","from":"0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9","to":"","amount":"100","sign":"0xaa7ff25bb0cd26e199d99cf2f2d771248c6e26b9c35eee62d70b45aa2f0284ea5dbf196d1159603fbdc23d6d6be73ce4235c97de36e3bc1e06f188a82a3d3f5c1b"},{"nonce":"1598169779000","type":"transfer","from":"0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9","to":"0xDc19464589c1cfdD10AEdcC1d09336622b282652","amount":"30","sign":"0x08a308a67440a07ab3bbd5075388732965c21f1885923de0429e2a98926a200339aff56240556df1e4c9486c8d94370a963d8d960e31730fa6eff49477d566bc1b"}]`, string(res))
}

func TestMustFetchIds(t *testing.T) {
	ids := MustFetchIds(`{
		"op": "and",
		"expr1": {
			"op": "equals",
			"expr1": "TokenSymbol",
			"expr2": "ROL"
		},
		"expr2": {
			"op": "equals",
			"expr1": "CreatedBy",
			"expr2": "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"
		}
	}`, nil)
	assert.True(t, len(ids) > 0)

	rIds := MustFetchIdsASC(`{
		"op": "and",
		"expr1": {
			"op": "equals",
			"expr1": "TokenSymbol",
			"expr2": "ROL"
		},
		"expr2": {
			"op": "equals",
			"expr1": "CreatedBy",
			"expr2": "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8"
		}
	}`, nil)
	assert.Equal(t, len(ids), len(rIds))
}
