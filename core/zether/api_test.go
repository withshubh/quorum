package zether

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestData(t *testing.T) {

	y := new(bn256.G1)
	yBytes, _ := hexutil.Decode("0x121f01c904b8502432b2e312f65fee3f63f91aae615da325d9b4d72a5e3ccafe157037e4a46359911bb30107731ec3d9a05d3749e79174e7a4c57670d2872f38")
	y.Unmarshal(yBytes)
	fmt.Println("y: " + y.String())

	CL := new(bn256.G1)
	CL.ScalarBaseMult(big.NewInt(int64(100)))
	CL.Add(CL, y)
	fmt.Println("CL: " + CL.String())

	CR := new(bn256.G1)
	CRBytes, _ := hexutil.Decode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002")
	CR.Unmarshal(CRBytes)
	fmt.Println("CR: " + CR.String())

	yBar := new(bn256.G1)
	yBarBytes, _ := hexutil.Decode("0x2c4fdb29d468bbacba5ae1f67c6a314f8b3724c64542507361eee196afe120df0fbf548c7a6a7ad0ab943df7a8ae1dc8ea873d482092a1a0f28df51ea0cdda3d")
	yBar.Unmarshal(yBarBytes)
	fmt.Println("yBar: " + yBar.String())

	myRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, _, _ := bn256.RandomG1(myRand)
	fmt.Println("r: " + r.String() + ", "+ common.BytesToHash(r.Bytes()).Hex())

	x := new(big.Int)
	xBytes, _ := hexutil.Bytes("0x1045af016c96d8fd062c342197d3696f0216c8181a38cbaec855b043abcf84e4").MarshalText()
	x.UnmarshalText(xBytes)
	fmt.Println("x: " + hexutil.Encode(x.Bytes()))

	fmt.Println(hexutil.MustDecodeBig("0x1c71d94150251d1229bc46397bf7fcf6fd0539cdd6c138c389862e11ee75e05e"))
	inOutR := new(bn256.G1)
	inOutR.ScalarBaseMult(hexutil.MustDecodeBig("0x1c71d94150251d1229bc46397bf7fcf6fd0539cdd6c138c389862e11ee75e05e"))
	fmt.Println("inOutR: " + inOutR.String())

}