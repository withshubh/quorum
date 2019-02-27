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

func (api *PublicZetherAPI) ProveTransfer(CLBytes [][2]common.Hash, CRBytes [][2]common.Hash, yBytes [][2]common.Hash, x common.Hash, bTransfer uint64, bDiff uint64, outIndex uint64, inIndex uint64) (map[string]interface{}, error) {
	// note: CL and CR here are before the debits are done, whereas verification takes them after the debits are done.
	// a bit weird, but makes sense: the contract will have to do them "anyway", whereas geth javascript will not.
	result := make(map[string]interface{})

	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	size := yBytes.length
	CL := make([]string, size)
	CR := make([]string, size)
	y := make([]string, size)
	for i := 0; i < size; i++ {
		CL[i] = hexutil.Encode(append(CLBytes[i][0].Bytes(), CLBytes[i][1].Bytes()))
		CR[i] = hexutil.Encode(append(CRBytes[i][0].Bytes(), CRBytes[i][1].Bytes()))
		y[i] = hexutil.Encode(append(yBytes[i][0].Bytes(), yBytes[i][1].Bytes()))
	}

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-transfer", nil)
	q := req.URL.Query()
	q.Add("CL", CL)
	q.Add("CR", CR)
	q.Add("y", y)
	q.Add("x", hexutil.Encode(x.Bytes()))
	q.Add("r", common.BytesToHash(r.Bytes()).Hex())
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	q.add("outIndex", hexutil.EncodeUint64(outIndex))
	q.add("inIndex", hexutil.EncodeUint64(inIndex))

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
	result["L"] = make([][2]common.Hash, size)
	for i := 0; i < size; i++ {
		L := new(bn256.G1)
		L.Unmarshal(yBytes[i])
		L.ScalarMult(L, r)
		if i == outIndex {
			gOut := new(bn256.G1)
			gOut.Unmarshal(gBytes)
			gOut.ScalarMult(big.NewInt(int64(-bTransfer)))
			L.Add(L, gOut)
		} else if i == inIndex {
			gIn := new(bn256.G1)
			gIn.ScalarMult(big.NewInt(int64(bTransfer)))
			L.Add(L, gIn)
		}
		result["L"][i] = [2]common.Hash{common.BytesToHash(L.Marshal()[:32]), common.BytesToHash(L.Marshal()[32:])}
	}
	R := new(bn256.G1)
	R.Unmarshal(gBytes)
	R.ScalarMult(R, r)
	if err != nil {
		return nil, err
	}
	result["R"] = [2]common.Hash{common.BytesToHash(R.Marshal()[:32]), common.BytesToHash(R.Marshal()[32:])}
	result["proof"] = proof // will have to concatenate the shuffle proof!

	return result, nil
}

func (api *PublicZetherAPI) ProveBurn(CL [2]common.Hash, CR [2]common.Hash, y [2]common.Hash, bTransfer uint64, x common.Hash, bDiff uint64) (interface{}, error) {
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
