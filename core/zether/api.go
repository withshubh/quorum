package zether

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
	"math/rand"
	"time"
)

type PublicZetherAPI struct {
}

func NewPublicZetherAPI() *PublicZetherAPI {
	return &PublicZetherAPI{}
}

func (api *PublicZetherAPI) CreateAccount() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // can i use the default source?
	x, y, err := bn256.RandomG1(r)
	if err != nil {
		return nil, err
	}
	result["x"] = common.BytesToHash(x.Bytes())
	result["y"] = [2]common.Hash{common.BytesToHash(y.Marshal()[:32]), common.BytesToHash(y.Marshal()[32:])}
	return result, nil
}

// CL and CR are 64-byte strings, x is a 32-byte string (in JS).
func (api *PublicZetherAPI) ReadBalance(CLBytes [2]common.Hash, CRBytes [2]common.Hash, xHash common.Hash) (int64, error) {
	CL := new(bn256.G1)
	if _, err := CL.Unmarshal(append(CLBytes[0].Bytes(), CLBytes[1].Bytes()...)); err != nil {
		return 0, err
	}
	CR := new(bn256.G1)
	if _, err := CR.Unmarshal(append(CRBytes[0].Bytes(), CRBytes[1].Bytes()...)); err != nil {
		return 0, err
	}
	x := new(big.Int)
	xBytes, err := xHash.MarshalText()
	if err != nil {
		return 0, err
	}
	x.UnmarshalText(xBytes)
	gb := new(bn256.G1)
	gb.Add(CL, CR.ScalarMult(CR, x.Neg(x)))

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
