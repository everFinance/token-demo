package token

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	token := New("BTC", "0x123")
	assert.Equal(t, "BTC", token.Symbol)

	assert.Equal(t, "nil amount", token.mint("0x123", nil).Error())

	// mint
	assert.NoError(t, token.mint("0x123", big.NewInt(10)))

	assert.Equal(t, big.NewInt(10), token.TotalSupply())
	assert.Equal(t, big.NewInt(10), token.BalanceOf("0x123"))
	assert.Equal(t, big.NewInt(0), token.BalanceOf("0x888"))

	// transfer
	assert.NoError(t, token.transfer("0x123", "0x888", big.NewInt(8)))

	assert.Equal(t, big.NewInt(10), token.TotalSupply())
	assert.Equal(t, big.NewInt(2), token.BalanceOf("0x123"))
	assert.Equal(t, big.NewInt(8), token.BalanceOf("0x888"))

	assert.Equal(t, "not enough amount", token.transfer("0x123", "0x888", big.NewInt(3)).Error())

	// burn
	assert.Equal(t, "not enough amount", token.burn("0x888", big.NewInt(9)).Error())
	assert.NoError(t, token.burn("0x888", big.NewInt(8)))
	assert.Equal(t, big.NewInt(2), token.TotalSupply())
	assert.Equal(t, big.NewInt(2), token.BalanceOf("0x123"))
	assert.Equal(t, big.NewInt(0), token.BalanceOf("0x888"))

	// executeTx
	assert.Equal(t, "not enough amount", token.ExecuteTx(
		Tx{
			Nonce:  "123",
			Type:   "transfer",
			From:   "0x5583900EB2d48761AEc176A39aB118bf5c4208BD",
			To:     "0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9",
			Amount: "1",
			Sign:   "0x90d4eea4019ef72b0cf4cc660ab7143889da751d748440eeab65640715b60aca1aa3fad04d6259bece790831ab684650f7d7d4d53db667a9c3e0c7221e0f70de1b",
			Owner:  "0x123",
		},
	).Error())
	token.mint("0x5583900EB2d48761AEc176A39aB118bf5c4208BD", big.NewInt(10))
	assert.NoError(t, token.ExecuteTx(
		Tx{
			Nonce:  "123",
			Type:   "transfer",
			From:   "0x5583900EB2d48761AEc176A39aB118bf5c4208BD",
			To:     "0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9",
			Amount: "1",
			Sign:   "0x90d4eea4019ef72b0cf4cc660ab7143889da751d748440eeab65640715b60aca1aa3fad04d6259bece790831ab684650f7d7d4d53db667a9c3e0c7221e0f70de1b",
			Owner:  "0x123",
		}))
	assert.Equal(t, big.NewInt(9), token.BalanceOf("0x5583900EB2d48761AEc176A39aB118bf5c4208BD"))
	assert.Equal(t, big.NewInt(1), token.BalanceOf("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"))
}

func TestTxVerify(t *testing.T) {
	token := New("BTC", "0x123")
	tx := Tx{
		Nonce:  "123",
		Type:   "transfer",
		From:   "0x5583900EB2d48761AEc176A39aB118bf5c4208BD",
		To:     "0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9",
		Amount: "1",
		Sign:   "0x90d4eea4019ef72b0cf4cc660ab7143889da751d748440eeab65640715b60aca1aa3fad04d6259bece790831ab684650f7d7d4d53db667a9c3e0c7221e0f70de1b",
		Owner:  "0x123",
	}

	assert.NoError(t, token.txVerify(tx))
}
