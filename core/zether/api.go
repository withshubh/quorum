package zether

import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type PublicZetherAPI struct {
}

func NewPublicZetherAPI() *PublicZetherAPI {
	return &PublicZetherAPI{}
}

func (api *PublicZetherAPI) TestConnect(b uint) (map[string]string, error) {
	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	x, _, err := bn256.RandomG1(myRand)
	// test java connectivity
	req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	q := req.URL.Query()
	q.Add("a", common.BytesToHash(x.Bytes()).Hex())
	q.Add("b", fmt.Sprint(b))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute at server")
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return map[string]string{
		"response": string(resp_body),
		"status":   "success",
	}, nil
}

func (api *PublicZetherAPI) CreateAccount() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	myRand := rand.New(rand.NewSource(time.Now().UnixNano())) // can i use the default source?
	x, _, err := bn256.RandomG1(myRand)
	y := new(bn256.G1)
	gBytes, _ := hexutil.Decode("0x077da99d806abd13c9f15ece5398525119d11e11e9836b2ee7d23f6159ad87d401485efa927f2ad41bff567eec88f32fb0a0f706588b4e41a8d587d008b7f875")
	y.Unmarshal(gBytes)
	y.ScalarMult(y, x)
	if err != nil {
		return nil, err
	}
	result["x"] = common.BytesToHash(x.Bytes())
	result["y"] = [2]common.Hash{common.BytesToHash(y.Marshal()[:32]), common.BytesToHash(y.Marshal()[32:])}
	return result, nil
}

