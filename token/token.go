package token

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	TxTypeTransfer = "transfer"
	TxTypeClaim    = "claim"
)

// Tx is transaction from API or Arweave
// Nonce: timestamp of transaction
// Type: transfer or claim, each address can only be claimed once
// Amount: only supports integer
// Sign: personalSign(hash(Nonce + Type + From + To + Amount))
// Owner: only accept tx from owner
type Tx struct {
	Nonce  string `json:"nonce"`
	Type   string `json:"type"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Sign   string `json:"sign"`

	Owner  string `json:"-"`
	ArTxID string `json:"-"`
}

// Token model
// ExecuteTx from API or Arweave
type Token struct {
	Symbol string
	Owner  string

	balances map[string]*big.Int // address -> blance
	nonces   map[string]int64    // address -> nonce
	minted   map[string]bool     // address -> minted(bool)

	lock sync.RWMutex
}

func New(symbol, owner string) *Token {
	return &Token{
		Symbol:   symbol,
		Owner:    owner,
		balances: make(map[string]*big.Int),
		nonces:   make(map[string]int64),
		minted:   make(map[string]bool),
	}
}

// VM for transaction execute
func (t *Token) ExecuteTx(tx Tx) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err = t.txVerify(tx); err != nil {
		return
	}

	amount, ok := new(big.Int).SetString(tx.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount: %v", amount)
	}

	switch tx.Type {
	case TxTypeTransfer:
		err = t.transfer(tx.From, tx.To, amount)
	case TxTypeClaim:
		err = t.mint(tx.From, amount)
	default:
		err = fmt.Errorf("invalid tx type: %v", tx.Type)
	}

	if err != nil {
		return
	}

	// update status
	nonce, _ := strconv.ParseInt(tx.Nonce, 10, 64)
	t.nonces[tx.From] = nonce
	return
}

func (t *Token) TotalSupply() *big.Int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	total := big.NewInt(0)
	for _, bal := range t.balances {
		total = new(big.Int).Add(total, bal)
	}
	return total
}

func (t *Token) BalanceOf(addr string) *big.Int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.balanceOf(addr)
}

func (t *Token) transfer(from, to string, amount *big.Int) error {
	if err := t.sub(from, amount); err != nil {
		return err
	}
	if err := t.add(to, amount); err != nil {
		return err
	}

	return nil
}

func (t *Token) balanceOf(addr string) *big.Int {
	if bal, ok := t.balances[addr]; ok {
		return bal
	}

	return big.NewInt(0)
}

func (t *Token) mint(addr string, amount *big.Int) error {
	if err := t.add(addr, amount); err != nil {
		return err
	}

	t.minted[addr] = true
	return nil
}

func (t *Token) burn(addr string, amount *big.Int) error {
	return t.sub(addr, amount)
}

func (t *Token) add(addr string, amount *big.Int) error {
	if amount == nil {
		return errors.New("nil amount")
	}

	bal := t.balanceOf(addr)
	t.balances[addr] = new(big.Int).Add(bal, amount)
	return nil
}

func (t *Token) sub(addr string, amount *big.Int) error {
	if amount == nil {
		return errors.New("nil amount")
	}

	bal := t.balanceOf(addr)
	if bal.Cmp(amount) < 0 {
		return errors.New("not enough amount")
	}

	if bal.Cmp(amount) == 0 {
		delete(t.balances, addr)
		return nil
	}

	t.balances[addr] = new(big.Int).Sub(bal, amount)
	return nil
}

func (t *Token) nonce(addr string) int64 {
	if nonce, ok := t.nonces[addr]; ok {
		return nonce
	}

	return 0
}

func (t *Token) txVerify(tx Tx) (err error) {
	// owner verify
	if tx.Owner != t.Owner {
		return fmt.Errorf("invalid owner: %v", tx.Owner)
	}

	// nonce verify
	nonce, err := strconv.ParseInt(tx.Nonce, 10, 64)
	if err != nil {
		return
	}
	currentNonce := t.nonce(tx.From)
	if nonce <= currentNonce {
		return fmt.Errorf("invalid nonce: %d, current nonce: %d", nonce, currentNonce)
	}

	// do not transfer to self
	if tx.From == tx.To {
		return fmt.Errorf("transfer to self: %v", tx.To)
	}

	// amount verify
	if tx.Amount == "0" {
		return fmt.Errorf("amount is zero")
	}
	amount, ok := new(big.Int).SetString(tx.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount: %v", tx.Amount)
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("amount is zero")
	}
	if amount.Cmp(big.NewInt(1)) < 0 {
		return fmt.Errorf("negative amount: %v", tx.Amount)
	}

	// only mint once
	if tx.Type == TxTypeClaim && t.minted[tx.From] {
		return fmt.Errorf("minted")
	}

	// signature verfiy
	sign := common.FromHex(tx.Sign)
	if len(sign) != 65 {
		return fmt.Errorf("invalid length of signture")
	}

	if sign[64] != 27 && sign[64] != 28 {
		return fmt.Errorf("invalid signature type")
	}
	sign[64] -= 27

	signData := tx.Nonce + tx.Type + tx.From + tx.To + tx.Amount
	hash, _ := accounts.TextAndHash([]byte(signData))

	recoverPub, err := crypto.Ecrecover(hash, sign)
	if err != nil {
		return
	}
	pubKey, err := crypto.UnmarshalPubkey(recoverPub)
	if err != nil {
		return
	}
	addr := crypto.PubkeyToAddress(*pubKey)

	if addr.String() != tx.From {
		return fmt.Errorf("invalid signature")
	}

	return
}
