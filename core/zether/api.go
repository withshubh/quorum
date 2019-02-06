package zether

import (
	"bytes"
	"encoding/binary"
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

func (api *PublicZetherAPI) ReadBalance(CLBytes [2]common.Hash, CRBytes [2]common.Hash, xHash common.Hash, startFloat float64, endFloat float64) (int64, error) {
	if int64(startFloat) < 0 || int64(endFloat) >= big.MaxPrec {
		return 0, errors.New("Invalid search range!")
	}
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
	end := big.NewInt(int64(endFloat))
	for i := big.NewInt(int64(startFloat)); i.Cmp(end) < 0; i.Add(i, one) {
		test := new(bn256.G1).ScalarBaseMult(i)
		if bytes.Compare(test.Marshal(), gb.Marshal()) == 0 {
			return i.Int64(), nil
		}
	}
	return 0, errors.New("Balance decryption failed!")
}

func (api *PublicZetherAPI) CreateTransfer(CL [2]common.Hash, CR [2]common.Hash, y [2]common.Hash, yBar [2]common.Hash, x common.Hash, bTransfer float64, bDiff float64) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bTransferBytes := make([]byte, 32)
	bDiffBytes := make([]byte, 32)
	binary.PutUvarint(bTransferBytes, uint64(bTransfer))
	binary.PutUvarint(bDiffBytes, uint64(bDiff))
	// consider sending these explicitly as uints instead of bytes

	// java.createTransfer(append(CL[0].Bytes(), CL[1].Bytes()...), append(CR[0].Bytes(), CR[1].Bytes()...), append(y[0].Bytes(), y[1].Bytes()...), append(yBar[0].Bytes(), yBar[1].Bytes()...), x.Bytes(), bTransferBytes, bDiffBytes)

	return result, nil
}

func (api *PublicZetherAPI) CreateBurn(CL [2]common.Hash, CR [2]common.Hash, y [2]common.Hash, bTransfer float64, x common.Hash, bDiff float64) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bTransferBytes := make([]byte, 32)
	bDiffBytes := make([]byte, 32)
	binary.PutUvarint(bTransferBytes, uint64(bTransfer))
	binary.PutUvarint(bDiffBytes, uint64(bDiff))
	// consider sending these explicitly as uints instead of bytes

	// java.createTransfer(append(CL[0].Bytes(), CL[1].Bytes()...), append(CR[0].Bytes(), CR[1].Bytes()...), append(y[0].Bytes(), y[1].Bytes()...), bTransferBytes, x.Bytes(), bDiffBytes)

	return result, nil
}
