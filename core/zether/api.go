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
		"response" : string(resp_body),
		"status" : "success",
	}, nil
}

func (api *PublicZetherAPI) CreateAccount() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	myRand := rand.New(rand.NewSource(time.Now().UnixNano())) // can i use the default source?
	x, y, err := bn256.RandomG1(myRand)
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
	gb.Add(CL, gb.ScalarMult(CR, x.Neg(x)))

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

func (api *PublicZetherAPI) Add(aBytes [2]common.Hash, bBytes [2]common.Hash) ([2]common.Hash, error) {
	A := new(bn256.G1)
	if _, err := A.Unmarshal(append(aBytes[0].Bytes(), aBytes[1].Bytes()...)); err != nil {
		return [2]common.Hash{common.BytesToHash(make([]byte, 32)), common.BytesToHash(make([]byte, 32))}, err
	}
	B := new(bn256.G1)
	if _, err := B.Unmarshal(append(bBytes[0].Bytes(), bBytes[1].Bytes()...)); err != nil {
		return [2]common.Hash{common.BytesToHash(make([]byte, 32)), common.BytesToHash(make([]byte, 32))}, err
	}

	sum := new(bn256.G1)
	sum.Add(A, B)
	sumBytes := sum.Marshal()

	return [2]common.Hash{common.BytesToHash(sumBytes[:32]), common.BytesToHash(sumBytes[32:])}, nil
}

func (api *PublicZetherAPI) ProveTransfer(CLBytes [2]common.Hash, CRBytes [2]common.Hash, yHash [2]common.Hash, yBarHash [2]common.Hash, xHash common.Hash, bTransfer uint64, bDiff uint64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	CL := new(bn256.G1)
	if _, err := CL.Unmarshal(append(CLBytes[0].Bytes(), CLBytes[1].Bytes()...)); err != nil {
		return nil, err
	}
	CR := new(bn256.G1)
	if _, err := CR.Unmarshal(append(CRBytes[0].Bytes(), CRBytes[1].Bytes()...)); err != nil {
		return nil, err
	}
	y := new(bn256.G1)
	if _, err := y.Unmarshal(append(yHash[0].Bytes(), yHash[1].Bytes()...)); err != nil {
		return nil, err
	}
	yBar := new(bn256.G1)
	if _, err := yBar.Unmarshal(append(yBarHash[0].Bytes(), yBarHash[1].Bytes()...)); err != nil {
		return nil, err
	}
	x := new(big.Int)
	xBytes, err := xHash.MarshalText()
	if err != nil {
		return nil, err
	}
	x.UnmarshalText(xBytes)

	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, inOutR, err := bn256.RandomG1(myRand)
	if err != nil {
		return nil, err
	}

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-transfer", nil)
	q := req.URL.Query()
	q.Add("CL", hexutil.Encode(append(CLBytes[0].Bytes(), CLBytes[1].Bytes()...)))
	q.Add("CR", hexutil.Encode(append(CRBytes[0].Bytes(), CRBytes[1].Bytes()...)))
	q.Add("y", hexutil.Encode(append(yHash[0].Bytes(), yHash[1].Bytes()...)))
	q.Add("yBar", hexutil.Encode(append(yBarHash[0].Bytes(), yBarHash[1].Bytes()...)))
	q.Add("x", hexutil.Encode(x.Bytes()))
	q.Add("r", common.BytesToHash(r.Bytes()).Hex())
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute at server")
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	proof := string(resp_body)

	// proof := java.proveTransfer(append(CL[0].Bytes(), CL[1].Bytes()...), append(CR[0].Bytes(), CR[1].Bytes()...), append(y[0].Bytes(), y[1].Bytes()...), append(yBar[0].Bytes(), yBar[1].Bytes()...), x.Bytes(), r.Bytes(), bTransfer.Bytes(), bDiff.Bytes())
	// warning: calling .Bytes() could yeild a slice of < 32 length. make sure this is ok with the RPC call, otherwise make([]byte, 32) beforehand and use PutUvarint
	//proof := common.Proof(make([]byte, 1216))

	gbTransfer := new(bn256.G1) // _recompute_ the following, which were computed within proveTransfer...
	gbTransfer.ScalarBaseMult(big.NewInt(int64(bTransfer)))
	outL := y.Add(gbTransfer, y.ScalarMult(y, r))         // base in inner expression is dummy
	inL := yBar.Add(gbTransfer, yBar.ScalarMult(yBar, r)) // value won't be used

	result["outL"] = [2]common.Hash{common.BytesToHash(outL.Marshal()[:32]), common.BytesToHash(outL.Marshal()[32:])}
	result["inL"] = [2]common.Hash{common.BytesToHash(inL.Marshal()[:32]), common.BytesToHash(inL.Marshal()[32:])}
	result["inOutR"] = [2]common.Hash{common.BytesToHash(inOutR.Marshal()[:32]), common.BytesToHash(inOutR.Marshal()[32:])}
	result["proof"] = proof // if had js elliptic packages, could just return only the proof and recompute the rest in web3

	return result, nil
}

func (api *PublicZetherAPI) ProveBurn(CLBytes [2]common.Hash, CRBytes [2]common.Hash, yHash [2]common.Hash, bTransfer uint64, x common.Hash, bDiff uint64) (interface{}, error) {
	bTransferBytes := make([]byte, 32)
	bDiffBytes := make([]byte, 32)
	binary.PutUvarint(bTransferBytes, uint64(bTransfer))
	binary.PutUvarint(bDiffBytes, uint64(bDiff))
	// consider sending these explicitly as uints instead of bytes

	CL := new(bn256.G1)
	if _, err := CL.Unmarshal(append(CLBytes[0].Bytes(), CLBytes[1].Bytes()...)); err != nil {
		return nil, err
	}
	CR := new(bn256.G1)
	if _, err := CR.Unmarshal(append(CRBytes[0].Bytes(), CRBytes[1].Bytes()...)); err != nil {
		return nil, err
	}
	y := new(bn256.G1)
	if _, err := y.Unmarshal(append(yHash[0].Bytes(), yHash[1].Bytes()...)); err != nil {
		return nil, err
	}

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-burn", nil)
	q := req.URL.Query()
	q.Add("CL", hexutil.Encode(append(CLBytes[0].Bytes(), CLBytes[1].Bytes()...)))
	q.Add("CR", hexutil.Encode(append(CRBytes[0].Bytes(), CRBytes[1].Bytes()...)))
	q.Add("y", hexutil.Encode(append(yHash[0].Bytes(), yHash[1].Bytes()...)))
	q.Add("x", hexutil.Encode(x.Bytes()))
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to execute at server")
	}
	resp_body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	proof := string(resp_body)

	// proof := java.proveBurn(append(CL[0].Bytes(), CL[1].Bytes()...), append(CR[0].Bytes(), CR[1].Bytes()...), append(y[0].Bytes(), y[1].Bytes()...), bTransferBytes, x.Bytes(), bDiffBytes)
	//proof := common.Proof(make([]byte, 1184))

	return proof, nil
}
