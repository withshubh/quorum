// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package backends

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// NanoTimeABI is the input ABI used to generate the binding from.
const NanoTimeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"timestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// NanoTimeBin is the compiled bytecode used for deploying new contracts.
const NanoTimeBin = `0x6080604052348015600f57600080fd5b50609e8061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063b80777ea14602d575b600080fd5b60336045565b60408051918252519081900360200190f35b604051438082526000916020818181620100015afa606257600080fd5b519291505056fea265627a7a72315820d6f3ca06d5bdea7774e62e7c50780de42ce5428728a61a6fd604dfba0d64c45c64736f6c634300050b0032`

// DeployNanoTime deploys a new Ethereum contract, binding an instance of NanoTime to it.
func DeployNanoTime(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NanoTime, error) {
	parsed, err := abi.JSON(strings.NewReader(NanoTimeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NanoTimeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NanoTime{NanoTimeCaller: NanoTimeCaller{contract: contract}, NanoTimeTransactor: NanoTimeTransactor{contract: contract}, NanoTimeFilterer: NanoTimeFilterer{contract: contract}}, nil
}

// NanoTime is an auto generated Go binding around an Ethereum contract.
type NanoTime struct {
	NanoTimeCaller     // Read-only binding to the contract
	NanoTimeTransactor // Write-only binding to the contract
	NanoTimeFilterer   // Log filterer for contract events
}

// NanoTimeCaller is an auto generated read-only Go binding around an Ethereum contract.
type NanoTimeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NanoTimeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NanoTimeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NanoTimeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NanoTimeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NanoTimeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NanoTimeSession struct {
	Contract     *NanoTime         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NanoTimeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NanoTimeCallerSession struct {
	Contract *NanoTimeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// NanoTimeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NanoTimeTransactorSession struct {
	Contract     *NanoTimeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// NanoTimeRaw is an auto generated low-level Go binding around an Ethereum contract.
type NanoTimeRaw struct {
	Contract *NanoTime // Generic contract binding to access the raw methods on
}

// NanoTimeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NanoTimeCallerRaw struct {
	Contract *NanoTimeCaller // Generic read-only contract binding to access the raw methods on
}

// NanoTimeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NanoTimeTransactorRaw struct {
	Contract *NanoTimeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNanoTime creates a new instance of NanoTime, bound to a specific deployed contract.
func NewNanoTime(address common.Address, backend bind.ContractBackend) (*NanoTime, error) {
	contract, err := bindNanoTime(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NanoTime{NanoTimeCaller: NanoTimeCaller{contract: contract}, NanoTimeTransactor: NanoTimeTransactor{contract: contract}, NanoTimeFilterer: NanoTimeFilterer{contract: contract}}, nil
}

// NewNanoTimeCaller creates a new read-only instance of NanoTime, bound to a specific deployed contract.
func NewNanoTimeCaller(address common.Address, caller bind.ContractCaller) (*NanoTimeCaller, error) {
	contract, err := bindNanoTime(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NanoTimeCaller{contract: contract}, nil
}

// NewNanoTimeTransactor creates a new write-only instance of NanoTime, bound to a specific deployed contract.
func NewNanoTimeTransactor(address common.Address, transactor bind.ContractTransactor) (*NanoTimeTransactor, error) {
	contract, err := bindNanoTime(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NanoTimeTransactor{contract: contract}, nil
}

// NewNanoTimeFilterer creates a new log filterer instance of NanoTime, bound to a specific deployed contract.
func NewNanoTimeFilterer(address common.Address, filterer bind.ContractFilterer) (*NanoTimeFilterer, error) {
	contract, err := bindNanoTime(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NanoTimeFilterer{contract: contract}, nil
}

// bindNanoTime binds a generic wrapper to an already deployed contract.
func bindNanoTime(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NanoTimeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NanoTime *NanoTimeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NanoTime.Contract.NanoTimeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NanoTime *NanoTimeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NanoTime.Contract.NanoTimeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NanoTime *NanoTimeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NanoTime.Contract.NanoTimeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NanoTime *NanoTimeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _NanoTime.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NanoTime *NanoTimeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NanoTime.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NanoTime *NanoTimeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NanoTime.Contract.contract.Transact(opts, method, params...)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() constant returns(result uint256)
func (_NanoTime *NanoTimeCaller) Timestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _NanoTime.contract.Call(opts, out, "timestamp")
	return *ret0, err
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() constant returns(result uint256)
func (_NanoTime *NanoTimeSession) Timestamp() (*big.Int, error) {
	return _NanoTime.Contract.Timestamp(&_NanoTime.CallOpts)
}

// Timestamp is a free data retrieval call binding the contract method 0xb80777ea.
//
// Solidity: function timestamp() constant returns(result uint256)
func (_NanoTime *NanoTimeCallerSession) Timestamp() (*big.Int, error) {
	return _NanoTime.Contract.Timestamp(&_NanoTime.CallOpts)
}
