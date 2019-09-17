// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

// NewEVMContext creates a new context for use in the EVM.
func NewEVMContext(msg Message, header *types.Header, chain ChainContext, author *common.Address) vm.Context {
	// If we don't have an explicit author (i.e. not mining), extract from the header
	var beneficiary common.Address
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	return vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		GetNano:     GetNanoFn(header, chain),
		Origin:      msg.From(),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		GasLimit:    header.GasLimit,
		GasPrice:    new(big.Int).Set(msg.GasPrice()),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	// in the original GetHashFn, the passed-in element `ref` was NOT included in the list of headers considered.
	// since our getHeaderFn below _does_ consider the most recent (non-strict inequality), we have to
	// artificially bump `ref` up by one before passing it over to getHeaderFn to preserve semantics. --BD
	return func(n uint64) common.Hash {
		header := getHeaderFn(ref, chain)(n)
		if header == nil || n == ref.Number.Uint64() {
			return common.Hash{}
		}
		return header.Hash()
	}
}

func GetNanoFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	return func(n uint64) common.Hash {
		header := getHeaderFn(ref, chain)(n)
		if header == nil {
			return common.Hash{}
		}
		extraDataBytes := header.Extra[types.ExtraVanity:types.ExtraVanity+types.ExtraDataLen]
		var extraData *types.ExtraData
		if err := rlp.DecodeBytes(extraDataBytes, &extraData); err != nil {
			return common.Hash{}
		}
		var result common.Hash
		copy(result[24:], extraData.NanoTime)
		return result
	}
}

func getHeaderFn(ref *types.Header, chain ChainContext) func(n uint64) *types.Header {
	var cache map[uint64]*types.Header

	return func(n uint64) *types.Header {
		// If there's no hash cache yet, make one
		if cache == nil {
			cache = map[uint64]*types.Header{} // why not initialize it right away above? --BD
		}
		// Try to fulfill the request from the cache
		if header, ok := cache[n]; ok {
			return header
		}
		// Not cached, iterate the blocks and cache the hashes
		for header := ref; header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
			cache[header.Number.Uint64()] = header
			if n == header.Number.Uint64() {
				return header
			}
		}
		return nil
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
