package raft

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestSignHeader(t *testing.T) {
	//create only what we need to test the seal
	var testRaftId uint16 = 5
	config := &node.Config{Name: "unit-test", DataDir: ""}

	nodeKey := config.NodeKey()

	raftProtocolManager := &ProtocolManager{raftId: testRaftId}
	raftService := &RaftService{nodeKey: nodeKey, raftProtocolManager: raftProtocolManager}
	minter := minter{eth: raftService}

	//create some fake header to sign
	fakeParentHash := common.HexToHash("0xc2c1dc1be8054808c69e06137429899d")

	now := time.Now()
	header := &types.Header{
		ParentHash: fakeParentHash,
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(1),
		GasLimit:   uint64(0),
		GasUsed:    uint64(0),
		Coinbase:   minter.coinbase,
		Time:       big.NewInt(now.Unix()),
	}

	headerHash := header.Hash()
	nanotime := make([]byte, 8)
	binary.BigEndian.PutUint64(nanotime, uint64(now.UnixNano()))
	extraData := &types.ExtraData{NanoTime: nanotime}
	extraDataBytes, err := rlp.EncodeToBytes(extraData)
	if err != nil {
		t.Errorf("RLP encoding of extra data struct failed!")
	}
	extraSealBytes := minter.buildExtraSeal(headerHash, extraDataBytes)
	var seal *extraSeal
	err = rlp.DecodeBytes(extraSealBytes[:], &seal)
	if err != nil {
		t.Fatalf("Unable to decode seal: %s", err.Error())
	}

	// Check raftId
	sealRaftId := binary.LittleEndian.Uint16(seal.RaftId)
	if sealRaftId != testRaftId {
		t.Errorf("RaftID does not match. Expected: %d, Actual: %d", testRaftId, sealRaftId)
	}

	//Identify who signed it
	sig := seal.Signature
	hw := sha3.NewKeccak256()
	hw.Write(headerHash.Bytes()) // write the header hash
	hw.Write(extraDataBytes) // write the extra data
	var total common.Hash
	hw.Sum(total[:])
	pubKey, err := crypto.SigToPub(total.Bytes(), sig)
	if err != nil {
		t.Fatalf("Unable to get public key from signature: %s", err.Error())
	}

	//Compare derived public key to original public key
	if pubKey.X.Cmp(nodeKey.X) != 0 {
		t.Errorf("Signature incorrect!")
	}

}
