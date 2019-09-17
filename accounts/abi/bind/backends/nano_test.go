package backends

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type NanoTimeTestSuite struct {
	suite.Suite
	auth     *bind.TransactOpts
	address  common.Address
	gAlloc   core.GenesisAlloc
	sim      *SimulatedBackend
	nanotime *NanoTime
}

func TestNanoTime(t *testing.T) {
	suite.Run(t, new(NanoTimeTestSuite))
}

func (s *NanoTimeTestSuite) SetupTest() {
	key, _ := crypto.GenerateKey()
	s.auth = bind.NewKeyedTransactor(key)

	s.address = s.auth.From
	s.gAlloc = map[common.Address]core.GenesisAccount{
		s.address: {Balance: big.NewInt(10000000000)},
	}

	s.sim = NanoSimulatedBackend(s.gAlloc, 1000000)

	_, _, nano, e := DeployNanoTime(s.auth, s.sim)
	s.Nil(e)
	s.nanotime = nano
	s.sim.Commit()
}

func (s *NanoTimeTestSuite) TestTime() {
	time, err := s.nanotime.Timestamp(nil)
	s.Nil(err)
	s.NotEqual(time.Uint64(), uint64(0))
}