func (api *PublicZetherAPI) ReadBalance(CLBytes [2]common.Hash, CRBytes [2]common.Hash, xHash common.Hash, start int64, endInt int64) (int64, error) {
	// using int64, not uint64, for args... make sure nothing goes wrong here
	if start < 0 || endInt > big.MaxPrec {
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
	gb.Add(CL, gb.ScalarMult(CR, x.Neg(x)))

	one := big.NewInt(1)
	end := big.NewInt(endInt)
	gBytes, _ := hexutil.Decode("0x077da99d806abd13c9f15ece5398525119d11e11e9836b2ee7d23f6159ad87d401485efa927f2ad41bff567eec88f32fb0a0f706588b4e41a8d587d008b7f875")
	for i := big.NewInt(start); i.Cmp(end) <= 0; i.Add(i, one) {
		test := new(bn256.G1)
		test.Unmarshal(gBytes)
		test.ScalarMult(test, i)
		if bytes.Compare(test.Marshal(), gb.Marshal()) == 0 {
			return i.Int64(), nil
		}
	}
	return 0, errors.New("Balance decryption failed!")
}

func (api *PublicZetherAPI) ProveTransfer(CLBytes [][2]common.Hash, CRBytes [][2]common.Hash, yBytes [][2]common.Hash, gEpoch [2]common.Hash, xBytes common.Hash, bTransfer uint64, bDiff uint64, index []uint64) (map[string]interface{}, error) {
	// note: CL and CR here are before the debits are done, whereas verification takes them after the debits are done.
	// a bit weird, but makes sense: the contract will have to do them "anyway", whereas geth javascript will not.
	// edit: actually, CL and CR should represent the state after any "overdue rollovers" are done, but before the debits are done.
	result := make(map[string]interface{})

	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, _, err := bn256.RandomG1(myRand)
	size := len(yBytes)
	var CL strings.Builder
	var CR strings.Builder
	var y strings.Builder
	for i := 0; i < size; i++ {
		CL.WriteString(hexutil.Encode(append(CLBytes[i][0].Bytes(), CLBytes[i][1].Bytes()...)))
		CR.WriteString(hexutil.Encode(append(CRBytes[i][0].Bytes(), CRBytes[i][1].Bytes()...)))
		y.WriteString(hexutil.Encode(append(yBytes[i][0].Bytes(), yBytes[i][1].Bytes()...)))
	}

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-transfer", nil)
	q := req.URL.Query()
	q.Add("CL", CL.String())
	q.Add("CR", CR.String())
	q.Add("y", y.String())
	q.Add("gEpoch", hexutil.Encode(append(gEpoch[0].Bytes(), gEpoch[1].Bytes()...)))
	q.Add("x", hexutil.Encode(xBytes.Bytes()))
	q.Add("r", common.BytesToHash(r.Bytes()).Hex())
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	q.Add("index", hexutil.EncodeUint64(index[0])+hexutil.EncodeUint64(index[1]))

	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to execute Java request.")
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	proof := string(resp_body)

	gBytes, _ := hexutil.Decode("0x077da99d806abd13c9f15ece5398525119d11e11e9836b2ee7d23f6159ad87d401485efa927f2ad41bff567eec88f32fb0a0f706588b4e41a8d587d008b7f875")
	L := make([][2]common.Hash, size)
	for i := 0; i < size; i++ {
		L_i := new(bn256.G1)
		L_i.Unmarshal(append(yBytes[i][0].Bytes(), yBytes[i][1].Bytes()...))
		L_i.ScalarMult(L_i, r)
		if uint64(i) == index[0] {
			gOut := new(bn256.G1)
			gOut.Unmarshal(gBytes)
			gOut.ScalarMult(gOut, big.NewInt(int64(-bTransfer)))
			L_i.Add(L_i, gOut)
		} else if uint64(i) == index[1] {
			gIn := new(bn256.G1)
			gIn.ScalarMult(gIn, big.NewInt(int64(bTransfer)))
			L_i.Add(L_i, gIn)
		}
		L[i] = [2]common.Hash{common.BytesToHash(L_i.Marshal()[:32]), common.BytesToHash(L_i.Marshal()[32:])}
	}
	R := new(bn256.G1)
	R.Unmarshal(gBytes)
	R.ScalarMult(R, r)
	if err != nil {
		return nil, err
	}
	x := new(big.Int)
	x.SetBytes(xBytes.Bytes())
	u := new(bn256.G1)
	u.Unmarshal(append(gEpoch[0].Bytes(), gEpoch[1].Bytes()...))
	u.ScalarMult(u, x)
	result["L"] = L
	result["R"] = [2]common.Hash{common.BytesToHash(R.Marshal()[:32]), common.BytesToHash(R.Marshal()[32:])}
	result["u"] = [2]common.Hash{common.BytesToHash(u.Marshal()[:32]), common.BytesToHash(u.Marshal()[32:])}
	// ^^^ again, all of these extra bits would be unnecessary if elliptic curve ops could be performed directly in javascript.
	// as another alternative, could add ec operations to the zether namespace, but this is sort of clunky.
	result["proof"] = proof

	return result, nil
}

func (api *PublicZetherAPI) ProveBurn(CL [2]common.Hash, CR [2]common.Hash, y [2]common.Hash, bTransfer uint64, gEpoch [2]common.Hash, x common.Hash, bDiff uint64) (interface{}, error) {
	bTransferBytes := make([]byte, 32)
	bDiffBytes := make([]byte, 32)
	binary.PutUvarint(bTransferBytes, uint64(bTransfer))
	binary.PutUvarint(bDiffBytes, uint64(bDiff))
	// consider sending these explicitly as uints instead of bytes

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-burn", nil)
	q := req.URL.Query()
	q.Add("CL", hexutil.Encode(append(CL[0].Bytes(), CL[1].Bytes()...)))
	q.Add("CR", hexutil.Encode(append(CR[0].Bytes(), CR[1].Bytes()...)))
	q.Add("gEpoch", hexutil.Encode(append(gEpoch[0].Bytes(), gEpoch[1].Bytes()...)))
	q.Add("y", hexutil.Encode(append(y[0].Bytes(), y[1].Bytes()...)))
	q.Add("x", hexutil.Encode(x.Bytes()))
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to execute Java request.")
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	proof := string(resp_body)

	return proof, nil
}
