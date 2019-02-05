package zsc

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
	"math/rand"
	"time"
)

type PublicZSCAPI struct {
}

func NewPublicZSCAPI() *PublicZSCAPI {
	return &PublicZSCAPI{}
}

func (api *PublicZSCAPI) CreateAccount() ([]byte, []byte, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // can i use the default source?
	x, y, err := bn256.RandomG1(r)
	if err != nil {
		return nil, nil, err
	}
	return x.Bytes(), y.Marshal(), nil
}

// CL and CR are 64-byte strings, x is a 32-byte string (in JS).
func (api *PublicZSCAPI) ReadBalance(CLHash common.Hash, CRHash common.Hash, xHash common.Hash) (int64, error) {
	CLBytes, err := CLHash.MarshalText()
	if err != nil {
		return 0, err
	}
	CL := new(bn256.G1)
	if _, err := CL.Unmarshal(CLBytes); err != nil {
		return 0, err
	}
	CRBytes, err := CRHash.MarshalText()
	if err != nil {
		return 0, err
	}
	CR := new(bn256.G1)
	if _, err := CR.Unmarshal(CRBytes); err != nil {
		return 0, err
	}
	x := new(big.Int)
	xBytes, err := xHash.MarshalText()
	if err != nil {
		return 0, err
	}
	x.UnmarshalText(xBytes)
	gb := CR.ScalarMult(CR, x)

	one := big.NewInt(1)
	end := big.NewInt(big.MaxPrec)
	for i := big.NewInt(0); i.Cmp(end) < 0; i.Add(i, one) {
		test := new(bn256.G1).ScalarBaseMult(i)
		if bytes.Compare(test.Marshal(), gb.Marshal()) == 0 {
			return i.Int64(), nil
		}
	}
	return 0, errors.New("Balance decryption failed!")
}
