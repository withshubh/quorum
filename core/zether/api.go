package zether

import "C"
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
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
		return nil, errors.New("Failed to execute Java request.")
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

func (api *PublicZetherAPI) Add(aBytes [2][2]common.Hash, bBytes [2][2]common.Hash) [2][2]common.Hash {
	// adds two elemengs of G x G.
	aL := new(bn256.G1)
	aL.Unmarshal(append(aBytes[0][0].Bytes(), aBytes[0][1].Bytes()...))
	aR := new(bn256.G1)
	aR.Unmarshal(append(aBytes[1][0].Bytes(), aBytes[1][1].Bytes()...))
	bL := new(bn256.G1)
	bL.Unmarshal(append(bBytes[0][0].Bytes(), bBytes[0][1].Bytes()...))
	bR := new(bn256.G1)
	bR.Unmarshal(append(bBytes[1][0].Bytes(), bBytes[1][1].Bytes()...))

	resultL := new(bn256.G1)
	resultL.Add(aL, bL)
	resultR := new(bn256.G1)
	resultR.Add(aR, bR)
	resultLBytes := resultL.Marshal()
	resultRBytes := resultR.Marshal()

	return [2][2]common.Hash{{common.BytesToHash(resultLBytes[:32]), common.BytesToHash(resultLBytes[32:])}, {common.BytesToHash(resultRBytes[:32]), common.BytesToHash(resultRBytes[32:])}}
}

func (api *PublicZetherAPI) ReadBalance(CBytes [2][2]common.Hash, xHash common.Hash, start int64, endInt int64) (int64, error) {
	// no longer checking whether start >= 0.
	if endInt > big.MaxPrec {
		return 0, errors.New("Invalid search range!")
	}
	CL := new(bn256.G1)
	if _, err := CL.Unmarshal(append(CBytes[0][0].Bytes(), CBytes[0][1].Bytes()...)); err != nil {
		return 0, err
	}
	CR := new(bn256.G1)
	if _, err := CR.Unmarshal(append(CBytes[1][0].Bytes(), CBytes[1][1].Bytes()...)); err != nil {
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

func computeU(input string, xBytes common.Hash) bn256.G1 { // could also conceivably just return the raw bytes...
	p := hexutil.MustDecodeBig("0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47") // field order
	x := new(big.Int)
	x.SetBytes(xBytes.Bytes())

	seed := new(big.Int)
	seed.SetBytes(crypto.Keccak256([]byte(input)))
	seed.Mod(seed, p)
	y := new(big.Int)
	for {
		ySquared := y.Add(y.Exp(seed, big.NewInt(3), p), big.NewInt(3)) // throw away base values
		if y.ModSqrt(ySquared, p) != nil {                              // y.Exp(ySquared, y.Div(y.Add(p, big.NewInt(1)), big.NewInt(4)), p) // why doesn't this work?!?!?!?
			break
		}
		seed.Add(seed, big.NewInt(1))
	}
	seedBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	copy(seedBytes[32-len(seed.Bytes()):], seed.Bytes()) // right-justify
	copy(yBytes[32-len(y.Bytes()):], y.Bytes())          // right-justify
	result := new(bn256.G1)
	result.Unmarshal(append(seedBytes, yBytes...))
	result.ScalarMult(result, x)
	return *result
}

func (api *PublicZetherAPI) ProveTransfer(CLBytes [][2]common.Hash, CRBytes [][2]common.Hash, yBytes [][2]common.Hash, epoch uint64, x common.Hash, bTransfer uint64, bDiff uint64, index []uint64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, _, err := bn256.RandomG1(myRand)
	size := len(yBytes)
	var CL strings.Builder
	var CR strings.Builder
	var y strings.Builder
	CL.WriteString("0x")
	CR.WriteString("0x")
	y.WriteString("0x")
	for i := 0; i < size; i++ {
		CL.WriteString(hexutil.Encode(append(CLBytes[i][0].Bytes(), CLBytes[i][1].Bytes()...))[2:])
		CR.WriteString(hexutil.Encode(append(CRBytes[i][0].Bytes(), CRBytes[i][1].Bytes()...))[2:])
		y.WriteString(hexutil.Encode(append(yBytes[i][0].Bytes(), yBytes[i][1].Bytes()...))[2:])
	}

	// RPC to Java service
	req, _ := http.NewRequest("GET", "http://localhost:8080/prove-transfer", nil)
	q := req.URL.Query()
	q.Add("CL", CL.String())
	q.Add("CR", CR.String())
	q.Add("y", y.String())
	q.Add("epoch", hexutil.EncodeUint64(epoch))
	q.Add("x", hexutil.Encode(x.Bytes()))
	q.Add("r", common.BytesToHash(r.Bytes()).Hex())
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("bDiff", hexutil.EncodeUint64(bDiff))
	q.Add("outIndex", hexutil.EncodeUint64(index[0]))
	q.Add("inIndex", hexutil.EncodeUint64(index[1]))

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
			gIn.Unmarshal(gBytes)
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
	u := computeU("Zether "+strconv.FormatUint(epoch, 10), x)
	result["L"] = L
	result["R"] = [2]common.Hash{common.BytesToHash(R.Marshal()[:32]), common.BytesToHash(R.Marshal()[32:])}
	result["u"] = [2]common.Hash{common.BytesToHash(u.Marshal()[:32]), common.BytesToHash(u.Marshal()[32:])}
	// ^^^ again, all of these extra bits would be unnecessary if elliptic curve ops could be performed directly in javascript.
	result["proof"] = proof

	return result, nil
}

func (api *PublicZetherAPI) ProveBurn(CL [2]common.Hash, CR [2]common.Hash, y [2]common.Hash, bTransfer uint64, epoch uint64, x common.Hash, bDiff uint64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

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
	q.Add("bTransfer", hexutil.EncodeUint64(bTransfer))
	q.Add("epoch", hexutil.EncodeUint64(epoch))
	q.Add("x", hexutil.Encode(x.Bytes()))
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

	u := computeU("Zether "+strconv.FormatUint(epoch, 10), x)
	result["u"] = [2]common.Hash{common.BytesToHash(u.Marshal()[:32]), common.BytesToHash(u.Marshal()[32:])}
	result["proof"] = proof

	return result, nil
}
